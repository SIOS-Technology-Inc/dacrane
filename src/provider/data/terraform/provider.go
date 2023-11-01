package terraform

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type TerraformDataProvider struct{}

type ProviderConfig struct {
	User     string `hcl:"user"`
	Password string `hcl:"password"`
}

type DataConfig struct {
	A int `hcl:"a"`
	B int `hcl:"b"`
}

var ctx = context.Background()

func (p TerraformDataProvider) Get(parameters map[string]any) (map[string]any, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// Setting Provider
	providerName, ok := parameters["provider"].(string)
	if !ok {
		return nil, fmt.Errorf("provider name is required and must be a string")
	}
	providerBlock := rootBody.AppendNewBlock("provider", []string{providerName})
	providerBody := providerBlock.Body()
	if configs, ok := parameters["configurations"].(map[string]interface{}); ok {
		for k, v := range configs {
			providerBody.SetAttributeValue(k, cty.StringVal(fmt.Sprintf("%v", v)))
		}
	}

	// Setting Resource
	resourceType, ok := parameters["resource"].(string)
	if !ok {
		return nil, fmt.Errorf("resource type is required and must be a string")
	}
	resourceName, ok := parameters["name"].(string)
	if !ok {
		return nil, fmt.Errorf("resource name is required and must be a string")
	}
	resourceBlock := rootBody.AppendNewBlock("data", []string{resourceType, resourceName})
	resourceBody := resourceBlock.Body()
	if args, ok := parameters["argument"].(map[string]interface{}); ok {
		for k, v := range args {
			switch v := v.(type) {
			case string:
				resourceBody.SetAttributeValue(k, cty.StringVal(v))
			case int, int64, float64:
				resourceBody.SetAttributeValue(k, cty.NumberFloatVal(float64(v.(int))))
			default:
				return nil, fmt.Errorf("unsupported type for argument: %v", v)
			}
		}
	}

	// write file
	instanceName := "your_instance_name"
	localModuleName := "your_module_name"
	filename := "your_filename.tf"
	dir := filepath.Join(".dacrane", "instances", instanceName, "custom_states", localModuleName)
	filePath := filepath.Join(dir, filename)

	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}
	
	// Write the file
	if err := os.WriteFile(filePath, f.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}
	
	fmt.Printf("HCL written to %s\n", filePath)
	
	// Terraform exec
	if err := p.ApplyTerraform(filePath); err != nil {
		return nil, fmt.Errorf("failed to apply terraform: %w", err)
	}
	
	return nil, nil
}


func (TerraformDataProvider) ApplyTerraform(filePath string) error {
	// Terraform init
	dir := filepath.Dir(filePath)
	
	initCmd := exec.Command("terraform", "init")
	initCmd.Dir = dir
	if output, err := initCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run terraform init: %s, %s", err, output)
	}

	// Terraform apply
	applyCmd := exec.Command("terraform", "apply", "-auto-approve", filePath)
	if output, err := applyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run terraform apply: %s, %s", err, output)
	}

	fmt.Println("Terraform apply complete")
	return nil
}