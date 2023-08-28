package docker

import (
	"fmt"
	"os/exec"
)

type DockerArtifactProvider struct{}

func (DockerArtifactProvider) Build(workingDir string, params map[string]any) ([]byte, error) {
	dockerfile := params["dockerfile"].(string)
	image := params["image"].(string)
	tag := params["tag"].(string)

	dockerCmd := fmt.Sprintf("docker build -t %s:%s -f %s .", image, tag, dockerfile)
	cmd := exec.Command("bash", "-c", dockerCmd)
	cmd.Dir = workingDir
	return cmd.CombinedOutput()
}

func (DockerArtifactProvider) Publish(workingDir string, params map[string]any) ([]byte, error) {
	image := params["image"].(string)
	tag := params["tag"].(string)
	repository := params["repository"].(map[string](any))
	url := repository["url"].(string)
	user := repository["user"].(string)
	password := repository["password"].(string)

	dockerLoginCmd := fmt.Sprintf("docker login -u %s -p %s %s", user, password, url)
	dockerImageTagCmd := fmt.Sprintf("docker image tag %s:%s %s/%s:%s", image, tag, url, image, tag)
	dockerPushCmd := fmt.Sprintf("docker image push %s/%s:%s", url, image, tag)

	cmds := []string{dockerLoginCmd, dockerImageTagCmd, dockerPushCmd}

	var log []byte
	for _, cmd := range cmds {
		cmd := exec.Command("bash", "-c", cmd)
		cmd.Dir = workingDir
		out, err := cmd.CombinedOutput()
		log = append(log, out...)
		if err != nil {
			return log, err
		}
	}

	return log, nil
}
