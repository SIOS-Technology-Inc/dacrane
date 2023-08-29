package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

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
