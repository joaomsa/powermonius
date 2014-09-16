package main

import (
	"github.com/guelfey/go.dbus"
	"log"
	"sync"
)

type Monitor struct {
	Resources *map[string]Resource
	DBusCh    chan bool
	DoneCh    chan bool
}

func NewMonitor(resources *map[string]Resource) *Monitor {
	m := &Monitor{
		Resources: resources,
		DoneCh:    make(chan bool),
		DBusCh:    make(chan bool),
	}
	return m
}

func (m *Monitor) UpdateResources(onBattery bool) {
	var wg sync.WaitGroup
	for _, resource := range *m.Resources {
		wg.Add(1)
		go func(r Resource) {
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
	m.DBusCh <- m.CheckDBus()

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
			props := v.Body[1].(map[string]dbus.Variant)
			m.DBusCh <- props["OnBattery"].Value().(bool)
		case <-m.DoneCh:
			return
		}
	}
}

func (m *Monitor) Listen() {
	go m.ListenDBus()

	for {
		select {
		case onBattery := <-m.DBusCh:
			log.Printf("DBus channel: %v\n", onBattery)
			m.UpdateResources(onBattery)
		case <-m.DoneCh:
			return
		}
	}
}

func (m *Monitor) Stop() {
	close(m.DoneCh)
}
