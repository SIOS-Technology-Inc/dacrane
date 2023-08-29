package docker

import (
	"dacrane/utils"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DockerArtifactProvider struct{}

func (DockerArtifactProvider) Build(params map[string]any) error {
	dockerfile := params["dockerfile"].(string)
	image := params["image"].(string)
	tag := params["tag"].(string)

	dockerCmd := fmt.Sprintf("docker build -t %s:%s -f %s .", image, tag, dockerfile)
	_, err := utils.RunOnBash(dockerCmd)
	return err
}

func (DockerArtifactProvider) Publish(params map[string]any) error {
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

	for _, cmd := range cmds {
		_, err := utils.RunOnBash(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (DockerArtifactProvider) Unpublish(params map[string]any) error {
	image := params["image"].(string)
	tag := params["tag"].(string)
	repository := params["repository"].(map[string](any))
	url := repository["url"].(string)
	user := repository["user"].(string)
	password := repository["password"].(string)

	// remove registry image
	dockerDigestCmd := fmt.Sprintf("docker images %s/%s --format {{.Digest}}", url, image)
	out, err := utils.RunOnBash(dockerDigestCmd)
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
	_, err = utils.RequestHttp(req)
	if err != nil {
		return err
	}

	// remove local image
	dockerRmiCmd := fmt.Sprintf("docker rmi %s/%s:%s", url, image, tag)
	_, err = utils.RunOnBash(dockerRmiCmd)
	if err != nil {
		return err
	}
	dockerRmiCmd = fmt.Sprintf("docker rmi %s:%s", image, tag)
	_, err = utils.RunOnBash(dockerRmiCmd)
	if err != nil {
		return err
	}
	return nil
}

func (DockerArtifactProvider) SearchVersions(params map[string]any) error {
	image := params["image"].(string)
	repository := params["repository"].(map[string](any))
	url := repository["url"].(string)
	user := repository["user"].(string)
	password := repository["password"].(string)

	localVersionsCmd := fmt.Sprintf("docker images %s --format {{.Tag}}", image)

	_, err := utils.RunOnBash(localVersionsCmd)
	if err != nil {
		return err
	}

	listUrl := fmt.Sprintf("https://%s/v2/%s/tags/list", url, image)
	req, err := http.NewRequest("GET", listUrl, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	res, err := utils.RequestHttp(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
