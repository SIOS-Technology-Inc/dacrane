package core

import (
	"dacrane/provider/data/environment"
	docker_container "dacrane/provider/resource/docker-container"
	docker_local_image "dacrane/provider/resource/docker-local-image"
	docker_remote_image "dacrane/provider/resource/docker-remote-image"
	file "dacrane/provider/resource/file"
)

type ResourceProvider interface {
	Create(parameter any) (any, error)
	Delete(parameter any) error
}

type DataProvider interface {
	Get(parameters any) (any, error)
}

var resourceProviders = map[string](ResourceProvider){
	"docker-container":    docker_container.DockerResourceProvider{},
	"docker-local-image":  docker_local_image.DockerArtifactProvider{},
	"docker-remote-image": docker_remote_image.DockerArtifactProvider{},
	"file":                file.FileProvider{},
}

var dataProviders = map[string](DataProvider){
	"environment": environment.EnvironmentDataProvider{},
}

func FindResourceProvider(providerName string) ResourceProvider {
	return resourceProviders[providerName]
}

func FindDataProvider(providerName string) DataProvider {
	return dataProviders[providerName]
}
