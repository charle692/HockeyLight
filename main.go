package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var gameStarted = false
var waitingForGameToStart = false

// TeamName - The team to listen to goals for
const TeamName = "Montréal Canadiens"

// Domain - The domain of the api
const Domain = "https://statsapi.web.nhl.com"

// FeedData - Contains game feed data
type FeedData struct {
	LiveData LiveData `json:"liveData"`
}

// LiveData - Contains Live Game data
type LiveData struct {
	LineScore LineScore `json:"linescore"`
}

// LineScore - Contains Period and Team related data
type LineScore struct {
	CurrentPeriodOrdinal string `json:"currentPeriodOrdinal"`
	Teams                Teams  `json:"teams"`
}

// Teams - Contains Home and Away teams
type Teams struct {
	Home Home `json:"home"`
	Away Away `json:"away"`
}

// Home - Contains home team data
type Home struct {
	Team  Team `json:"team"`
	Goals int  `json:"goals"`
}

// Away - Contains away team data
type Away struct {
	Team  Team `json:"team"`
	Goals int  `json:"goals"`
}

// Team - Contains team name
type Team struct {
	Name string `json:"name"`
}

// TeamData - Contains the team name
type TeamData struct {
	Name string `json:"name"`
}

func main() {
	gameChan := make(chan Game)
	gameStartedChan := make(chan string)
	goalChan := make(chan bool)
	isItWorking := make(chan string)
	pin := initializeGPIOPin()
	go retrieveSchedule(gameChan)

	for {
		select {
		case game := <-gameChan:
			go waitForGameToStart(game, gameStartedChan)
			fmt.Printf("The %s are playing today!\n", TeamName)
		case link := <-gameStartedChan:
			go turnOnLight(pin)
			go listenForGoals(link, goalChan, isItWorking)
			fmt.Println("The game has started!")
		case <-goalChan:
			fmt.Printf("The %s have scored!", TeamName)
			go turnOnLight(pin)
		}
	}
}

func listenForGoals(link string, goalChan chan bool, isItWorking chan string) {
	ticker := time.NewTicker(time.Second * 3)
	feedData := &FeedData{}
	awayGoals := 0
	homeGoals := 0

	for range ticker.C {
		resp, err := http.Get(Domain + link)

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &feedData); err != nil {
			panic(err)
		}

		if feedData.LiveData.LineScore.Teams.Away.Team.Name == TeamName && feedData.LiveData.LineScore.Teams.Away.Goals > awayGoals {
			awayGoals = feedData.LiveData.LineScore.Teams.Away.Goals
			goalChan <- true
		}

		if feedData.LiveData.LineScore.Teams.Home.Team.Name == TeamName && feedData.LiveData.LineScore.Teams.Home.Goals > homeGoals {
			homeGoals = feedData.LiveData.LineScore.Teams.Home.Goals
			goalChan <- true
		}

		if feedData.LiveData.LineScore.CurrentPeriodOrdinal == "Final" {
			return
		}
	}
}

func setWaitingForGameToStart(waiting bool) {
	waitingForGameToStart = waiting
}
