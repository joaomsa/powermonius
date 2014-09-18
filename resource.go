package main

import (
	"log"
	"os/exec"
)

type Resource struct {
	Name      string
	when      string
	startCmd  string
	stopCmd   string
	statusCmd string
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
	r.startCmd = start

	stop, ok := data.(map[interface{}]interface{})["stop"].(string)
	if !ok {
		log.Fatalf("Failed to define stop command")
	}
	r.stopCmd = stop

	status, ok := data.(map[interface{}]interface{})["status"].(string)
	if !ok {
		log.Fatalf("Failed to define status command")
	}
	r.statusCmd = status

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

func (r *Resource) Status() error {
	cmd := exec.Command("bash", "-c", r.statusCmd)
	log.Printf("[status] %v\n", r.statusCmd)
	return cmd.Run()
}

func (r *Resource) Start() error {
	cmd := exec.Command("bash", "-c", r.startCmd)
	log.Printf("[start] %v\n", r.startCmd)
	err := cmd.Start()
	if err != nil {
		log.Printf("[%v] %v", r.startCmd, err.Error())
	}
	return err
}

func (r *Resource) Stop() error {
	cmd := exec.Command("bash", "-c", r.stopCmd)
	log.Printf("[stop] %v\n", r.stopCmd)
	err := cmd.Start()
	if err != nil {
		log.Printf("[%v] %v", r.stopCmd, err.Error())
	}
	return err
}

func (r *Resource) Plug() {
	status := r.Status()
	if status != nil {
		if r.WhenPlugged() {
			r.Start()
		}
	} else {
		if !r.WhenPlugged() {
			r.Stop()
		}
	}
}

func (r *Resource) Unplug() {
	status := r.Status()
	if status != nil {
		if r.WhenUnplugged() {
			r.Start()
		}
	} else {
		if !r.WhenUnplugged() {
			r.Stop()
		}
	}
}

func (r *Resource) Suspend() {
	status := r.Status()
	if status == nil {
		r.Stop()
	}
}
