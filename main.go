package main

import (
	"log"
	"os"
	"os/signal"
	"os/user"
	"path"
	"syscall"
)

func main() {
	// Load and parse resource file
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	resourceFile := path.Join(usr.HomeDir, ".config/powermonius/resources.yaml")
	resources := loadResourceFile(resourceFile)

	// Monitor for battery changes and update resources
	monitor := NewMonitor(&resources)
	go monitor.Listen()
	defer monitor.Stop()

	// Listen for interrupts and quit
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
