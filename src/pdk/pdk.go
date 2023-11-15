package pdk

type ProviderMeta struct {
	CustomStateDir string
}

type Provider struct {
	Apply   func(current any, previous any, meta ProviderMeta) (any, error)
	Destroy func(current any, meta ProviderMeta) error
}

func NewResourceProvider(resource Resource) Provider {
	if resource.Update == nil {
		resource.Update = func(current any, previous any, meta ProviderMeta) (any, error) {
			err := resource.Delete(previous, meta)
			if err != nil {
				return nil, err
			}
			return resource.Create(current, meta)
		}
	}
	return Provider{
		Apply: func(current any, previous any, meta ProviderMeta) (any, error) {
			if previous == nil {
				return resource.Create(current, meta)
			} else {
				return resource.Update(current, previous, meta)
			}
		},
		Destroy: func(current any, meta ProviderMeta) error {
			return resource.Delete(current, meta)
		},
	}
}

func NewDataProvider(data Data) Provider {
	return Provider{
		Apply: func(current any, previous any, meta ProviderMeta) (any, error) {
			return data.Get(current, meta)
		},
		Destroy: nil,
	}
}

type Resource struct {
	Create func(parameter any, meta ProviderMeta) (any, error)
	Update func(current any, previous any, meta ProviderMeta) (any, error)
	Delete func(parameter any, meta ProviderMeta) error
}

type Data struct {
	Get func(parameters any, meta ProviderMeta) (any, error)
}
