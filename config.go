package main

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
)

func readResourceFile(filename string) (contents []byte) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return contents
}

func parseResourceFile(filedata []byte) (resources map[string]Resource) {
	resources = make(map[string]Resource)
	if err := yaml.Unmarshal(filedata, resources); err != nil {
		log.Fatalf("error: %v", err)
	}

	for key, resource := range resources {
		resource.Name = key
	}

	return resources
}

func loadResourceFile(filename string) map[string]Resource {
	return parseResourceFile(readResourceFile(filename))
}
