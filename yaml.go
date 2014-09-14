package main

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"log"
)

var data = `
dropbox:
  # Indicate when to run program
  plugged: running
  battery: stopped

  # Command that starts program
  start: dropbox start
  # Command that stops program
  stop: dropbox stop
  # Command that returns 0 if running and nonzero otherwise
  status: "! dropbox running"

transmission:
  # "plugged" defaults to "running"
  # "battery" defaults to "stopped"

  start: transmission-gtk -m
  stop: pkill -f transmission-gtk
  status: pgrep -f transmission-gtk

tracker:
  # Only start program that requires gui if display server is running, defaults to yes
  require_gui: no

  start: tracker-control -s
  stop: tracker-control -t
  status: tracker-control -l | grep --q "Found [^0]\d* miners running"
`

type Config struct {
	Plugged     string
	Battery     string
	Require_Gui bool
	Start       string
	Stop        string
	Status      string
}

func main() {
	entry := make(map[string]Config)

	err := yaml.Unmarshal([]byte(data), entry)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- entry:\n|%v|\n\n", entry)
	fmt.Printf("--- entry:\n|%s|\n\n", entry["dropbox"].Require_Gui)
}