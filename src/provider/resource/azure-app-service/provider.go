package azureappservice

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type AzureAppServiceResourceProvider struct{}

var ctx = context.Background()

func (AzureAppServiceResourceProvider) Create(parameters map[string]any, credentials map[string]any) error {
	subscriptionId := credentials["subscription_id"].(string)
	tenantId := credentials["tenant_id"].(string)
	clientId := credentials["client_id"].(string)
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	name := parameters["name"].(string)
	resourceGroupName := parameters["resource_group_name"].(string)
	location := parameters["location"].(string)
	appServicePlanId := parameters["app_service_plan_id"].(string)
	siteConfig := parameters["site_config"].(map[string]any)
	appSettings := map[string]*string{}
	for k, v := range parameters["app_settings"].(map[string]any) {
		s := v.(string)
		appSettings[k] = &s
	}

	linuxFxVersion := siteConfig["linux_fx_version"].(string)
	println(linuxFxVersion)

	client := web.NewAppsClient(subscriptionId)

	cred := auth.NewUsernamePasswordConfig(username, password, clientId, tenantId)
	auth, err := cred.Authorizer()
	if err != nil {
		return nil
	}

	client.Authorizer = auth

	siteEnvelope := web.Site{
		Location: &location,
		SiteProperties: &web.SiteProperties{
			ServerFarmID: &appServicePlanId,
			SiteConfig: &web.SiteConfig{
				LinuxFxVersion: &linuxFxVersion,
			},
		},
	}

	_, err = client.CreateOrUpdate(ctx, resourceGroupName, name, siteEnvelope)
	if err != nil {
		return err
	}

	settings := web.StringDictionary{
		Properties: appSettings,
	}

	if _, err := client.UpdateApplicationSettings(ctx, resourceGroupName, name, settings); err != nil {
		return fmt.Errorf("updating Application Settings for App Service %q: %+v", name, err)
	}

	return nil
}

func (AzureAppServiceResourceProvider) Delete(parameters map[string]any, credentials map[string]any) error {
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
