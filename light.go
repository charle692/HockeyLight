package main

import (
	"fmt"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

func turnOnLight(pin *rpio.Pin) {
	pin.Low()
	ticker := time.NewTicker(time.Second * 30)

	for range ticker.C {
		fmt.Printf("30 seconds is up!\n")
		pin.High()
	}
}

func initializeGPIOPin() *rpio.Pin {
	err := rpio.Open()

	if err != nil {
		panic(err)
	}

	pin := rpio.Pin(4)
	pin.Output()
	pin.High()

	return &pin
}
