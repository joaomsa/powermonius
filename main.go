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
	monitor := NewMonitor()
	monitor.SetResources(&resources)
	go monitor.Listen()
	defer monitor.Stop()

	// Listen for reload, interrupts and quit
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	for {
		select {
		case sig := <-sigCh:
			log.Printf("Received signal: %v", sig)
			if sig == syscall.SIGUSR1 {
				log.Printf("Reloading configuration")
				resources = loadResourceFile(resourceFile)
				monitor.SetResources(&resources)
				onBattery := monitor.CheckDBus()
				monitor.UpdateResources(onBattery)
			} else {
				return
			}
		}
	}
}
