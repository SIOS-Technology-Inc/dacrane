package main

import (
	"dacrane/pdk"
	"os"
)

func main() {
	pdk.ExecPluginJob(pdk.Plugin{
		Config: pdk.NewDefaultPluginConfig(),
		Resources: pdk.MapToFunc(map[string]pdk.Resource{
			"file": FileResource,
		}),
		Data: pdk.MapToFunc(map[string]pdk.Data{}),
	})
}

var FileResource = pdk.Resource{
	Create: func(parameter any, _ pdk.PluginMeta) (any, error) {
		params := parameter.(map[string]any)
		contents := params["contents"].(string)
		filename := params["filename"].(string)

		e := os.WriteFile(filename, []byte(contents), 0644)
		if e != nil {
			return nil, e
		}

		return parameter, nil
	},
	Delete: func(parameter any, _ pdk.PluginMeta) error {
		params := parameter.(map[string]any)
		filename := params["filename"].(string)
		err := os.Remove(filename)
		if err != nil {
			return err
		}
		return nil
	},
}
