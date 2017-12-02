package main

import (
	"log"
	"os"
)

func initLogFile() *os.File {
	f, err := os.OpenFile("/home/pi/hockey_light.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Printf("error opening file: %v", err)
	}

	log.SetOutput(f)
	return f
}
