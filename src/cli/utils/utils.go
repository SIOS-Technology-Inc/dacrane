package utils

import (
	"errors"

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

func FillDefault(schema any, document any) any {
	if schema == nil {
		return nil
	}
	defaultValue := schema.(map[string]any)["default"]
	if document == nil {
		return defaultValue
	}

	switch schema.(map[string]any)["type"] {
	case "object":
		properties, hasProperties := schema.(map[string]any)["properties"].(map[string]any)
		if !hasProperties {
			return document
		}
		result := map[string]any{}
		for key, propSchema := range properties {
			filledDoc := FillDefault(propSchema, document.(map[string]any)[key])
			result[key] = filledDoc
		}
		return result
	case "array":
		itemsSchema, hasItems := schema.(map[string]any)["items"]
		if !hasItems {
			return document
		}
		result := []any{}
		for _, item := range document.([]any) {
			filledDoc := FillDefault(itemsSchema, item)
			result = append(result, filledDoc)
		}
		return result
	default:
		return document
	}
}
