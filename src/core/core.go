package core

import (
	"bytes"
	"dacrane/provider/artifact/docker"
	azureappservice "dacrane/provider/resource/azure-app-service"
	azureappserviceplan "dacrane/provider/resource/azure-app-service-plan"
	azurecontainerregistry "dacrane/provider/resource/azure-container-registry"
	azureresourcegroup "dacrane/provider/resource/azure-resource-group"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseCode(codeBytes []byte) ([]Code, error) {
	r := bytes.NewReader([]byte(os.ExpandEnv(string(codeBytes))))
	dec := yaml.NewDecoder(r)

	var codes []Code
	for {
		var code Code
		if dec.Decode(&code) != nil {
			break
		}
		codes = append(codes, code)
	}

	return codes, nil
}

type Code struct {
	Kind        string         `yaml:"kind"`
	Name        string         `yaml:"name"`
	Provider    string         `yaml:"provider"`
	Parameters  map[string]any `yaml:"parameters"`
	Credentials map[string]any `yaml:"credentials"`
}

type ArtifactProvider interface {
	Build(parameters map[string]any) error
	Publish(parameters map[string]any) error
	Unpublish(parameters map[string]any) error
	SearchVersions(map[string]any) error
}

type ResourceProvider interface {
	Create(parameters map[string]any, credentials map[string]any) error
	Delete(parameters map[string]any, credentials map[string]any) error
}

var artifactProviders = map[string](ArtifactProvider){
	"docker": docker.DockerArtifactProvider{},
}

var resourceProviders = map[string](ResourceProvider){
	"azure-resource-group":     azureresourcegroup.AzureResourceGroupResourceProvider{},
	"azure-app-service-plan":   azureappserviceplan.AzureAppServicePlanResourceProvider{},
	"azure-app-service":        azureappservice.AzureAppServiceResourceProvider{},
	"azure-container-registry": azurecontainerregistry.AzureContainerRegistryResourceProvider{},
}

func FindArtifactProvider(providerName string) ArtifactProvider {
	return artifactProviders[providerName]
}

func FindResourceProvider(providerName string) ResourceProvider {
	return resourceProviders[providerName]
}
