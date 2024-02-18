package repository

import (
	"os"
	"path/filepath"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/utils"

	"gopkg.in/yaml.v3"
)

type DocumentRepository struct {
	documentFile string
	document     map[string]any
}

// TODO Specify from entry point
var documentFile = ".dacrane/instance.yaml"

func InitDocumentRepositoryFile() {
	dir := filepath.Dir(documentFile)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
	bytes, err := yaml.Marshal(map[string]any{})
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(documentFile, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func LoadDocumentRepository() DocumentRepository {
	repository := DocumentRepository{
		documentFile: documentFile,
	}
	repository.read()
	return repository
}

func (repository *DocumentRepository) read() {
	bytes, err := os.ReadFile(repository.documentFile)
	if err != nil {
		panic(err)
	}
	var doc map[string]any
	err = yaml.Unmarshal(bytes, &doc)
	if err != nil {
		panic(err)
	}
	repository.document = doc
}

func (repository DocumentRepository) write() {
	bytes, err := yaml.Marshal(repository.document)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(repository.documentFile, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func (repository DocumentRepository) Exists(id string) bool {
	_, exists := repository.document[id]
	return exists
}

func (repository DocumentRepository) Find(id string) any {
	return repository.document[id]
}

func (repository *DocumentRepository) Upsert(id string, document any) {
	repository.document[id] = document
	repository.write()
	repository.read()
}

func (repository *DocumentRepository) Delete(id string) {
	delete(repository.document, id)
	repository.write()
	repository.read()
}

func (repository DocumentRepository) Ids() []string {
	return utils.Keys(repository.document)
}

func (repository DocumentRepository) Document() map[string]any {
	return repository.document
}
