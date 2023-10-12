package file

import (
	"context"
	"os"
)
type FileProvider struct{}
var ctx = context.Background()

func (FileProvider) Create(parameters map[string]any) (map[string]any, error) {
	statesYaml := []byte{}
	contents := parameters["contents"].(string)
	filename := parameters["filename"].(string)

	statesYaml = append(statesYaml, []byte(contents)...)
	e := os.WriteFile(filename, statesYaml, 0644)
	if e != nil {
    return nil, e
	}

 return nil, nil
}

func (fp FileProvider) Delete(parameters map[string]interface{}) error {

	filename := parameters["filename"].(string)
	err := os.Remove(filename)
	if err != nil {
			return err
	}
	return nil
}