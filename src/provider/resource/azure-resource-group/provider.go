package azureresourcegroup

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type DockerArtifactProvider struct{}

var ctx = context.Background()

func (DockerArtifactProvider) Create(parameters map[string]any, credentials map[string]any) error {
	subscriptionId := credentials["subscription_id"].(string)
	tenantId := credentials["tenant_id"].(string)
	clientId := credentials["client_id"].(string)
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	name := parameters["name"].(string)
	location := parameters["location"].(string)

	cred, err := azidentity.NewUsernamePasswordCredential(tenantId, clientId, username, password, nil)
	if err != nil {
		return err
	}

	rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		return err
	}

	param := armresources.ResourceGroup{
		Location: &location,
	}

	_, err = rgClient.CreateOrUpdate(ctx, name, param, nil)
	if err != nil {
		return err
	}

	return nil
}

func (DockerArtifactProvider) Delete(parameters map[string]any, credentials map[string]any) error {
	subscriptionId := credentials["subscription_id"].(string)
	tenantId := credentials["tenant_id"].(string)
	clientId := credentials["client_id"].(string)
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	name := parameters["name"].(string)

	cred, err := azidentity.NewUsernamePasswordCredential(tenantId, clientId, username, password, nil)
	if err != nil {
		return err
	}

	rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		return err
	}

	_, err = rgClient.BeginDelete(ctx, name, nil)
	if err != nil {
		return err
	}

	return nil
}
