package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charle692/hockeyLight/mp3"
	rpio "github.com/stianeikeland/go-rpio"
)

var gpioPin *rpio.Pin
var teamSelected chan bool

// Settings - contains the hockey light settings
type Settings struct {
	Team  Team  `json:"team"`
	Delay Delay `json:"delay"`
}

func saveSettings(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	settings := &Settings{}
	if err := decoder.Decode(settings); err != nil {
		fmt.Println(err)
	}

	team := getTeamByName(settings.Team.Name)
	if team.Name != "" {
		db.Model(getSelectedTeam()).Update("selected", false)
		db.Model(team).Update("selected", true)
		teamSelected <- true
	}

	db.Model(getDelay()).Update("value", settings.Delay.Value)
	settings.Team = *getSelectedTeam()
	settings.Delay = *getDelay()

	json, err := json.Marshal(settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
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
	json, err := json.Marshal(getDelay())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func playHornHandler(w http.ResponseWriter, r *http.Request) {
	go turnOnLight(gpioPin)
	go mp3.Play(hornFilePath())
}

func startHTTPServer(pin *rpio.Pin, newTeamSelected chan bool) {
	gpioPin = pin
	teamSelected = newTeamSelected
	http.HandleFunc("/post/settings", saveSettings)
	http.HandleFunc("/get/teams", getTeamsHandler)
	http.HandleFunc("/get/delay", getDelayHandler)
	http.HandleFunc("/play_horn", playHornHandler)
	http.ListenAndServe(":8080", nil)
}
