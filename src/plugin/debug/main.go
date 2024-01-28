package main

import (
	"dacrane/pdk"
	"encoding/json"
)

func main() {
	pdk.ExecPluginJob(pdk.Plugin{
		Config: pdk.NewDefaultPluginConfig(),
		Resources: pdk.MapToFunc(map[string]pdk.Resource{
			"print": PrintResource,
			"dummy": DummyResource,
		}),
	})
}

var PrintResource = pdk.Resource{
	Create: func(parameter any, meta pdk.PluginMeta) (any, error) {
		v, err := json.MarshalIndent(parameter, "", "  ")
		if err != nil {
			return nil, err
		}
		meta.Log(string(v))

		return parameter, nil
	},
	Update: func(current, previous any, meta pdk.PluginMeta) (any, error) {
		v, err := json.MarshalIndent(current, "", "  ")
		if err != nil {
			return nil, err
		}
		meta.Log(string(v))

		return current, nil
	},
	Delete: func(parameter any, meta pdk.PluginMeta) error {
		v, err := json.MarshalIndent(parameter, "", "  ")
		if err != nil {
			return err
		}

		meta.Log(string(v))

		return nil
	},
}

var DummyResource = pdk.Resource{
	Create: func(parameter any, meta pdk.PluginMeta) (any, error) {
		return parameter, nil
	},
	Update: func(current, previous any, meta pdk.PluginMeta) (any, error) {
		return current, nil
	},
	Delete: func(parameter any, meta pdk.PluginMeta) error {
		return nil
	},
}
