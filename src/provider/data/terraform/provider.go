package terraform

import (
	"dacrane/pdk"
	"dacrane/utils"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclparse"
)

var TerraformData = pdk.Data{
	Get: func(parameter any, meta pdk.ProviderMeta) (any, error) {
		parameters := parameter.(map[string]any)
		resource := parameters["resource"].(string)
		name := parameters["name"].(string)
		argument := parameters["argument"].(map[string]any)

		mainTf := map[string]any{
			"data": map[string]any{
				resource: map[string]any{
					name: argument,
				},
			},
		}

		byteData, err := json.Marshal(mainTf)
		if err != nil {
			fmt.Println("Error marshalling to JSON:", err)
			return nil, nil
		}

		parser := hclparse.NewParser()

		// Parse the JSON string to obtain an HCL file object.
		_, diags := parser.ParseJSON(byteData, "config0.json")

		if diags.HasErrors() {
			diags.Errs()
			return nil, nil
		}

		if err != nil {
			fmt.Println("Error saving HCL file:", err)
		}

		filename := "main.tf.json"
			dir := meta.CustomStateDir
			filePath := filepath.Join(dir, filename)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("Error creating directory: %v", err)
			}
		}

		if err := os.WriteFile(filePath, byteData, 0644); err != nil {
			return nil, fmt.Errorf("Error writing JSON file: %v", err)
		}

		if err := ApplyTerraform(filePath); err != nil {
			return nil, err
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
			return r.(map[string]any)["mode"] == "managed" &&
				r.(map[string]any)["type"] == resourceType &&
				r.(map[string]any)["name"] == resourceName
		})

		instances := resource.(map[string]any)["instances"]
		instance := instances.([]any)[0]
		attributes := instance.(map[string]any)["attributes"]
		return attributes.(map[string]any), nil
	},
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
