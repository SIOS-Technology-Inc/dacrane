package docker

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
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

func (DockerArtifactProvider) Unpublish(workingDir string, params map[string]any) ([]byte, error) {
	image := params["image"].(string)
	tag := params["tag"].(string)
	repository := params["repository"].(map[string](any))
	url := repository["url"].(string)
	user := repository["user"].(string)
	password := repository["password"].(string)

	// remove registry image
	dockerDigestCmd := fmt.Sprintf("docker images %s/%s --format {{.Digest}}", url, image)
	cmd := exec.Command("bash", "-c", dockerDigestCmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, err
	}
	digest := strings.ReplaceAll(string(out), "\n", "")

	// cf. https://docs.docker.com/registry/spec/api/#deleting-an-image
	client := http.DefaultClient
	deleteUrl := fmt.Sprintf("https://%s/v2/%s/manifests/%s", url, image, digest)
	req, err := http.NewRequest("DELETE", deleteUrl, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, password)
	_, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	// remove local image
	dockerRmiCmd := fmt.Sprintf("docker rmi %s/%s:%s", url, image, tag)
	cmd = exec.Command("bash", "-c", dockerRmiCmd)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return out, err
	}
	return nil, nil
}
