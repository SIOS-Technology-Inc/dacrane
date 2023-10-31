package pdk

type Module struct {
	Apply   func(current any, previous any) (any, error)
	Destroy func(current any) error
}

func NewResourceModule(resource Resource) Module {
	if resource.Update == nil {
		resource.Update = func(current any, previous any) (any, error) {
			err := resource.Delete(previous)
			if err != nil {
				return nil, err
			}
			return resource.Create(current)
		}
	}
	return Module{
		Apply: func(current any, previous any) (any, error) {
			if previous == nil {
				return resource.Create(current)
			} else {
				return resource.Update(current, previous)
			}
		},
		Destroy: func(current any) error {
			return resource.Delete(current)
		},
	}
}

func NewDataModule(data Data) Module {
	return Module{
		Apply: func(current any, previous any) (any, error) {
			return data.Get(current)
		},
		Destroy: nil,
	}
}

type Resource struct {
	Create func(parameter any) (any, error)
	Update func(current any, previous any) (any, error)
	Delete func(parameter any) error
}

type Data struct {
	Get func(parameters any) (any, error)
}
