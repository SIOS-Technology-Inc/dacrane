package core

import (
	"dacrane/utils"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ContextConfig struct {
	CurrentContextName string    `yaml:"current"`
	Contexts           []Context `yaml:"contexts"`
}

type Context struct {
	Name string `yaml:"name"`
}

var contextConfigDir = ".dacrane"
var contextConfigFilePath = fmt.Sprintf("%s/context.yaml", contextConfigDir)

func LoadContextConfig() ContextConfig {
	data, err := os.ReadFile(contextConfigFilePath)
	if err != nil {
		panic(err)
	}
	var config ContextConfig
	yaml.Unmarshal(data, &config)
	return config
}

func (config ContextConfig) Init() {
	err := os.Mkdir(contextConfigDir, 0755)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(contextConfigFilePath, config.GenerateYaml(), 0644)
	if err != nil {
		panic(err)
	}

	for _, c := range config.Contexts {
		c.Init()
	}
}

func (config *ContextConfig) Add(context Context) {
	config.Contexts = append(config.Contexts, context)
	context.Init()
	config.save()
}

func (config *ContextConfig) Delete(name string) {
	context := config.GetContext(name)
	context.Delete()
	config.Contexts = utils.Filter(config.Contexts, func(c Context) bool {
		return c.Name != name
	})
	config.save()
}

func (config *ContextConfig) Switch(name string) {
	config.CurrentContextName = name
	config.save()
}

func (config ContextConfig) PrettyList() string {
	s := ""
	for _, c := range config.Contexts {
		if config.IsCurrent(c) {
			s = s + fmt.Sprintf("* %s\n", c.Name)
		} else {
			s = s + fmt.Sprintf("  %s\n", c.Name)
		}
	}
	return s
}

func (config ContextConfig) save() {
	data := config.GenerateYaml()
	os.WriteFile(contextConfigFilePath, data, 0644)
}

func (config ContextConfig) IsCurrent(context Context) bool {
	return config.CurrentContextName == context.Name
}

func (context Context) Init() {
	err := os.Mkdir(context.Dir(), 0755)
	if err != nil {
		panic(err)
	}
}

func (context Context) Delete() {
	err := os.RemoveAll(context.Dir())
	if err != nil {
		panic(err)
	}
}

func (context Context) Dir() string {
	return fmt.Sprintf("%s/%s", contextConfigDir, context.Name)
}

func (context Context) StateFilePath() string {
	return fmt.Sprintf("%s/state.yaml", context.Dir())
}

func (context Context) EnvFilePath() string {
	return fmt.Sprintf("%s/env.yaml", context.Dir())
}

func NewDefaultContextConfig() ContextConfig {
	return ContextConfig{
		CurrentContextName: "default",
		Contexts: []Context{
			{
				Name: "default",
			},
		},
	}
}

func (config ContextConfig) CurrentContext() Context {
	return utils.Find(config.Contexts, func(c Context) bool {
		return c.Name == config.CurrentContextName
	})
}

func (config ContextConfig) GetContext(name string) Context {
	context := utils.Find(config.Contexts, func(c Context) bool {
		return c.Name == name
	})
	if context.Name == "" {
		panic(fmt.Sprintf("%s context is not found.", name))
	}
	return context
}

func (config ContextConfig) GenerateYaml() []byte {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	return data
}

func (context Context) ReadState() []byte {
	data, err := os.ReadFile(context.StateFilePath())
	if err != nil {
		panic(err)
	}
	return data
}

func (context Context) WriteState(data []byte) {
	err := os.WriteFile(context.StateFilePath(), data, 0644)
	if err != nil {
		panic(err)
	}
}

func (context Context) ReadEnv() []byte {
	data, err := os.ReadFile(context.EnvFilePath())
	if err != nil {
		panic(err)
	}
	return data
}

func (context Context) WriteEnv(data []byte) {
	err := os.WriteFile(context.EnvFilePath(), data, 0644)
	if err != nil {
		panic(err)
	}
}
