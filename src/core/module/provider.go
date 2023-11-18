package module

import (
	"dacrane/core/repository"
	"dacrane/pdk"
	"fmt"
)

type Provider struct {
	Name  string
	Apply func(
		instanceAddress string,
		argument any,
		instances *repository.DocumentRepository,
	)
	Destroy func(
		instanceAddress string,
		instances *repository.DocumentRepository,
	)
}

func NewResourceProvider(providerName string, resource pdk.Resource) Provider {
	if resource.Update == nil {
		resource.Update = func(current any, previous any, meta pdk.ProviderMeta) (any, error) {
			err := resource.Delete(previous, meta)
			if err != nil {
				return nil, err
			}
			return resource.Create(current, meta)
		}
	}
	return Provider{
		Name: providerName,
		Apply: func(
			instanceAddress string,
			argument any,
			instances *repository.DocumentRepository,
		) {
			if instances.Exists(instanceAddress) {
				fmt.Printf("[%s (%s)] Updating...\n", instanceAddress, providerName)

				document := instances.Find(instanceAddress)
				instance := NewInstanceFromDocument(document).(providerInstance)

				meta := pdk.ProviderMeta{
					CustomStateDir: instance.CustomStateDir,
				}
				output, err := resource.Update(argument, instance.ToState(*instances), meta)
				if err != nil {
					panic(err)
				}
				instance.Output = output
				instances.Upsert(instanceAddress, instance)
				fmt.Printf("[%s (%s)] Updated.\n", instanceAddress, providerName)
			} else {
				fmt.Printf("[%s (%s)] Creating...\n", instanceAddress, providerName)
				// TODO Specify from entry point
				meta := pdk.ProviderMeta{CustomStateDir: fmt.Sprintf(".dacrane/custom_state/%s", instanceAddress)}
				output, err := resource.Create(argument, meta)
				if err != nil {
					panic(err)
				}
				instance := NewProviderInstance(providerName, meta.CustomStateDir, argument, output)
				instances.Upsert(instanceAddress, instance)
				fmt.Printf("[%s (%s)] Created.\n", instanceAddress, providerName)
			}
		},
		Destroy: func(instanceAddress string, instances *repository.DocumentRepository) {
			if !instances.Exists(instanceAddress) {
				fmt.Printf("[%s (%s)] Skipped. %s is not exist.\n",
					instanceAddress, providerName, instanceAddress)
			}
			fmt.Printf("[%s (%s)] Deleting...\n", instanceAddress, providerName)
			document := instances.Find(instanceAddress)
			instance := NewInstanceFromDocument(document).(providerInstance)

			meta := pdk.ProviderMeta{
				CustomStateDir: instance.CustomStateDir,
			}
			err := resource.Delete(instance.Output, meta)
			if err != nil {
				panic(err)
			}
			instances.Delete(instanceAddress)
			fmt.Printf("[%s (%s)] Deleted.\n", instanceAddress, providerName)
		},
	}
}

func NewDataProvider(providerName string, data pdk.Data) Provider {
	return Provider{
		Name: providerName,
		Apply: func(
			instanceAddress string,
			argument any,
			instances *repository.DocumentRepository,
		) {
			fmt.Printf("[%s (%s)] Reading...\n", instanceAddress, providerName)
			// TODO Specify from entry point
			meta := pdk.ProviderMeta{CustomStateDir: fmt.Sprintf(".dacrane/custom_state/%s", instanceAddress)}
			output, err := data.Get(argument, meta)
			if err != nil {
				panic(err)
			}
			instance := NewProviderInstance(providerName, meta.CustomStateDir, argument, output)
			instances.Upsert(instanceAddress, instance)
			fmt.Printf("[%s (%s)] Read.\n", instanceAddress, providerName)
		},
		Destroy: nil,
	}
}
