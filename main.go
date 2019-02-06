package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/catho/Sem-Release/semrelease"
)

func main() {

	repository := semrelease.NewRepository()
	service := semrelease.NewService(repository)
	owner := os.Getenv("owner")
	repo := os.Getenv("repository")
	owner = "eduardokenjimiura"
	repo = "job_configuration"
	service.CreateRelease(owner, repo)
	fmt.Println("finish")

}
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
