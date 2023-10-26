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

type TerraformResourceProvider struct{}

type ProviderConfig struct {
	User     string `hcl:"user"`
	Password string `hcl:"password"`
}

type ResourceConfig struct {
	A int `hcl:"a"`
	B int `hcl:"b"`
}

var ctx = context.Background()

func (p TerraformResourceProvider) Create(parameters map[string]interface{}) (map[string]interface{}, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// Provider
	providerBlock := rootBody.AppendNewBlock("provider", []string{"bar"})
	providerBody := providerBlock.Body()
	providerConfig := &ProviderConfig{
		User:     "Alice", // ここどうやって抽象化する？
		Password: "abc123", // ここどうやって抽象化する？
	}
	providerBody.SetAttributeRaw("user", hclwrite.TokensForValue(cty.StringVal(providerConfig.User)))
	providerBody.SetAttributeRaw("password", hclwrite.TokensForValue(cty.StringVal(providerConfig.Password)))

	// Resource
	resourceBlock := rootBody.AppendNewBlock("resource", []string{"baz", "quz"})
	resourceBody := resourceBlock.Body()
	resourceConfig := &ResourceConfig{
		A: 1,
		B: 2,
	}
	resourceBody.SetAttributeRaw("a", hclwrite.TokensForValue(cty.NumberIntVal(int64(resourceConfig.A))))
	resourceBody.SetAttributeRaw("b", hclwrite.TokensForValue(cty.NumberIntVal(int64(resourceConfig.B))))

	// Output to a file
	instanceName := "your_instance_name"  // インスタンス名を適切に設定してください
	localModuleName := "your_module_name" // モジュール名を適切に設定してください
	filename := "your_filename.tf"        // ファイル名を適切に設定してください
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
	
	// Terraformを実行
	if err := p.ApplyTerraform(filePath); err != nil {
		return nil, fmt.Errorf("failed to apply terraform: %w", err)
	}
	
	return nil, nil
}

func (TerraformResourceProvider) ApplyTerraform(filePath string) error {
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


func (TerraformResourceProvider) Delete(parameters map[string]interface{}) error {
	return nil
}
