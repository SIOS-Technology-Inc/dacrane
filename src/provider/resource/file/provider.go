package file

import (
	"os"
)

type FileProvider struct{}

func (FileProvider) Create(parameter any) (any, error) {
	params := parameter.(map[string]any)
	statesYaml := []byte{}
	contents := params["contents"].(string)
	filename := params["filename"].(string)

	statesYaml = append(statesYaml, []byte(contents)...)
	e := os.WriteFile(filename, statesYaml, 0644)
	if e != nil {
		return nil, e
	}

	return nil, nil
}

func (fp FileProvider) Delete(parameter any) error {
	params := parameter.(map[string]any)
	filename := params["filename"].(string)
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}
