package main

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load and parse resource file
	resourceFile := "test.yaml"
	resourceData, readErr := ioutil.ReadFile(resourceFile)
	if readErr != nil {
		log.Fatalf("error: %v", readErr)
	}

	resources := make(map[string]Resource)
	yamlErr := yaml.Unmarshal([]byte(resourceData), resources)
	if yamlErr != nil {
		log.Fatalf("error: %v", yamlErr)
	}

	// Monitor for battery changes and update resources
	monitor := NewMonitor(&resources)
	go monitor.Listen()
	defer monitor.Stop()

	// Listen for interrupts and quit
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
