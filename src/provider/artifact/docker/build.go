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
