package main

import (
	"fmt"
	"github.com/guelfey/go.dbus"
	"os"
)

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to bus:", err)
		os.Exit(1)
	}

	var onBattery bool

	msg, getErr := conn.
		Object("org.freedesktop.UPower", "/org/freedesktop/UPower").
		GetProperty("org.freedesktop.UPower.OnBattery")
	if getErr != nil {
		fmt.Fprintln(os.Stderr, "Failed to get property:", getErr)
		os.Exit(1)
	}
	onBattery = msg.Value().(bool)

	fmt.Printf("onBattery: %v\n", onBattery)
	fmt.Println("-------------------")

	matchRule := "type='signal',sender='org.freedesktop.UPower',path='/org/freedesktop/UPower'"
	call := conn.
		BusObject().
		Call("org.freedesktop.DBus.AddMatch", 0, matchRule)
)
	if call.Err != nil {
		fmt.Fprintln(os.Stderr, "Failed to add match:", call.Err)
		os.Exit(1)
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	for v := range c {
		props := v.Body[1].(map[string]dbus.Variant)
		onBattery = props["OnBattery"].Value().(bool)

		fmt.Printf("onBattery: %v\n", onBattery)
		fmt.Println("-------------------")
	}
}
