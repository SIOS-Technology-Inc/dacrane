package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/xeipuuv/gojsonschema"
)

func Contains[T any](array []T, f func(T) bool) bool {
	for _, value := range array {
		if f(value) {
			return true
		}
	}
	return false
}

func Find[T any](array []T, f func(T) bool) (result T) {
	for _, value := range array {
		if f(value) {
			return value
		}
	}
	return
}

func Filter[T any](array []T, f func(T) bool) (result []T) {
	for _, value := range array {
		if f(value) {
			result = append(result, value)
		}
	}
	return
}

func Map[T, T2 any](array []T, f func(T) T2) (result []T2) {
	for _, value := range array {
		result = append(result, f(value))
	}
	return
}

func Reverse[T any](array []T) []T {
	for i := 0; i < len(array)/2; i++ {
		array[i], array[len(array)-i-1] = array[len(array)-i-1], array[i]
	}
	return array
}

func Keys[T comparable](array map[T]any) (result []T) {
	for k := range array {
		result = append(result, k)
	}
	return
}

func Values[T comparable, T2 any](array map[T]T2) (result []T2) {
	for _, v := range array {
		result = append(result, v)
	}
	return
}

func RunOnBash(script string) ([]byte, error) {
	fmt.Printf("> %s\n", script)
	cmd := exec.Command("bash", "-c", script)
	writer := new(bytes.Buffer)
	cmd.Stdout = io.MultiWriter(os.Stdout, writer)
	cmd.Stderr = io.MultiWriter(os.Stderr, writer)
	err := cmd.Run()
	return writer.Bytes(), err
}

func RequestHttp(req *http.Request) (*http.Response, error) {
	fmt.Printf("> %s %s\n", req.Method, req.URL)
	res, err := http.DefaultClient.Do(req)
	fmt.Printf("> %s\n", res.Status)
	return res, err
}

func Validate(schema any, document any) error {
	if schema == nil {
		return nil
	}
	schemaLoader := gojsonschema.NewGoLoader(schema)
	documentLoader := gojsonschema.NewGoLoader(document)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err)
	}

	if result.Valid() {
		return nil
	} else {
		err := result.Errors()[0]
		return errors.New(err.String())
	}
}

func FillDefault(schema any, document any) (any, error) {
	defaultValue := schema.(map[string]any)["default"]
	if document == nil {
		return defaultValue, nil
	}

	switch schema.(map[string]any)["type"] {
	case "object":
		properties, hasProperties := schema.(map[string]any)["properties"].(map[string]any)
		if !hasProperties {
			return document, nil
		}
		result := map[string]any{}
		for key, propSchema := range properties {
			filledDoc, err := FillDefault(propSchema, document.(map[string]any)[key])
			if err != nil {
				return nil, err
			}
			result[key] = filledDoc
		}
		return result, nil
	case "array":
		itemsSchema, hasItems := schema.(map[string]any)["items"]
		if !hasItems {
			return document, nil
		}
		result := []any{}
		for _, item := range document.([]any) {
			filledDoc, err := FillDefault(itemsSchema, item)
			if err != nil {
				return nil, err
			}
			result = append(result, filledDoc)
		}
		return result, nil
	default:
		return document, nil
	}
}
