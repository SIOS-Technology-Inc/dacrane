package main

import (
	"bytes"
	"dacrane/pdk"
	"fmt"
	"io"
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
			"shell": ShellResource,
		}),
	})
}

var ShellResource = pdk.Resource{
	Create: func(parameter any, meta pdk.PluginMeta) (any, error) {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		env := params["env"].([]any)
		tag := params["tag"].(string)
		shell := params["shell"].(string)
		script, ok := params["create"].(string)

		if !ok {
			return parameter, nil
		}

		envOpts := []string{}
		for _, e := range env {
			name := e.(map[string]any)["name"].(string)
			value := e.(map[string]any)["value"].(string)
			opt := fmt.Sprintf(`-e "%s=%s"`, name, value)
			envOpts = append(envOpts, opt)
		}

		netOpt := ""
		if network, ok := params["network"].(string); ok {
			netOpt = fmt.Sprintf("--net %s", network)
		}

		cmd := fmt.Sprintf(
			`docker run --rm -v $HOST_WORKING_DIR:/work %s %s %s:%s %s -c "%s"`,
			strings.Join(envOpts, " "), netOpt, image, tag, shell, script)

		_, err := RunOnSh(cmd, meta)
		if err != nil {
			panic(err)
		}

		return parameter, nil
	},
	Delete: func(parameter any, meta pdk.PluginMeta) error {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		env := params["env"].([]any)
		tag := params["tag"].(string)
		shell := params["shell"].(string)
		script, ok := params["delete"].(string)

		if !ok {
			return nil
		}

		envOpts := []string{}
		for _, e := range env {
			name := e.(map[string]any)["name"].(string)
			value := e.(map[string]any)["value"].(string)
			opt := fmt.Sprintf(`-e "%s=%s"`, name, value)
			envOpts = append(envOpts, opt)
		}

		netOpt := ""
		if network, ok := params["network"].(string); ok {
			netOpt = fmt.Sprintf("--net %s", network)
		}

		cmd := fmt.Sprintf(
			`docker run --rm -v $HOST_WORKING_DIR:/work %s %s %s:%s %s -c "%s"`,
			strings.Join(envOpts, " "), netOpt, image, tag, shell, script)

		_, err := RunOnSh(cmd, meta)
		if err != nil {
			panic(err)
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
