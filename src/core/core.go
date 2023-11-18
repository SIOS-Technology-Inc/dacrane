package core

import (
	"dacrane/core/module"
	"dacrane/provider/data/environment"
	terraform_data "dacrane/provider/data/terraform"
	docker_container "dacrane/provider/resource/docker-container"
	docker_local_image "dacrane/provider/resource/docker-local-image"
	docker_remote_image "dacrane/provider/resource/docker-remote-image"
	file "dacrane/provider/resource/file"
	terraform_resource "dacrane/provider/resource/terraform"
)

var Providers = []module.Provider{
	module.NewDataProvider("data/environment", environment.EnvironmentData),
	module.NewDataProvider("data/terraform", terraform_data.TerraformData),
	module.NewResourceProvider("resource/docker-container", docker_container.DockerContainerResource),
	module.NewResourceProvider("resource/docker-local-image", docker_local_image.DockerLocalImageResource),
	module.NewResourceProvider("resource/docker-remote-image", docker_remote_image.DockerRemoteImage),
	module.NewResourceProvider("resource/file", file.FileResource),
	module.NewResourceProvider("resource/terraform", terraform_resource.TerraformResource),
}
