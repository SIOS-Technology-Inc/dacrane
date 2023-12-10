package main

import (
	"bytes"
	"dacrane/pdk"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	config := pdk.NewDefaultPluginConfig()
	dockerHost := "/var/run/docker.sock"
	config.DockerHost = &dockerHost
	pdk.ExecPluginJob(pdk.Plugin{
		Config: config,
		Resources: pdk.MapToFunc(map[string]pdk.Resource{
			"container":    DockerContainerResource,
			"local-image":  DockerLocalImageResource,
			"remote-image": DockerRemoteImage,
		}),
	})
}

var DockerContainerResource = pdk.Resource{
	Create: func(parameter any, _ pdk.PluginMeta) (any, error) {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		name := params["name"].(string)
		env := params["env"].([]any)
		port := params["port"].(string)
		tag := params["tag"].(string)

		envOpts := []string{}
		for _, e := range env {
			name := e.(map[string]any)["name"].(string)
			value := e.(map[string]any)["value"].(string)
			opt := fmt.Sprintf(`-e "%s=%s"`, name, value)
			envOpts = append(envOpts, opt)
		}

		cmd := fmt.Sprintf("docker run -d --name %s -p %s %s %s:%s", name, port, strings.Join(envOpts, " "), image, tag)

		_, err := RunOnBash(cmd)
		if err != nil {
			panic(err)
		}

		return parameter, nil
	},
	Delete: func(parameter any, _ pdk.PluginMeta) error {
		params := parameter.(map[string]any)
		name := params["name"].(string)
		_, err := RunOnBash(fmt.Sprintf("docker stop %s", name))
		if err != nil {
			panic(err)
		}
		_, err = RunOnBash(fmt.Sprintf("docker rm %s", name))
		if err != nil {
			panic(err)
		}
		return nil
	},
}

var DockerLocalImageResource = pdk.Resource{
	Create: func(parameter any, _ pdk.PluginMeta) (any, error) {
		params := parameter.(map[string]any)
		dockerfile := params["dockerfile"].(string)
		image := params["image"].(string)
		tag := params["tag"].(string)

		dockerCmd := fmt.Sprintf("docker build -t %s:%s -f %s .", image, tag, dockerfile)
		_, err := RunOnBash(dockerCmd)
		return params, err
	},
	Delete: func(parameter any, _ pdk.PluginMeta) error {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		tag := params["tag"].(string)

		// remove local image
		dockerRmiCmd := fmt.Sprintf("docker rmi %s:%s", image, tag)
		_, err := RunOnBash(dockerRmiCmd)
		if err != nil {
			return err
		}
		return nil
	},
}

var DockerRemoteImage = pdk.Resource{
	Create: func(parameter any, _ pdk.PluginMeta) (any, error) {
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
			_, err := RunOnBash(cmd)
			if err != nil {
				return nil, err
			}
		}

		return params, nil
	},
	Delete: func(parameter any, _ pdk.PluginMeta) error {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		tag := params["tag"].(string)

		remote := params["remote"].(map[string](any))
		url := remote["url"].(string)
		user := remote["user"].(string)
		password := remote["password"].(string)

		// remove registry image
		dockerDigestCmd := fmt.Sprintf("docker images %s/%s --format {{.Digest}}", url, image)
		out, err := RunOnBash(dockerDigestCmd)
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
		res, err := RequestHttp(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		dockerRmiCmd := fmt.Sprintf("docker rmi %s/%s:%s", url, image, tag)
		_, err = RunOnBash(dockerRmiCmd)
		if err != nil {
			return err
		}

		return nil
	},
}

func RunOnBash(script string) ([]byte, error) {
	fmt.Printf("> %s\n", script)
	cmd := exec.Command("bash", "-c", script)
	writer := new(bytes.Buffer)
	cmd.Stdout = io.MultiWriter(os.Stdout, writer)
	cmd.Stderr = io.MultiWriter(os.Stderr, writer)
	err := cmd.Run()
	return writer.Bytes(), err
}

func RequestHttp(req *http.Request) (*http.Response, error) {
	fmt.Printf("> %s %s\n", req.Method, req.URL)
	res, err := http.DefaultClient.Do(req)
	fmt.Printf("> %s\n", res.Status)
	return res, err
}
