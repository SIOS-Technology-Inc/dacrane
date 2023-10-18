package docker

import (
	"dacrane/utils"
	"fmt"
	"strings"
)

type DockerResourceProvider struct{}

func (DockerResourceProvider) Create(parameters map[string]any) (map[string]any, error) {
	image := parameters["image"].(string)
	name := parameters["name"].(string)
	env := parameters["env"].([]any)
	port := parameters["port"].(string)
	tag := parameters["tag"].(string)

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

	return parameters, nil
}

func (DockerResourceProvider) Delete(parameters map[string]any) error {
	name := parameters["name"].(string)
	_, err := utils.RunOnBash(fmt.Sprintf("docker stop %s", name))
	if err != nil {
		panic(err)
	}
	_, err = utils.RunOnBash(fmt.Sprintf("docker rm %s", name))
	if err != nil {
		panic(err)
	}
	return nil
}
