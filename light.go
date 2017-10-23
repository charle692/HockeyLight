package main

import (
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

func turnOnLight(pin *rpio.Pin) {
	ticker := time.NewTicker(time.Second * 69)
	pin.Low()

	for range ticker.C {
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
