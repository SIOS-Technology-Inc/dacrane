package file

import (
	"os"
)

type FileProvider struct{}

func (FileProvider) Create(parameter any) (any, error) {
	params := parameter.(map[string]any)
	contents := params["contents"].(string)
	filename := params["filename"].(string)

	e := os.WriteFile(filename, []byte(contents), 0644)
	if e != nil {
		return nil, e
	}

	return nil, nil
}

func (provider FileProvider) Update(current any, previous any) (any, error) {
	err := provider.Delete(previous)
	if err != nil {
		return nil, err
	}
	return provider.Create(current)
}

func (FileProvider) Delete(parameter any) error {
	params := parameter.(map[string]any)
	filename := params["filename"].(string)
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}
