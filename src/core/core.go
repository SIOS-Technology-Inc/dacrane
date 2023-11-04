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

var pluginModules = map[string](pdk.Provider){
	"data/environment":             environment.EnvironmentDataModule,
	"data/terraform":               terraform_data.TerraformDataModule,
	"resource/docker-container":    docker_container.DockerContainerResourceModule,
	"resource/docker-local-image":  docker_local_image.DockerLocalImageResourceModule,
	"resource/docker-remote-image": docker_remote_image.DockerRemoteImageProvider,
	"resource/file":                file.FileProvider,
	"resource/terraform":           terraform_resource.TerraformResourceModule,
}
