package main

import (
	"github.com/guelfey/go.dbus"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"sync"
)

func checkOnBattery() bool {
	conn, connErr := dbus.SystemBus()
	if connErr != nil {
		log.Fatalf("Failed to connect to bus: %v", connErr)
	}

	msg, getErr := conn.
		Object("org.freedesktop.UPower", "/org/freedesktop/UPower").
		GetProperty("org.freedesktop.UPower.OnBattery")
	if getErr != nil {
		log.Fatalf("Failed to get property: ", getErr)
	}
	return msg.Value().(bool)
}

func listenForBattery(battery chan bool) {
	conn, connErr := dbus.SystemBus()
	if connErr != nil {
		log.Fatalf("Failed to connect to bus: %v", connErr)
	}

	matchRule := "type='signal',sender='org.freedesktop.UPower',path='/org/freedesktop/UPower'"
	call := conn.
		BusObject().
		Call("org.freedesktop.DBus.AddMatch", 0, matchRule)

	if call.Err != nil {
		log.Fatalf("Failed to add match", call.Err)
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	for v := range c {
		props := v.Body[1].(map[string]dbus.Variant)
		battery <- props["OnBattery"].Value().(bool)
	}
}

func updateResources(resources map[string]Resource, onBattery bool) {
	var wg sync.WaitGroup
	for key := range resources {
		wg.Add(1)
		go func(resource Resource) {
			defer wg.Done()
			if onBattery {
				resource.Unplug()
			} else {
				resource.Plug()
			}
		}(resources[key])
	}
	wg.Wait()
}

func main() {
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

	onBattery := checkOnBattery()
	updateResources(resources, onBattery)

	c := make(chan bool)
	go listenForBattery(c)
	for v := range c {
		updateResources(resources, v)
	}
}
