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
			log.Fatalf("\"%v\" not a valid value for \"when\"\n", when)
		}
	} else {
		r.when = "charging"
	}

	start, ok := data.(map[interface{}]interface{})["start"].(string)
	if !ok {
		log.Fatalf("Failed to define start command\n")
	}
	r.startCmd = start

	stop, ok := data.(map[interface{}]interface{})["stop"].(string)
	if !ok {
		log.Fatalf("Failed to define stop command\n")
	}
	r.stopCmd = stop

	status, ok := data.(map[interface{}]interface{})["status"].(string)
	if !ok {
		log.Fatalf("Failed to define status command\n")
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
	log.Printf("[status/%v] %v\n", r.Name, r.statusCmd)
	return cmd.Run()
}

func (r *Resource) Start() error {
	cmd := exec.Command("bash", "-c", "nohup "+r.startCmd+" &")
	log.Printf("[start/%v] %v\n", r.Name, r.startCmd)
	err := cmd.Start()
	if err != nil {
		log.Printf("[start/%v/err] %v\n", r.Name, err.Error())
	}
	exit_err := cmd.Wait()
	if exit_err != nil {
		log.Printf("[start/%v/exit] %v\n", r.Name, exit_err)
	}
	return err
}

func (r *Resource) Stop() error {
	cmd := exec.Command("bash", "-c", "nohup "+r.stopCmd+" &")
	log.Printf("[stop/%v] %v\n", r.Name, r.stopCmd)
	err := cmd.Start()
	if err != nil {
		log.Printf("[stop/%v/err] %v\n", r.Name, err.Error())
	}
	exit_err := cmd.Wait()
	if exit_err != nil {
		log.Printf("[stop/%v/exit] %v\n", r.Name, exit_err)
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
