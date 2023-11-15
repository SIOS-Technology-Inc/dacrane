package terraform

import (
	"dacrane/pdk"
	"dacrane/utils"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

var TerraformDataProvider = pdk.NewDataProvider(pdk.Data{
	Get: func(parameter any, meta pdk.ProviderMeta) (any, error) {
		parameters := parameter.(map[string]any)
		f := hclwrite.NewEmptyFile()
		rootBody := f.Body()

		// Setting Provider
		if provider, ok := parameters["provider"].(string); ok {
			providerBlock := rootBody.AppendNewBlock("provider", []string{provider})
			providerBody := providerBlock.Body()
			if configs, ok := parameters["configurations"].(map[string]interface{}); ok {
				for k, v := range configs {
					writeHCL(providerBody, k, v)
				}
			}
		} else {
			return nil, fmt.Errorf("provider name is required and must be a string")
		}

		// Setting Resource
		resourceType, resourceName := "", "main"
		if resType, ok := parameters["resource"].(string); ok {
			resourceType = resType
		} else {
			return nil, fmt.Errorf("resource type is required and must be a string")
		}

		resourceBlock := rootBody.AppendNewBlock("data", []string{resourceType, resourceName})
		resourceBody := resourceBlock.Body()
		if args, ok := parameters["argument"].(map[string]interface{}); ok {
			for k, v := range args {
				writeHCL(resourceBody, k, v)
			}
		}

		// write file
		filename := "main.tf"
		dir := meta.CustomStateDir
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
		if err := ApplyTerraform(filePath); err != nil {
			return nil, fmt.Errorf("failed to apply terraform: %w", err)
		}

		// Get Terraform State
		bytes, err := os.ReadFile(dir + "/terraform.tfstate")
		if err != nil {
			return nil, err
		}

		var state map[string]any
		err = json.Unmarshal(bytes, &state)
		if err != nil {
			return nil, err
		}

		resource := utils.Find(state["resources"].([]any), func(r any) bool {
			return r.(map[string]any)["mode"] == "data" &&
				r.(map[string]any)["type"] == resourceType &&
				r.(map[string]any)["name"] == resourceName
		})

		instances := resource.(map[string]any)["instances"]
		instance := instances.([]any)[0]
		attributes := instance.(map[string]any)["attributes"]
		return attributes.(map[string]any), nil
	},
})

func writeHCL(body *hclwrite.Body, key string, value interface{}) {
	switch v := value.(type) {
	case map[string]interface{}:
		block := body.AppendNewBlock(key, nil)
		blockBody := block.Body()
		for k, val := range v {
			writeHCL(blockBody, k, val)
		}
	case string:
		body.SetAttributeValue(key, cty.StringVal(v))
	case []interface{}:
		values := make([]cty.Value, len(v))
		for i, val := range v {
			values[i] = cty.StringVal(val.(string))
		}
		body.SetAttributeValue(key, cty.ListVal(values))
	default:
		fmt.Printf("Unsupported type: %T\n", v)
	}
}

func ApplyTerraform(filePath string) error {
	// Terraform init
	dir := filepath.Dir(filePath)

	initCmd := exec.Command("terraform", "init")
	initCmd.Dir = dir
	if output, err := initCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run terraform init: %s, %s", err, output)
	}

	// Terraform apply
	applyCmd := exec.Command("terraform", "apply", "-auto-approve")
	applyCmd.Dir = dir
	if output, err := applyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run terraform apply: %s, %s", err, output)
	}

	fmt.Println("Terraform apply complete")
	return nil
}
