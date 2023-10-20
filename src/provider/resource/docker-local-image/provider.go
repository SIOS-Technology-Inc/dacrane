package docker

import (
	"dacrane/utils"
	"fmt"
)

type DockerArtifactProvider struct{}

func (DockerArtifactProvider) Create(params map[string]any) (map[string]any, error) {
	dockerfile := params["dockerfile"].(string)
	image := params["image"].(string)
	tag := params["tag"].(string)

	dockerCmd := fmt.Sprintf("docker build -t %s:%s -f %s .", image, tag, dockerfile)
	_, err := utils.RunOnBash(dockerCmd)
	return params, err
}

func (DockerArtifactProvider) Delete(params map[string]any) error {
	image := params["image"].(string)
	tag := params["tag"].(string)

	// remove local image
	dockerRmiCmd := fmt.Sprintf("docker rmi %s:%s", image, tag)
	_, err := utils.RunOnBash(dockerRmiCmd)
	if err != nil {
		return err
	}
	return nil
}
