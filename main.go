package main

import (
	"github.com/guelfey/go.dbus"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os/exec"
	"sync"
)

type Resource struct {
	when   string
	Start  string
	Stop   string
	Status string
}

func (r *Resource) SetYAML(tag string, data interface{}) bool {
	when, ok := data.(map[interface{}]interface{})["when"].(string)
	if ok {
		switch when {
		case "always", "charging", "discharging", "never":
			r.when = when
		default:
			log.Fatalf("\"%v\" not a valid value for \"when\"", when)
		}
	} else {
		r.when = "charging"
	}

	start, ok := data.(map[interface{}]interface{})["start"].(string)
	if !ok {
		log.Fatalf("Failed to define start command")
	}
	r.Start = start

	stop, ok := data.(map[interface{}]interface{})["stop"].(string)
	if !ok {
		log.Fatalf("Failed to define stop command")
	}
	r.Stop = stop

	status, ok := data.(map[interface{}]interface{})["status"].(string)
	if !ok {
		log.Fatalf("Failed to define status command")
	}
	r.Status = status

	return true
}

func (r *Resource) WhenPlugged() bool {
	switch r.when {
	case "always", "charging":
		return true
	}
	return false
}

func (r *Resource) WhenUnplugged() bool {
	switch r.when {
	case "always", "discharging":
		return true
	}
	return false
}

func (r Resource) GetStatus() error {
	cmd := exec.Command("bash", "-c", r.Status)
	log.Printf("[status] %v\n", r.Status)
	return cmd.Run()
}

func (r Resource) RunStart() error {
	cmd := exec.Command("bash", "-c", r.Start)
	log.Printf("[start] %v\n", r.Start)
	return cmd.Start()
}

func (r Resource) RunStop() error {
	cmd := exec.Command("bash", "-c", r.Stop)
	log.Printf("[stop] %v\n", r.Stop)
	return cmd.Start()
}

func (r Resource) Plug() {
	status := r.GetStatus()
	if status != nil {
		if r.WhenPlugged() {
			err := r.RunStart()
			if err != nil {
				log.Printf("[%v] %v", r.Start, err.Error())
			}
		}
	} else {
		if !r.WhenPlugged() {
			err := r.RunStop()
			if err != nil {
				log.Printf("[%v] %v", r.Stop, err.Error())
			}
		}
	}
}

func (r Resource) Unplug() {
	status := r.GetStatus()
	if status != nil {
		if r.WhenUnplugged() {
			err := r.RunStart()
			if err != nil {
				log.Printf("[%v] %v", r.Start, err.Error())
			}
		}
	} else {
		if !r.WhenUnplugged() {
			err := r.RunStop()
			if err != nil {
				log.Printf("[%v] %v", r.Stop, err.Error())
			}
		}
	}
}

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
