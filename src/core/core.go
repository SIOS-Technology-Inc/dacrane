package core

import (
	"dacrane/provider/data/environment"
	azure_app_service "dacrane/provider/resource/azure-app-service"
	azure_app_service_plan "dacrane/provider/resource/azure-app-service-plan"
	azure_container_registry "dacrane/provider/resource/azure-container-registry"
	azure_resource_group "dacrane/provider/resource/azure-resource-group"
	docker_container "dacrane/provider/resource/docker-container"
	docker_local_image "dacrane/provider/resource/docker-local-image"
	docker_remote_image "dacrane/provider/resource/docker-remote-image"
)

type ResourceProvider interface {
	Create(parameters map[string]any) (map[string]any, error)
	Delete(parameters map[string]any) error
}

type DataProvider interface {
	Get(parameters map[string]any) (map[string]any, error)
}

var resourceProviders = map[string](ResourceProvider){
	"azure-resource-group":     azure_resource_group.AzureResourceGroupResourceProvider{},
	"azure-app-service-plan":   azure_app_service_plan.AzureAppServicePlanResourceProvider{},
	"azure-app-service":        azure_app_service.AzureAppServiceResourceProvider{},
	"azure-container-registry": azure_container_registry.AzureContainerRegistryResourceProvider{},
	"docker-container":         docker_container.DockerResourceProvider{},
	"docker_local_image":       docker_local_image.DockerArtifactProvider{},
	"docker_remote_image":      docker_remote_image.DockerArtifactProvider{},
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
