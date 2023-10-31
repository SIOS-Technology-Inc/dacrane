package core

import (
	"dacrane/pdk"
	"dacrane/provider/data/environment"
	docker_container "dacrane/provider/resource/docker-container"
	docker_local_image "dacrane/provider/resource/docker-local-image"
	docker_remote_image "dacrane/provider/resource/docker-remote-image"
	file "dacrane/provider/resource/file"
)

var pluginModules = map[string](pdk.Module){
	"data/environment":             environment.EnvironmentDataModule,
	"resource/docker-container":    docker_container.DockerContainerResourceModule,
	"resource/docker-local-image":  docker_local_image.DockerLocalImageResourceModule,
	"resource/docker-remote-image": docker_remote_image.DockerRemoteImageProvider,
	"resource/file":                file.FileProvider,
}
