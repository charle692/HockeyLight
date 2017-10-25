package main

import rpio "github.com/stianeikeland/go-rpio"

func turnOnLight(pin *rpio.Pin, hornDone chan bool) {
	pin.Low()

	select {
	case <-hornDone:
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
