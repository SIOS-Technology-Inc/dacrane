package core

import (
	"dacrane/pdk"
	"dacrane/utils"
	"fmt"
	"os"
	"path/filepath"

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

func (config ProjectConfig) HasInstance(name string) bool {
	return slices.Contains(config.InstanceNames, name)
}

func (config ProjectConfig) GetInstance(name string) Instance {
	instanceDir := config.InstanceDir(name)
	return Instance{
		Name:   name,
		Module: loadModule(instanceDir),
		State:  loadState(instanceDir),
	}
}

func (config ProjectConfig) GetStates() map[string]any {
	ret := map[string]any{}
	for _, instanceName := range config.InstanceNames {
		instance := config.GetInstance(instanceName)
		ret[instanceName] = instance.State
	}
	return ret
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
	modules []Module,
	writesInstance bool,
) map[string]any {
	module := utils.Find(modules, func(m Module) bool {
		return m.Name == moduleName
	})
	err := utils.Validate(module.Parameter, argument)
	if err != nil {
		panic(err)
	}
	var instance Instance
	if config.HasInstance(instanceName) {
		instance = config.GetInstance(instanceName)
	} else {
		state := map[string]any{
			"parameter": argument,
			"modules":   map[string]any{},
		}
		instance = Instance{
			Name:   instanceName,
			Module: module,
			State:  state,
		}
	}

	if writesInstance {
		config.UpsertInstance(instance)
	}
	moduleCalls := module.TopologicalSortedModuleCalls()
	for _, moduleCall := range moduleCalls {
		fmt.Printf("[%s (%s)] Evaluating...\n", moduleCall.Name, moduleCall.Module)
		evaluatedModuleCall := moduleCall.Evaluate(instance.State)
		fmt.Printf("[%s (%s)] Evaluated.\n", moduleCall.Name, moduleCall.Module)
		if evaluatedModuleCall == nil {
			fmt.Printf("[%s (%s)] Skipped.\n", moduleCall.Name, moduleCall.Module)
			continue
		}
		pluginModule, isPluginModules := providers[evaluatedModuleCall.Module]

		if isPluginModules {
			previous := instance.State["modules"].(map[string]any)[evaluatedModuleCall.Name]
			meta := pdk.ProviderMeta{CustomStateDir: filepath.Join(projectConfigDir, config.InstancesDir, instanceName, moduleCall.Name, "custom_state")}
			fmt.Printf("[%s (%s)] Applying...\n", moduleCall.Name, moduleCall.Module)
			ret, err := pluginModule.Apply(evaluatedModuleCall.Argument, previous, meta)
			if err != nil {
				panic(err)
			}
			fmt.Printf("[%s (%s)] Applied.\n", moduleCall.Name, moduleCall.Module)
			instance.State["modules"].(map[string]any)[evaluatedModuleCall.Name] = ret
		} else {
			localState := config.Apply(instanceName+"/"+moduleCall.Name, evaluatedModuleCall.Module, evaluatedModuleCall.Argument, modules, false)
			instance.State["modules"].(map[string]any)[evaluatedModuleCall.Name] = localState
		}
		if writesInstance {
			config.UpsertInstance(instance)
		}
	}
	return instance.State
}

func (config ProjectConfig) Destroy(
	instanceName string,
	modules []Module,
) {
	instance := config.GetInstance(instanceName)

	sortedModuleCalls := utils.Reverse(instance.Module.TopologicalSortedModuleCalls())
	for _, moduleCall := range sortedModuleCalls {
		state := instance.State["modules"].(map[string]any)[moduleCall.Name]
		if state == nil {
			fmt.Printf("[%s (%s)] Skipped.\n", moduleCall.Name, moduleCall.Module)
			continue
		}
		pluginModule, isPluginModules := providers[moduleCall.Module]

		if isPluginModules {
			if pluginModule.Destroy == nil {
				fmt.Printf("[%s (%s)] Skipped.\n", moduleCall.Name, moduleCall.Module)
			} else {
				fmt.Printf("[%s (%s)] Destroying...\n", moduleCall.Name, moduleCall.Module)
				filepath.Join()
				meta := pdk.ProviderMeta{CustomStateDir: filepath.Join(projectConfigDir, config.InstancesDir, instanceName, moduleCall.Name, "custom_state")}
				err := pluginModule.Destroy(state, meta)
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s (%s)] Destroyed.\n", moduleCall.Name, moduleCall.Module)
			}
		} else {
			module := utils.Find(modules, func(m Module) bool {
				return m.Name == moduleCall.Module
			})
			config.destroy(instanceName+"/"+moduleCall.Name, module, state.(map[string]any), modules)
		}
		delete(instance.State["modules"].(map[string]any), moduleCall.Name)
		config.UpsertInstance(instance)
	}
	config.DeleteInstance(instanceName)
}

func (config ProjectConfig) destroy(
	instanceName string,
	module Module,
	moduleState map[string]any,
	modules []Module,
) {
	sortedModuleCalls := utils.Reverse(module.TopologicalSortedModuleCalls())
	for _, moduleCall := range sortedModuleCalls {
		state := moduleState["modules"].(map[string]any)[moduleCall.Name]
		if state == nil {
			fmt.Printf("[%s (%s)] Skipped.\n", moduleCall.Name, moduleCall.Module)
			continue
		}
		pluginModule, isPluginModules := providers[moduleCall.Module]

		if isPluginModules {
			if pluginModule.Destroy == nil {
				fmt.Printf("[%s (%s)] Skipped.\n", moduleCall.Name, moduleCall.Module)
			} else {
				fmt.Printf("[%s (%s)] Destroying...\n", moduleCall.Name, moduleCall.Module)
				filepath.Join()
				meta := pdk.ProviderMeta{CustomStateDir: filepath.Join(projectConfigDir, config.InstancesDir, instanceName, moduleCall.Name, "custom_state")}
				err := pluginModule.Destroy(state, meta)
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s (%s)] Destroyed.\n", moduleCall.Name, moduleCall.Module)
			}
		} else {
			module := utils.Find(modules, func(m Module) bool {
				return m.Name == moduleCall.Module
			})
			config.destroy(instanceName+"/"+moduleCall.Name, module, state.(map[string]any), modules)
		}
	}
}
