package core

import (
	"bytes"
	"dacrane/provider/artifact/docker"

	"gopkg.in/yaml.v3"
)

func ParseCode(codeBytes []byte) ([]Code, error) {
	r := bytes.NewReader(codeBytes)
	dec := yaml.NewDecoder(r)

	var codes []Code
	var code Code
	for dec.Decode(&code) == nil {
		codes = append(codes, code)
	}

	return codes, nil
}

type Code struct {
	Kind       string         `yaml:"kind"`
	Name       string         `yaml:"name"`
	Provider   string         `yaml:"provider"`
	Parameters map[string]any `yaml:"parameters"`
}

type ArtifactProvider interface {
	Build(parameters map[string]any) error
	Publish(parameters map[string]any) error
	Unpublish(parameters map[string]any) error
	SearchVersions(map[string]any) error
}

var artifactProviders = map[string](ArtifactProvider){
	"docker": docker.DockerArtifactProvider{},
}

func FindArtifactProvider(providerName string) ArtifactProvider {
	return artifactProviders[providerName]
}
