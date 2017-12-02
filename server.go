package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	rpio "github.com/stianeikeland/go-rpio"
)

var gpioPin *rpio.Pin

func postTeamHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	team := &Team{}
	if err := decoder.Decode(team); err != nil {
		fmt.Println(err)
	}

	team = getTeamByName(team.Name)
	if team.Name != "" {
		db.Model(getSelectedTeam()).Update("selected", false)
		db.Model(team).Update("selected", true)
	}

	json, err := json.Marshal(team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func postDelayHandler(w http.ResponseWriter, r *http.Request) {
	// password := r.FormValue("password")
	// networkData := strings.Split(r.FormValue("networkName"), " - ")
	// ssid := networkData[0]
	// securityType := networkData[1]

	http.Redirect(w, r, "/views/success", http.StatusFound)
}

func getTeamsHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(getTeams())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
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
	http.HandleFunc("/get/teams", getTeamsHandler)
	http.HandleFunc("/get/delay", getDelayHandler)
	http.HandleFunc("/play_horn", playHornHandler)
	http.ListenAndServe(":8080", nil)
}
