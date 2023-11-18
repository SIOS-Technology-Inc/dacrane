package docker

import (
	"dacrane/pdk"
	"dacrane/utils"
	"fmt"
	"strings"
)

var DockerContainerResource = pdk.Resource{
	Create: func(parameter any, _ pdk.ProviderMeta) (any, error) {
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

		_, err := utils.RunOnBash(cmd)
		if err != nil {
			panic(err)
		}

		return parameter, nil
	},
	Delete: func(parameter any, _ pdk.ProviderMeta) error {
		params := parameter.(map[string]any)
		name := params["name"].(string)
		_, err := utils.RunOnBash(fmt.Sprintf("docker stop %s", name))
		if err != nil {
			panic(err)
		}
		_, err = utils.RunOnBash(fmt.Sprintf("docker rm %s", name))
		if err != nil {
			panic(err)
		}
		return nil
	},
}
