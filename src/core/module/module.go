package module

import (
	"bytes"
	"dacrane/core/evaluator"
	"dacrane/core/repository"
	"dacrane/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"gopkg.in/yaml.v3"
)

type Module struct {
	Name         string       `yaml:"name"`
	Import       []string     `yaml:"import"`
	Parameter    any          `yaml:"parameter"`
	Dependencies []Dependency `yaml:"dependencies"`
	ModuleCalls  []ModuleCall `yaml:"modules"`
}

type Dependency struct {
	Name   string `yaml:"name"`
	Module string `yaml:"module"`
}

type ModuleCall struct {
	Name         string            `yaml:"name"`
	DependsOn    []string          `yaml:"depends_on"`
	Module       string            `yaml:"module"`
	Argument     any               `yaml:"argument"`
	Dependencies map[string]string `yaml:"dependencies"`
	If           any               `yaml:"if"`
}

func (module Module) Apply(
	instanceAddress string,
	argument any,
	instances *repository.DocumentRepository,
	importedModule []Module,
	importedProvider []Provider,
) {
	// Validation Argument
	err := utils.Validate(module.Parameter, argument)
	if err != nil {
		panic(err)
	}

	// Create or get the instance
	var instance moduleInstance
	if instances.Exists(instanceAddress) {
		document := instances.Find(instanceAddress)
		instance = NewInstanceFromDocument(document).(moduleInstance)
	} else {
		instance = NewModuleInstance(module, instanceAddress, argument)
		instances.Upsert(instanceAddress, instance)
	}

	// Import external modules
	// TODO scope handling
	for _, urlOrFilepath := range module.Import {
		importedModule = append(importedModule, Import(urlOrFilepath)...)
	}

	moduleCalls := module.TopologicalSortedModuleCalls()
	for _, moduleCall := range moduleCalls {
		fmt.Printf("[%s (%s)] Evaluating...\n", instanceAddress, moduleCall.Module)
		evaluatedModuleCall := moduleCall.Evaluate(instance.ToState(*instances).(map[string]any))
		fmt.Printf("[%s (%s)] Evaluated.\n", instanceAddress, moduleCall.Module)
		if evaluatedModuleCall == nil {
			fmt.Printf("[%s (%s)] Skipped.\n", instanceAddress, moduleCall.Module)
			continue
		}

		childRelAddr := evaluatedModuleCall.Name
		childAbsAddr := instanceAddress + "." + evaluatedModuleCall.Name

		providerExists := utils.Contains(importedProvider, func(provider Provider) bool {
			return provider.Name == evaluatedModuleCall.Module
		})
		moduleExists := utils.Contains(importedModule, func(module Module) bool {
			return module.Name == evaluatedModuleCall.Module
		})

		if providerExists {
			provider := utils.Find(importedProvider, func(provider Provider) bool {
				return provider.Name == evaluatedModuleCall.Module
			})
			provider.Apply(childAbsAddr, evaluatedModuleCall.Argument, instances)
		} else if moduleExists {
			childModule := utils.Find(importedModule, func(module Module) bool {
				return module.Name == evaluatedModuleCall.Module
			})
			childModule.Apply(childAbsAddr, evaluatedModuleCall.Argument, instances, importedModule, importedProvider)
		} else {
			panic(fmt.Sprintf("undefined module or provider: %s", evaluatedModuleCall.Module))
		}
		instance.Instances = append(instance.Instances, childRelAddr)
		instances.Upsert(instanceAddress, instance)
	}
}

func Import(urlOrFilepath string) []Module {
	_, err := url.ParseRequestURI(urlOrFilepath)
	if err == nil {
		return ImportFromUrl(urlOrFilepath)
	} else {
		return ImportFromFilepath(urlOrFilepath)
	}
}

func ImportFromUrl(url string) []Module {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return ParseModules(body)
}

func ImportFromFilepath(filepath string) []Module {
	codeBytes, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return ParseModules(codeBytes)
}

func ParseModules(codeBytes []byte) []Module {
	r := bytes.NewReader(codeBytes)
	dec := yaml.NewDecoder(r)

	modules := []Module{}
	for {
		var module Module
		if err := dec.Decode(&module); err != nil {
			if err.Error() == "EOF" {
				break
			}
			panic(err)
		}
		modules = append(modules, module)
	}

	return modules
}

func (module Module) FindModuleCall(name string) ModuleCall {
	return utils.Find(module.ModuleCalls, func(mc ModuleCall) bool {
		return mc.Name == name
	})
}

func (module Module) TopologicalSortedModuleCalls() []ModuleCall {
	g := simple.NewDirectedGraph()

	idToName := map[int64]string{}
	nodes := map[string]graph.Node{}
	for _, moduleCall := range module.ModuleCalls {
		node := g.NewNode()
		nodes[moduleCall.Name] = node
		g.AddNode(node)
		idToName[node.ID()] = moduleCall.Name
	}

	for _, mc := range module.ModuleCalls {
		ds := mc.Dependency()
		for _, d := range ds {
			g.SetEdge(g.NewEdge(nodes[d], nodes[mc.Name]))
		}
	}

	sorted, err := topo.Sort(g)
	if err != nil {
		panic(err)
	}
	return utils.Map(sorted, func(node graph.Node) ModuleCall {
		return module.FindModuleCall(idToName[node.ID()])
	})
}

