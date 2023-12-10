package pdk

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type Log func(string)

type PluginMeta struct {
	CustomStateDir string
	Log            Log
}

func NewPluginMeta(customStateDir string) PluginMeta {
	return PluginMeta{
		CustomStateDir: customStateDir,
		Log: func(s string) {
			fmt.Fprintln(os.Stderr, s)
		},
	}
}

type PluginConfig struct {
	DockerHost *string
	WorkingDir string
}

func NewDefaultPluginConfig() PluginConfig {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return PluginConfig{
		DockerHost: nil,
		WorkingDir: currentDir,
	}
}

type Resource struct {
	Create func(parameter any, meta PluginMeta) (any, error)
	Update func(current any, previous any, meta PluginMeta) (any, error)
	Delete func(parameter any, meta PluginMeta) error
}

type Data struct {
	Get func(parameters any, meta PluginMeta) (any, error)
}

type Plugin struct {
	Config    PluginConfig
	Resources func(string) (Resource, bool)
	Data      func(string) (Data, bool)
}

func MapToFunc[T any](m map[string]T) func(string) (T, bool) {
	return func(name string) (T, bool) {
		e, ok := m[name]
		return e, ok
	}
}

func ExecPluginJob(plugin Plugin) {
	rawArg := os.Args[1]
	decodeArg, err := base64.StdEncoding.DecodeString(rawArg)
	if err != nil {
		panic(err)
	}
	var arg map[string]any
	err = json.Unmarshal([]byte(decodeArg), &arg)
	if err != nil {
		panic(err)
	}
	kind := arg["kind"].(string)
	name := arg["name"].(string)
	operation := arg["operation"].(string)
	customStateDir := arg["custom_state_dir"].(string)
	arguments := arg["arguments"].([]any)

	if err != nil {
		panic(err)
	}

	meta := NewPluginMeta(customStateDir)

	output, err := exec(plugin, meta, kind, name, operation, arguments)
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

func exec(plugin Plugin, meta PluginMeta, kind string, name string, operation string, arguments []any) (any, error) {
	switch kind {
	case "resource":
		resource, ok := plugin.Resources(name)
		if !ok {
			return nil, fmt.Errorf("not find resource: %s", name)
		}
		return execResource(resource, meta, operation, arguments)
	case "data":
		data, ok := plugin.Data(name)
		if !ok {
			return nil, fmt.Errorf("not find data: %s", name)
		}
		return execData(data, meta, operation, arguments)
	case "preflight":
		// It returns the directory path that the plugin wants to be mounted volumes by dacrane core.
		// This process is executed before "data" or "resource" is called.
		return map[string]any{
			"working_dir": plugin.Config.WorkingDir,
			"docker_host": plugin.Config.DockerHost,
		}, nil
	default:
		return nil, fmt.Errorf("not supported kind: %s", kind)
	}
}

func execResource(resource Resource, meta PluginMeta, operation string, arguments []any) (any, error) {
	switch operation {
	case "create":
		return resource.Create(arguments[0], meta)
	case "update":
		if resource.Update == nil {
			resource.Update = func(current any, previous any, meta PluginMeta) (any, error) {
				err := resource.Delete(previous, meta)
				if err != nil {
					return nil, err
				}
				return resource.Create(current, meta)
			}
		}
		return resource.Update(arguments[0], arguments[1], meta)
	case "delete":
		err := resource.Delete(arguments[0], meta)
		return nil, err
	default:
		return nil, fmt.Errorf("not supported operation: %s", operation)
	}
}

func execData(data Data, meta PluginMeta, operation string, arguments []any) (any, error) {
	switch operation {
	case "get":
		return data.Get(arguments, meta)
	default:
		return nil, fmt.Errorf("not supported operation: %s", operation)
	}
}
