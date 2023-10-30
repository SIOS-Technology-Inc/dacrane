package environment

import "os"

type EnvironmentDataProvider struct{}

func (EnvironmentDataProvider) Get(parameter any) (any, error) {
	params := parameter.(map[string]any)
	data := map[string]any{}
	for key, name := range params {
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
