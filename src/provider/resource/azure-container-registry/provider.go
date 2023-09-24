package azurecontainerregistry

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

type AzureContainerRegistryResourceProvider struct{}

var ctx = context.Background()

func (AzureContainerRegistryResourceProvider) Create(parameters map[string]any) (map[string]any, error) {
	credentials := parameters["credentials"].(map[string]any)
	subscriptionId := credentials["subscription_id"].(string)
	tenantId := credentials["tenant_id"].(string)
	clientId := credentials["client_id"].(string)
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	name := parameters["name"].(string)
	resourceGroupName := parameters["resource_group_name"].(string)
	location := parameters["location"].(string)

	client := containerregistry.NewRegistriesClient(subscriptionId)
	cred := auth.NewUsernamePasswordConfig(username, password, clientId, tenantId)
	auth, err := cred.Authorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = auth

	params := containerregistry.Registry{
		Location: &location,
		Sku: &containerregistry.Sku{
			Name: containerregistry.Basic,
		},
		RegistryProperties: &containerregistry.RegistryProperties{
			AdminUserEnabled: to.BoolPtr(true),
		},
	}

	_, err = client.Create(ctx, resourceGroupName, name, params)
	if err != nil {
		return nil, err
	}

	creds, err := client.ListCredentials(ctx, resourceGroupName, name)
	if err != nil {
		return nil, err
	}

	parameters["url"] = fmt.Sprintf("%s.azurecr.io", name)

	if creds.Username != nil {
		parameters["user"] = *creds.Username
	}
	if creds.Passwords != nil && len(*creds.Passwords) > 0 {
		parameters["password"] = *(*creds.Passwords)[0].Value
	}

	return parameters, nil
}

func (AzureContainerRegistryResourceProvider) Delete(parameters map[string]any) error {
	credentials := parameters["credentials"].(map[string]any)
	subscriptionId := credentials["subscription_id"].(string)
	tenantId := credentials["tenant_id"].(string)
	clientId := credentials["client_id"].(string)
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	name := parameters["name"].(string)
	resourceGroupName := parameters["resource_group_name"].(string)

	client := containerregistry.NewRegistriesClient(subscriptionId)
	cred := auth.NewUsernamePasswordConfig(username, password, clientId, tenantId)
	auth, err := cred.Authorizer()
	if err != nil {
		return err
	}

	client.Authorizer = auth

	_, err = client.Delete(ctx, resourceGroupName, name)
	if err != nil {
		return err
	}

	return nil
}
