package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

// TeamName - The team to listen to goals for
const TeamName = "New York Rangers"

// Domain - The domain of the api
const Domain = "https://statsapi.web.nhl.com"

// FeedData - Contains game feed data
type FeedData struct {
	LiveData LiveData `json:"liveData"`
	GameData GameData `json:"gameData"`
}

// LiveData - Contains Live Game data
type LiveData struct {
	LineScore LineScore `json:"linescore"`
}

// GameData - Contains the game status
type GameData struct {
	GameStatus GameStatus `json:"status"`
}

// GameStatus - Contains the game status
type GameStatus struct {
	AbstractGameState string `json:"abstractGameState"`
}

// LineScore - Contains Period and Team related data
type LineScore struct {
	Teams Teams `json:"teams"`
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
	waitingForGameToStart, gameStarted := false, false
	gameChan, gameStartedChan, goalChan, winningTeam := make(chan Game), make(chan string), make(chan bool), make(chan string)
	pin := initializeGPIOPin()
	go retrieveSchedule(gameChan, &waitingForGameToStart, &gameStarted)

	for {
		select {
		case game := <-gameChan:
			fmt.Printf("The %s are playing today!\n", TeamName)
			go waitForGameToStart(game, gameStartedChan, &waitingForGameToStart)
		case link := <-gameStartedChan:
			gameStarted = true
			fmt.Println("The game has started!")
			playHornAndTurnOnLight(pin)
			go listenForGoals(link, goalChan, winningTeam)
		case <-goalChan:
			fmt.Printf("The %s have scored!\n", TeamName)
			playHornAndTurnOnLight(pin)
		case team := <-winningTeam:
			if team == TeamName {
				fmt.Printf("The %s have won!\n", TeamName)
				playHornAndTurnOnLight(pin)
			}
		}
	}
}

func listenForGoals(link string, goalChan chan bool, winningTeam chan string) {
	ticker := time.NewTicker(time.Second * 3)
	feedData := &FeedData{}
	awayGoals := 0
	homeGoals := 0
	homeTeam := Home{}
	awayTeam := Away{}

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

		awayTeam = feedData.LiveData.LineScore.Teams.Away
		homeTeam = feedData.LiveData.LineScore.Teams.Home

		if awayTeam.Team.Name == TeamName && awayTeam.Goals > awayGoals {
			awayGoals = awayTeam.Goals
			goalChan <- true
		}

		if homeTeam.Team.Name == TeamName && homeTeam.Goals > homeGoals {
			homeGoals = homeTeam.Goals
			goalChan <- true
		}

		if feedData.GameData.GameStatus.AbstractGameState == "Final" {
			if awayGoals > homeGoals {
				winningTeam <- awayTeam.Team.Name
			} else {
				winningTeam <- homeTeam.Team.Name
			}

			return
		}
	}
}

func playHornAndTurnOnLight(pin *rpio.Pin) {
	hornDone := make(chan struct{})
	go playHorn(hornDone)
	go turnOnLight(pin, hornDone)
}
