package docker

import (
	"dacrane/pdk"
	"dacrane/utils"
	"fmt"
	"net/http"
	"strings"
)

var DockerRemoteImageProvider = pdk.NewResourceModule(pdk.Resource{
	Create: func(parameter any, _ pdk.ProviderMeta) (any, error) {
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
			_, err := utils.RunOnBash(cmd)
			if err != nil {
				return nil, err
			}
		}

		return params, nil
	},
	Delete: func(parameter any, _ pdk.ProviderMeta) error {
		params := parameter.(map[string]any)
		image := params["image"].(string)
		tag := params["tag"].(string)

		remote := params["remote"].(map[string](any))
		url := remote["url"].(string)
		user := remote["user"].(string)
		password := remote["password"].(string)

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
		res, err := utils.RequestHttp(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		dockerRmiCmd := fmt.Sprintf("docker rmi %s/%s:%s", url, image, tag)
		_, err = utils.RunOnBash(dockerRmiCmd)
		if err != nil {
			return err
		}

		return nil
	},
})
