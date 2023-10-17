package core

import (
	"dacrane/utils"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ProjectConfig struct {
	InstancesDir  string   `yaml:"instances_dir"`
	InstanceNames []string `yaml:"instances"`
}

type Instance struct {
	Name   string
	Module Module
	State  map[string]any
}

var projectConfigDir = ".dacrane"
var configFilePath = fmt.Sprintf("%s/config.yaml", projectConfigDir)

func StateFilePath(instanceDir string) string {
	return fmt.Sprintf("%s/state.yaml", instanceDir)
}

func ModuleFilePath(instanceDir string) string {
	return fmt.Sprintf("%s/module.yaml", instanceDir)
}

func NewProjectConfig() ProjectConfig {
	return ProjectConfig{
		InstancesDir:  "instances",
		InstanceNames: []string{},
	}
}

func LoadProjectConfig() ProjectConfig {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	var config ProjectConfig
	yaml.Unmarshal(data, &config)
	return config
}

func (config ProjectConfig) Init() {
	err := os.Mkdir(projectConfigDir, 0755)
	if err != nil {
		panic(err)
	}

	instanceDir := fmt.Sprintf("%s/%s", projectConfigDir, config.InstancesDir)
	err = os.Mkdir(instanceDir, 0755)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(configFilePath, config.GenerateYaml(), 0644)
	if err != nil {
		panic(err)
	}
}

func (config *ProjectConfig) CreateInstance(instance Instance) {
	config.InstanceNames = append(config.InstanceNames, instance.Name)
	instance.save(config.InstanceDir(instance.Name))
	config.save()
}

func (config ProjectConfig) UpdateInstance(instance Instance) {
	instance.save(config.InstanceDir(instance.Name))
}

func (config *ProjectConfig) DeleteInstance(name string) {
	instanceDir := config.InstanceDir(name)
	err := os.RemoveAll(instanceDir)
	if err != nil {
		panic(err)
	}
	config.InstanceNames = utils.Filter(config.InstanceNames, func(n string) bool {
		return n != name
	})
	config.save()
}

func (config ProjectConfig) PrettyList() string {
	s := ""
	for _, name := range config.InstanceNames {
		instance := config.GetInstance(name)
		s = s + fmt.Sprintf("%s (%s)\n", name, instance.Module.Name)
	}
	return s
}

func (config ProjectConfig) InstanceDir(instanceName string) string {
	return fmt.Sprintf("%s/%s/%s", projectConfigDir, config.InstancesDir, instanceName)
}

func (config ProjectConfig) save() {
	data := config.GenerateYaml()

	os.WriteFile(configFilePath, data, 0644)
}

func (config ProjectConfig) GetInstance(name string) Instance {
	instanceDir := config.InstanceDir(name)
	return Instance{
		Name:   name,
		Module: loadModule(instanceDir),
		State:  loadState(instanceDir),
	}
}

func (config ProjectConfig) GenerateYaml() []byte {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	return data
}

func loadState(instanceDir string) map[string]any {
	data, err := os.ReadFile(StateFilePath(instanceDir))
	if err != nil {
		panic(err)
	}
	var state map[string]any
	err = yaml.Unmarshal(data, &state)
	if err != nil {
		panic(err)
	}

	return state
}

func loadModule(instanceDir string) Module {
	data, err := os.ReadFile(ModuleFilePath(instanceDir))
	if err != nil {
		panic(err)
	}
	var module Module
	err = yaml.Unmarshal(data, &module)
	if err != nil {
		panic(err)
	}

	return module
}

func (instance Instance) WriteState(instanceDir string) {
	data, err := yaml.Marshal(instance.State)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(StateFilePath(instanceDir), data, 0644)
	if err != nil {
		panic(err)
	}
}

func (instance Instance) WriteModule(instanceDir string) {
	data, err := yaml.Marshal(instance.Module)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(ModuleFilePath(instanceDir), data, 0644)
	if err != nil {
		panic(err)
	}
}

func (instance Instance) save(instanceDir string) {
	err := os.MkdirAll(instanceDir, 0755)
	if err != nil {
		panic(err)
	}
	instance.WriteModule(instanceDir)
	instance.WriteState(instanceDir)
}
