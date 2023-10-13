package core

import (
	"dacrane/provider/artifact/docker"
	environemnt "dacrane/provider/data/environment"
	azureappservice "dacrane/provider/resource/azure-app-service"
	azureappserviceplan "dacrane/provider/resource/azure-app-service-plan"
	azurecontainerregistry "dacrane/provider/resource/azure-container-registry"
	azureresourcegroup "dacrane/provider/resource/azure-resource-group"
	dockerresource "dacrane/provider/resource/docker"
	file "dacrane/provider/resource/file"
)

type ArtifactProvider interface {
	Build(parameters map[string]any) error
	Publish(parameters map[string]any) (map[string]any, error)
	Unpublish(parameters map[string]any) error
	SearchVersions(map[string]any) error
}

type ResourceProvider interface {
	Create(parameters map[string]any) (map[string]any, error)
	Delete(parameters map[string]any) error
}

type DataProvider interface {
	Get(parameters map[string]any) (map[string]any, error)
}

var artifactProviders = map[string](ArtifactProvider){
	"docker": docker.DockerArtifactProvider{},
}

var resourceProviders = map[string](ResourceProvider){
	"azure-resource-group":     azureresourcegroup.AzureResourceGroupResourceProvider{},
	"azure-app-service-plan":   azureappserviceplan.AzureAppServicePlanResourceProvider{},
	"azure-app-service":        azureappservice.AzureAppServiceResourceProvider{},
	"azure-container-registry": azurecontainerregistry.AzureContainerRegistryResourceProvider{},
	"docker":                   dockerresource.DockerResourceProvider{},
	"file": 										file.FileProvider{},
}

var dataProviders = map[string](DataProvider){
	"environment": environemnt.EnvironmentDataProvider{},
}

func FindArtifactProvider(providerName string) ArtifactProvider {
	return artifactProviders[providerName]
}

func FindResourceProvider(providerName string) ResourceProvider {
	return resourceProviders[providerName]
}

func FindDataProvider(providerName string) DataProvider {
	return dataProviders[providerName]
}
