package environment

import "os"

type EnvironmentDataProvider struct{}

func (EnvironmentDataProvider) Get(parameters map[string]any) (map[string]any, error) {
	data := map[string]any{}
	for key, name := range parameters {
		name := name.(string)
		value, exists := os.LookupEnv(name)
		if exists {
			data[key] = value
		} else {
			data[key] = nil
		}
	}
	return data, nil
}
