package environment

import (
	"dacrane/pdk"
	"os"
)

var EnvironmentData = pdk.Data{
	Get: func(parameter any, _ pdk.ProviderMeta) (any, error) {
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
	},
}
