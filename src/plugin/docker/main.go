package main

import (
	"bytes"
	"dacrane/pdk"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	config := pdk.NewDefaultPluginConfig()
	dockerHost := "/var/run/docker.sock"
	config.DockerHost = &dockerHost
	pdk.ExecPluginJob(pdk.Plugin{
		Config: config,
		Resources: pdk.MapToFunc(map[string]pdk.Resource{
			"container":    DockerContainerResource,
			"network":      DockerNetworkResource,
			"local-image":  DockerLocalImageResource,
			"remote-image": DockerRemoteImage,
		}),
	})
}

var DockerContainerResource = pdk.Resource{
	Create: func(parameter any, meta pdk.PluginMeta) (any, error) {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		name := params["name"].(string)
		tag := params["tag"].(string)

		cmd := fmt.Sprintf("docker run -d --name %s", name)

		env, ok := params["env"].([]any)
		if ok {
			for _, e := range env {
				name := e.(map[string]any)["name"].(string)
				value := e.(map[string]any)["value"].(string)
				cmd = fmt.Sprintf(`%s -e "%s=%s"`, cmd, name, value)
			}
		}

		port, ok := params["port"].(string)
		if ok {
			cmd = fmt.Sprintf("%s -p %s", cmd, port)
		}

		network, ok := params["network"].(string)
		if ok {
			cmd = fmt.Sprintf("%s --net %s", cmd, network)
		}

		if healthcheck, ok := params["healthcheck"].(map[string]any); ok {
			cmd = fmt.Sprintf(`%s --health-cmd "%s¥"`, cmd, healthcheck["cmd"])
			cmd = fmt.Sprintf("%s --health-interval %s", cmd, healthcheck["interval"])
			cmd = fmt.Sprintf("%s --health-retries %s", cmd, healthcheck["retries"])
			cmd = fmt.Sprintf("%s --health-start-period %s", cmd, healthcheck["start_period"])
			cmd = fmt.Sprintf("%s --health-timeout %s", cmd, healthcheck["timeout"])
		}

		cmd = fmt.Sprintf("%s %s:%s", cmd, image, tag)

		_, err := RunOnSh(cmd, meta)
		if err != nil {
			panic(err)
		}

		// TODO Design for waiting
		if _, ok := params["healthcheck"].(map[string]any); ok {
			for i := 0; i <= 60; i++ {
				output, err := RunOnSh(fmt.Sprintf("docker inspect --format='{{json .State.Health}}' %s", name), meta)
				if err != nil {
					panic(err)
				}
				var res map[string]any
				json.Unmarshal(output, &res)
				if res["Status"] == "healthy" {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}

		return parameter, nil
	},
	Delete: func(parameter any, meta pdk.PluginMeta) error {
		params := parameter.(map[string]any)
		name := params["name"].(string)
		_, err := RunOnSh(fmt.Sprintf("docker stop %s", name), meta)
		if err != nil {
			panic(err)
		}
		_, err = RunOnSh(fmt.Sprintf("docker rm %s", name), meta)
		if err != nil {
			panic(err)
		}
		return nil
	},
}

var DockerNetworkResource = pdk.Resource{
	Create: func(parameter any, meta pdk.PluginMeta) (any, error) {
		params := parameter.(map[string]any)
		name := params["name"].(string)

		cmd := fmt.Sprintf("docker network create %s", name)

		_, err := RunOnSh(cmd, meta)
		if err != nil {
			panic(err)
		}

		return parameter, nil
	},
	Delete: func(parameter any, meta pdk.PluginMeta) error {
		params := parameter.(map[string]any)
		name := params["name"].(string)

		cmd := fmt.Sprintf("docker network rm %s", name)

		_, err := RunOnSh(cmd, meta)
		if err != nil {
			panic(err)
		}

		return nil
	},
}

var DockerLocalImageResource = pdk.Resource{
	Create: func(parameter any, meta pdk.PluginMeta) (any, error) {
		params := parameter.(map[string]any)
		dockerfile := params["dockerfile"].(string)
		image := params["image"].(string)
		tag := params["tag"].(string)

		dockerCmd := fmt.Sprintf("docker build -t %s:%s -f %s .", image, tag, dockerfile)
		_, err := RunOnSh(dockerCmd, meta)
		return params, err
	},
	Delete: func(parameter any, meta pdk.PluginMeta) error {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		tag := params["tag"].(string)

		// remove local image
		dockerRmiCmd := fmt.Sprintf("docker rmi %s:%s", image, tag)
		_, err := RunOnSh(dockerRmiCmd, meta)
		if err != nil {
			return err
		}
		return nil
	},
}

var DockerRemoteImage = pdk.Resource{
	Create: func(parameter any, meta pdk.PluginMeta) (any, error) {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		tag := params["tag"].(string)
		remote := params["remote"].(map[string]any)
		url := remote["url"].(string)
		user := remote["user"].(string)
		password := remote["password"].(string)

		dockerLoginCmd := fmt.Sprintf("docker login -u %s -p %s %s", user, password, url)
		dockerImageTagCmd := fmt.Sprintf("docker image tag %s:%s %s/%s:%s", image, tag, url, image, tag)
		dockerPushCmd := fmt.Sprintf("docker image push %s/%s:%s", url, image, tag)

		cmds := []string{dockerLoginCmd, dockerImageTagCmd, dockerPushCmd}

		for _, cmd := range cmds {
			_, err := RunOnSh(cmd, meta)
			if err != nil {
				return nil, err
			}
		}

		return params, nil
	},
	Delete: func(parameter any, meta pdk.PluginMeta) error {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		tag := params["tag"].(string)

		remote := params["remote"].(map[string](any))
		url := remote["url"].(string)
		user := remote["user"].(string)
		password := remote["password"].(string)

		// remove registry image
		dockerDigestCmd := fmt.Sprintf("docker images %s/%s --format {{.Digest}}", url, image)
		out, err := RunOnSh(dockerDigestCmd, meta)
		if err != nil {
			return err
		}
		digest := strings.ReplaceAll(string(out), "\n", "")

		// cf. https://docs.docker.com/registry/spec/api/#deleting-an-image
		deleteUrl := fmt.Sprintf("https://%s/v2/%s/manifests/%s", url, image, digest)
		req, err := http.NewRequest("DELETE", deleteUrl, nil)
		if err != nil {
			return err
		}
		req.SetBasicAuth(user, password)
		res, err := RequestHttp(req, meta)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		dockerRmiCmd := fmt.Sprintf("docker rmi %s/%s:%s", url, image, tag)
		_, err = RunOnSh(dockerRmiCmd, meta)
		if err != nil {
			return err
		}

		return nil
	},
}

func RunOnSh(script string, m pdk.PluginMeta) ([]byte, error) {
	m.Log(fmt.Sprintf("> %s\n", script))
	cmd := exec.Command("sh", "-c", script)
	writer := new(bytes.Buffer)
	cmd.Stdout = io.MultiWriter(os.Stderr, writer)
	cmd.Stderr = io.MultiWriter(os.Stderr, writer)
	err := cmd.Run()
	return writer.Bytes(), err
}

func RequestHttp(req *http.Request, m pdk.PluginMeta) (*http.Response, error) {
	m.Log(fmt.Sprintf("> %s %s\n", req.Method, req.URL))
	res, err := http.DefaultClient.Do(req)
	m.Log(fmt.Sprintf("> %s\n", res.Status))
	return res, err
}
