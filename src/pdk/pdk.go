package pdk

type ProviderMeta struct {
	CustomStateDir string
}

type Resource struct {
	Create func(parameter any, meta ProviderMeta) (any, error)
	Update func(current any, previous any, meta ProviderMeta) (any, error)
	Delete func(parameter any, meta ProviderMeta) error
}

type Data struct {
	Get func(parameters any, meta ProviderMeta) (any, error)
}
