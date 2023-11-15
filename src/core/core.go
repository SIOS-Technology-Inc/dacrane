package core

import (
	"dacrane/pdk"
	"dacrane/provider/data/environment"
	terraform_data "dacrane/provider/data/terraform"
	docker_container "dacrane/provider/resource/docker-container"
	docker_local_image "dacrane/provider/resource/docker-local-image"
	docker_remote_image "dacrane/provider/resource/docker-remote-image"
	file "dacrane/provider/resource/file"
	terraform_resource "dacrane/provider/resource/terraform"
)

var providers = map[string](pdk.Provider){
	"data/environment":             environment.EnvironmentDataProvider,
	"data/terraform":               terraform_data.TerraformDataProvider,
	"resource/docker-container":    docker_container.DockerContainerResourceProvider,
	"resource/docker-local-image":  docker_local_image.DockerLocalImageResourceProvider,
	"resource/docker-remote-image": docker_remote_image.DockerRemoteImageProvider,
	"resource/file":                file.FileResourceProvider,
	"resource/terraform":           terraform_resource.TerraformResourceProvider,
}
