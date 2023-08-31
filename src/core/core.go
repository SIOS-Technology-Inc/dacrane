package core

import (
	"bytes"
	"dacrane/provider/artifact/docker"
	azureresourcegroup "dacrane/provider/resource/azure-resource-group"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseCode(codeBytes []byte) ([]Code, error) {
	r := bytes.NewReader([]byte(os.ExpandEnv(string(codeBytes))))
	dec := yaml.NewDecoder(r)

	var codes []Code
	var code Code
	for dec.Decode(&code) == nil {
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
	"azure-resource-group": azureresourcegroup.AzureResourceGroupArtifactProvider{},
}

func FindArtifactProvider(providerName string) ArtifactProvider {
	return artifactProviders[providerName]
}

func FindResourceProvider(providerName string) ResourceProvider {
	return resourceProviders[providerName]
}
