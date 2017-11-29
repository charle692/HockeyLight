package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

// TeamName - The team to listen to goals for
const TeamName = "Montréal Canadiens"

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
	f := initLogFile()
	defer f.Close()
	go retrieveSchedule(gameChan, &waitingForGameToStart, &gameStarted)
	go startSSDPServer()

	for {
		select {
		case game := <-gameChan:
			log.Printf("The %s are playing today!\n", TeamName)
			go waitForGameToStart(game, gameStartedChan, &waitingForGameToStart)
		case link := <-gameStartedChan:
			gameStarted = true
			log.Println("The game has started!")
			fmt.Println("The game has started!")
			playHornAndTurnOnLight(pin)
			go listenForGoals(link, goalChan, winningTeam)
		case <-goalChan:
			log.Printf("The %s have scored!\n", TeamName)
			fmt.Println("They have scored!")
			playHornAndTurnOnLight(pin)
		case team := <-winningTeam:
			if team == TeamName {
				fmt.Println("They have won!")
				log.Printf("The %s have won!\n", TeamName)
				playHornAndTurnOnLight(pin)
			}
		}
	}
}

func listenForGoals(link string, goalChan chan bool, winningTeam chan string) {
	ticker := time.NewTicker(time.Second * 2)
	feedData := &FeedData{}
	awayGoals, homeGoals := 0, 0
	homeTeam, awayTeam := Home{}, Away{}
	firstPull := true

	for range ticker.C {
		resp, err := http.Get(Domain + link)

		if err != nil {
			log.Printf("An error while getting live game data: %s\n", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &feedData); err != nil {
			log.Printf("An error while unmarshalling live game data: %s\n", err)
		}

		fmt.Println("Pulling data")
		fmt.Printf("First Pull?: %v\n", firstPull)

		awayTeam = feedData.LiveData.LineScore.Teams.Away
		homeTeam = feedData.LiveData.LineScore.Teams.Home

		if awayTeam.Team.Name == TeamName && awayTeam.Goals > awayGoals {
			if firstPull {
				awayGoals = awayTeam.Goals
			} else {
				awayGoals = awayTeam.Goals
				goalChan <- true
				fmt.Println("They have scored!")
			}
		}

		if homeTeam.Team.Name == TeamName && homeTeam.Goals > homeGoals {
			if firstPull {
				homeGoals = homeTeam.Goals
			} else {
				homeGoals = homeTeam.Goals
				goalChan <- true
				fmt.Println("They have scored!")
			}
		}

		if firstPull {
			firstPull = false
		}

		if feedData.GameData.GameStatus.AbstractGameState == "Final" {
			if awayGoals > homeGoals {
				winningTeam <- awayTeam.Team.Name
			} else {
				winningTeam <- homeTeam.Team.Name
			}

			ticker.Stop()
			return
		}
	}
}

func playHornAndTurnOnLight(pin *rpio.Pin) {
	go turnOnLight(pin)
	go playHorn()
}

func initLogFile() *os.File {
	f, err := os.OpenFile("hockey_light.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Printf("error opening file: %v", err)
	}

	log.SetOutput(f)
	return f
}
