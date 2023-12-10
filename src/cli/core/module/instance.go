package module

import (
	"dacrane/cli/core/repository"
	"dacrane/cli/utils"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Instance interface {
	ToState(instances repository.DocumentRepository) any
	Destroy(instanceAddress string, instances *repository.DocumentRepository)
}

type moduleInstance struct {
	Type      string   `yaml:"type"`
	Module    Module   `yaml:"module"`
	Argument  any      `yaml:"argument"`
	Address   string   `yaml:"address"`
	Instances []string `yaml:"instances"`
}

type providerInstance struct {
	Type           string `yaml:"type"`
	Provider       string `yaml:"provider"`
	CustomStateDir string `yaml:"custom_state_dir"`
	Argument       any    `yaml:"argument"`
	Output         any    `yaml:"output"`
}

func NewModuleInstance(module Module, address string, argument any) moduleInstance {
	return moduleInstance{
		Type:      "module",
		Module:    module,
		Address:   address,
		Argument:  argument,
		Instances: []string{},
	}
}

func NewPluginInstance(provider string, customStateDir string, argument any, output any) providerInstance {
	return providerInstance{
		Type:           "provider",
		Provider:       provider,
		CustomStateDir: customStateDir,
		Argument:       argument,
		Output:         output,
	}
}

func NewInstanceFromDocument(document any) Instance {
	t := document.(map[string]any)["type"]
	switch t {
	case "module":
		bytes, err := yaml.Marshal(document)
		if err != nil {
			panic(err)
		}
		var instance moduleInstance
		yaml.Unmarshal(bytes, &instance)
		return instance
	case "provider":
		bytes, err := yaml.Marshal(document)
		if err != nil {
			panic(err)
		}
		var instance providerInstance
		yaml.Unmarshal(bytes, &instance)
		return instance
	default:
		panic(fmt.Sprintf("unknown instance type: %s", t))
	}
}

func (instance moduleInstance) ToState(instances repository.DocumentRepository) any {
	state := map[string]any{
		"parameter": instance.Argument,
		"modules":   map[string]any{},
	}
	for _, address := range instance.Instances {
		childAbsAddr := instance.Address + "." + address
		doc := instances.Find(childAbsAddr)
		child := NewInstanceFromDocument(doc)
		state["modules"].(map[string]any)[address] = child.ToState(instances)
	}
	return state
}

func (instance moduleInstance) Destroy(
	instanceAddress string,
	instances *repository.DocumentRepository,
) {
	sortedModuleCalls := utils.Reverse(instance.Module.TopologicalSortedModuleCalls())
	for _, moduleCall := range sortedModuleCalls {
		childAbsAddr := instanceAddress + "." + moduleCall.Name
		childRelAddr := moduleCall.Name
		if !instances.Exists(childAbsAddr) {
			fmt.Printf("[%s (%s)] Skipped. %s is not exist.\n",
				instanceAddress, moduleCall.Module, childAbsAddr)
			continue
		}

		document := instances.Find(childAbsAddr)
		child := NewInstanceFromDocument(document)
		child.Destroy(childAbsAddr, instances)
		customStatePath := filepath.Join(".dacrane/custom_state", childAbsAddr)
		err := os.RemoveAll(customStatePath)
		if err != nil {
			panic(err)
		}

		instance.Instances = utils.Filter(instance.Instances, func(instance string) bool {
			return instance != childRelAddr
		})
		instances.Upsert(instanceAddress, instance)
	}
	instances.Delete(instanceAddress)
}

func (instance providerInstance) ToState(_ repository.DocumentRepository) any {
	return instance.Output
}

func (instance providerInstance) Destroy(instanceAddress string, instances *repository.DocumentRepository) {
	plugin := NewPlugin(instance.Provider)
	if plugin.Destroy == nil {
		fmt.Printf("[%s (%s)] Skipped. Deletion is not needed.\n", instanceAddress, instance.Provider)
	}
	plugin.Destroy(instanceAddress, instances)
}
