package core

import (
	"dacrane/utils"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
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

func (config *ProjectConfig) UpsertInstance(instance Instance) {
	if !slices.Contains(config.InstanceNames, instance.Name) {
		config.InstanceNames = append(config.InstanceNames, instance.Name)
		config.save()
	}
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

func (config ProjectConfig) Apply(
	instanceName string,
	moduleName string,
	argument any,
	dependencies map[string]string,
	modules []Module,
) map[string]any {
	module := utils.Find(modules, func(m Module) bool {
		return m.Name == moduleName
	})
	state := map[string]any{
		"parameter": argument,
		"module":    map[string]any{},
	}
	instance := Instance{
		Name:   instanceName,
		Module: module,
		State:  state,
	}
	config.UpsertInstance(instance)
	moduleCalls := module.TopologicalSortedModuleCalls()
	for _, moduleCall := range moduleCalls {
		fmt.Printf("[%s (%s)] Evaluating...\n", moduleCall.Name, moduleCall.Module)
		evaluatedModuleCall := moduleCall.Evaluate(state)
		fmt.Printf("[%s (%s)] Evaluated\n", moduleCall.Name, moduleCall.Module)

		modulePaths := strings.Split(evaluatedModuleCall.Module, "/")
		kind := modulePaths[0]

		switch kind {
		case "resource":
			name := modulePaths[1]
			resourceProvider := FindResourceProvider(name)
			fmt.Printf("[%s (%s)] Crating...\n", moduleCall.Name, moduleCall.Module)
			ret, err := resourceProvider.Create(evaluatedModuleCall.Argument.(map[string]any))
			if err != nil {
				panic(err)
			}
			fmt.Printf("[%s (%s)] Created.\n", moduleCall.Name, moduleCall.Module)
			state["module"].(map[string]any)[evaluatedModuleCall.Name] = ret
		case "artifact":
			name := modulePaths[1]
			artifactProvider := FindArtifactProvider(name)
			fmt.Printf("[%s (%s)] Building...\n", moduleCall.Name, moduleCall.Module)
			err := artifactProvider.Build(evaluatedModuleCall.Argument.(map[string]any))
			if err != nil {
				panic(err)
			}
			fmt.Printf("[%s (%s)] Built.\n", moduleCall.Name, moduleCall.Module)
			fmt.Printf("[%s (%s)] Publishing...\n", moduleCall.Name, moduleCall.Module)
			ret, err := artifactProvider.Publish(evaluatedModuleCall.Argument.(map[string]any))
			if err != nil {
				panic(err)
			}
			fmt.Printf("[%s (%s)] Published.\n", moduleCall.Name, moduleCall.Module)
			state["module"].(map[string]any)[evaluatedModuleCall.Name] = ret
		case "data":
			name := modulePaths[1]
			dataProvider := FindDataProvider(name)
			fmt.Printf("[%s (%s)] Reading...\n", moduleCall.Name, moduleCall.Module)
			ret, err := dataProvider.Get(evaluatedModuleCall.Argument.(map[string]any))
			if err != nil {
				panic(err)
			}
			fmt.Printf("[%s (%s)] Read.\n", moduleCall.Name, moduleCall.Module)
			state["module"].(map[string]any)[evaluatedModuleCall.Name] = ret
		default:
			localState := config.Apply(instanceName, kind, moduleCall.Argument, map[string]string{}, modules)
			state["module"].(map[string]any)[evaluatedModuleCall.Name] = localState["module"]
		}
		instance.State = state
		config.UpsertInstance(instance)
	}
	return state
}

func (config ProjectConfig) Destroy(
	instanceName string,
) {
	instance := config.GetInstance(instanceName)

	sortedModuleCalls := utils.Reverse(instance.Module.TopologicalSortedModuleCalls())
	for _, moduleCall := range sortedModuleCalls {
		state := instance.State["module"].(map[string]any)[moduleCall.Name]
		if state == nil {
			fmt.Printf("[%s (%s)] Skipped.\n", moduleCall.Name, moduleCall.Module)
			continue
		}

		modulePaths := strings.Split(moduleCall.Module, "/")
		kind := modulePaths[0]

		switch kind {
		case "resource":
			name := modulePaths[1]
			resourceProvider := FindResourceProvider(name)
			fmt.Printf("[%s (%s)] Deleting...\n", moduleCall.Name, moduleCall.Module)
			err := resourceProvider.Delete(state.(map[string]any))
			if err != nil {
				panic(err)
			}
			fmt.Printf("[%s (%s)] Deleted.\n", moduleCall.Name, moduleCall.Module)
		case "artifact":
			name := modulePaths[1]
			artifactProvider := FindArtifactProvider(name)
			fmt.Printf("[%s (%s)] Unpublish...\n", moduleCall.Name, moduleCall.Module)
			err := artifactProvider.Unpublish(state.(map[string]any))
			if err != nil {
				panic(err)
			}
			fmt.Printf("[%s (%s)] Unpublished.\n", moduleCall.Name, moduleCall.Module)
		case "data":

		}
		delete(instance.State["module"].(map[string]any), "")
		config.UpsertInstance(instance)
	}
	config.DeleteInstance(instanceName)
}
