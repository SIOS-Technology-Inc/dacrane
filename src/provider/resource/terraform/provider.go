package terraform

import (
	"fmt"

	"context"
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

func (TerraformResourceProvider) Create(parameters map[string]any) (map[string]any, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// Provider
	providerBlock := rootBody.AppendNewBlock("provider", []string{"bar"})
	providerBody := providerBlock.Body()
	providerConfig := &ProviderConfig{
		User:     "Alice",
		Password: "abc123",
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

	// Output HCL
	fmt.Println(string(f.Bytes()))

 return nil, nil
}

func (tf TerraformResourceProvider) Delete(parameters map[string]interface{}) error {

	return nil
}