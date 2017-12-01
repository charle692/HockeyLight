package main

import (
	"net/http"

	rpio "github.com/stianeikeland/go-rpio"
)

var gpioPin *rpio.Pin

func postTeamHandler(w http.ResponseWriter, r *http.Request) {
	// password := r.FormValue("password")
	// networkData := strings.Split(r.FormValue("networkName"), " - ")
	// ssid := networkData[0]
	// securityType := networkData[1]

	http.Redirect(w, r, "/views/success", http.StatusFound)
}

func postDelayHandler(w http.ResponseWriter, r *http.Request) {
	// password := r.FormValue("password")
	// networkData := strings.Split(r.FormValue("networkName"), " - ")
	// ssid := networkData[0]
	// securityType := networkData[1]

	http.Redirect(w, r, "/views/success", http.StatusFound)
}

func getTeamHandler(w http.ResponseWriter, r *http.Request) {
}

func getDelayHandler(w http.ResponseWriter, r *http.Request) {
}

func playHornHandler(w http.ResponseWriter, r *http.Request) {
	playHornAndTurnOnLight(gpioPin)
}

func startHTTPServer(pin *rpio.Pin) {
	gpioPin = pin
	http.HandleFunc("/post/team", postTeamHandler)
	http.HandleFunc("/post/delay", postDelayHandler)
	http.HandleFunc("/get/team", getTeamHandler)
	http.HandleFunc("/get/delay", getDelayHandler)
	http.HandleFunc("/play_horn", playHornHandler)
	http.ListenAndServe(":8080", nil)
}
