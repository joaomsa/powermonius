package main

import (
	"github.com/guelfey/go.dbus"
	"log"
	"sync"
)

type Monitor struct {
	Resources   *map[string]*Resource
	ResourceLck sync.Mutex
	OnBatteryCh chan bool
	DoneCh      chan bool
}

func NewMonitor() *Monitor {
	resources := make(map[string]*Resource)
	m := &Monitor{
		Resources:   &resources,
		OnBatteryCh: make(chan bool),
		DoneCh:      make(chan bool),
	}
	return m
}

func (m *Monitor) SetResources(resources *map[string]*Resource) {
	m.ResourceLck.Lock()
	defer m.ResourceLck.Unlock()

	m.Resources = resources
}

func (m *Monitor) UpdateResources(onBattery bool) {
	m.ResourceLck.Lock()
	defer m.ResourceLck.Unlock()

	var wg sync.WaitGroup
	for _, resource := range *m.Resources {
		wg.Add(1)
		go func(r *Resource) {
			defer wg.Done()
			if onBattery {
				r.Unplug()
			} else {
				r.Plug()
			}
		}(resource)
	}
	wg.Wait()
}

func (m *Monitor) CheckDBus() (onBattery bool) {
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

	onBattery = msg.Value().(bool)
	return
}

func (m *Monitor) ListenDBus() {
	m.OnBatteryCh <- m.CheckDBus()

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

	c := make(chan *dbus.Signal)
	conn.Signal(c)

	for {
		select {
		case v := <-c:
			switch v.Path {
			case "/org/freedesktop/UPower":
				changedProperties := v.Body[1].(map[string]dbus.Variant)
				if property, ok := changedProperties["OnBattery"]; ok {
					onBattery := property.Value().(bool)
					m.OnBatteryCh <- onBattery
				}
			}
		case <-m.DoneCh:
			return
		}
	}
}

func (m *Monitor) Listen() {
	go m.ListenDBus()

	for {
		select {
		case onBattery := <-m.OnBatteryCh:
			m.UpdateResources(onBattery)
		case <-m.DoneCh:
			return
		}
	}
}

func (m *Monitor) Stop() {
	log.Printf("Stopping monitor")
	close(m.DoneCh)
}
