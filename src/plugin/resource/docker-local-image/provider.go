package docker

import (
	"dacrane/pdk"
	"dacrane/utils"
	"fmt"
)

var DockerLocalImageResource = pdk.Resource{
	Create: func(parameter any, _ pdk.ProviderMeta) (any, error) {
		params := parameter.(map[string]any)
		dockerfile := params["dockerfile"].(string)
		image := params["image"].(string)
		tag := params["tag"].(string)

		dockerCmd := fmt.Sprintf("docker build -t %s:%s -f %s .", image, tag, dockerfile)
		_, err := utils.RunOnBash(dockerCmd)
		return params, err
	},
	Delete: func(parameter any, _ pdk.ProviderMeta) error {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		tag := params["tag"].(string)

		// remove local image
		dockerRmiCmd := fmt.Sprintf("docker rmi %s:%s", image, tag)
		_, err := utils.RunOnBash(dockerRmiCmd)
		if err != nil {
			return err
		}
		return nil
	},
}
