package azureappserviceplan

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
)

type AzureAppServicePlanResourceProvider struct{}

var ctx = context.Background()

func (AzureAppServicePlanResourceProvider) Create(parameters map[string]any) (map[string]any, error) {
	credentials := parameters["credentials"].(map[string]any)
	subscriptionId := credentials["subscription_id"].(string)
	tenantId := credentials["tenant_id"].(string)
	clientId := credentials["client_id"].(string)
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	name := parameters["name"].(string)
	resourceGroupName := parameters["resource_group_name"].(string)
	location := parameters["location"].(string)
	kind := parameters["kind"].(string)
	sku := parameters["sku"].(map[string]any)
	sku_tier := sku["tier"].(string)
	sku_name := sku["name"].(string)

	cred, err := azidentity.NewUsernamePasswordCredential(tenantId, clientId, username, password, nil)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armappservice.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	client := clientFactory.NewPlansClient()

	_, err = client.BeginCreateOrUpdate(ctx, resourceGroupName, name, armappservice.Plan{
		Location: &location,
		Kind:     &kind,
		SKU: &armappservice.SKUDescription{
			Name: &sku_name,
			Tier: &sku_tier,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	return parameters, nil
}

func (AzureAppServicePlanResourceProvider) Delete(parameters map[string]any) error {
	credentials := parameters["credentials"].(map[string]any)
	subscriptionId := credentials["subscription_id"].(string)
	tenantId := credentials["tenant_id"].(string)
	clientId := credentials["client_id"].(string)
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	name := parameters["name"].(string)
	resourceGroupName := parameters["resource_group_name"].(string)

	cred, err := azidentity.NewUsernamePasswordCredential(tenantId, clientId, username, password, nil)
	if err != nil {
		return err
	}

	clientFactory, err := armappservice.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return err
	}

	client := clientFactory.NewPlansClient()

	_, err = client.Delete(ctx, resourceGroupName, name, nil)
	if err != nil {
		return err
	}

	return nil
}