// returns dependency module name
func (mc ModuleCall) Dependency() []string {
	return append(mc.ExplicitDependency(), mc.ImplicitDependency()...)
}

func (mc ModuleCall) ExplicitDependency() []string {
	return mc.DependsOn
}

func (mc ModuleCall) ImplicitDependency() []string {
	var paths []string
	paths = append(paths, references(mc.Name)...)
	paths = append(paths, references(mc.Module)...)
	paths = append(paths, references(mc.Argument)...)
	return paths
}

func (mc ModuleCall) Evaluate(data map[string]any) *ModuleCall {

	mapMc := mc.toMap()

	evaluated := Evaluate(mapMc, data)

	if evaluated == nil {
		return nil
	}

	return toModuleCall(evaluated.(map[string]any))
}

func (mc ModuleCall) toMap() map[string]any {
	if mc.If == nil {
		mc.If = true
	}
	return map[string]any{
		"name":       mc.Name,
		"depends_on": mc.DependsOn,
		"module":     mc.Module,
		"argument":   mc.Argument,
		"if":         mc.If,
	}
}

func toModuleCall(mc map[string]any) *ModuleCall {
	var dependsOn []string
	if mc["depends_on"] == nil {
		dependsOn = []string{}
	} else {
		dependsOn = mc["depends_on"].([]string)
	}

	return &ModuleCall{
		Name:      mc["name"].(string),
		DependsOn: dependsOn,
		Module:    mc["module"].(string),
		Argument:  mc["argument"],
	}
}

func Evaluate(prop any, data map[string]any) any {
	switch prop := prop.(type) {
	case string:
		single := isSingleExprString(prop)
		if single {
			r, e := regexp.Compile(`^\$\{\{(.*?)\}\}$`)
			if e != nil {
				panic(e)
			}
			exprStr := r.FindStringSubmatch(prop)[1]
			expr := evaluator.Parse(exprStr)
			return evaluator.Evaluate(expr, data)
		} else {
			return expandExpr(prop, data)
		}
	case map[string]any:
		prop, exists := evalIfProp(prop, data)
		if !exists {
			return nil
		}
		output := map[string]any{}
		for k, v := range prop {
			output[k] = Evaluate(v, data)
		}
		return output
	case []any:
		output := []any{}
		for _, v := range prop {
			output = append(output, Evaluate(v, data))
		}
		return output
	default:
		return prop
	}
}

func expandExpr(prop string, data map[string]any) string {
	r, e := regexp.Compile(`\$\{\{(.*?)\}\}`)
	if e != nil {
		panic(e)
	}
	return r.ReplaceAllStringFunc(prop, func(s string) string {
		exprStr := r.FindStringSubmatch(s)
		expr := evaluator.Parse(exprStr[1])
		v := evaluator.Evaluate(expr, data)
		return convertToString(v)
	})
}

func evalIfProp(prop map[string]any, data map[string]any) (map[string]any, bool) {
	if condition, ok := prop["if"]; ok {
		if !Evaluate(condition, data).(bool) {
			return nil, false
		}
	}
	delete(prop, "if")
	return prop, true
}

func isSingleExprString(s string) bool {
	r, e := regexp.Compile(`^\$\{\{.*?\}\}$`)
	if e != nil {
		panic(e)
	}
	return r.MatchString(s)
}

func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func references(raw any) []string {
	switch raw := raw.(type) {
	case map[string]any:
		var paths []string
		for _, v := range raw {
			paths = append(paths, references(v)...)
		}
		return paths
	case []any:
		var paths []string
		for _, v := range raw {
			paths = append(paths, references(v)...)
		}
		return paths
	case string:
		r, e := regexp.Compile(`\$\{\{(.*?)\}\}`)
		if e != nil {
			panic(e)
		}
		res := r.FindAllStringSubmatch(raw, -1)
		var paths []string
		for _, exprStr := range res {
			expr := evaluator.Parse(exprStr[1])
			paths = append(paths, evaluator.CollectReferences(expr)...)
		}
		return paths
	default:
		return []string{}
	}
}

func (module Module) GenerateYaml() []byte {
	data, err := yaml.Marshal(module)
	if err != nil {
		panic(err)
	}
	return data
}

func PrettyInstanceList(instances repository.DocumentRepository) string {
	s := ""
	for _, address := range instances.Ids() {
		document := instances.Find(address)
		instance := NewInstanceFromDocument(document)
		if !strings.Contains(address, ".") {
			s = s + fmt.Sprintf("%s (%s)\n", address, instance.(moduleInstance).Module.Name)
		}
	}
	return s
}
