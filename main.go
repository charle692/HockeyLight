package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/charle692/hockeyLight/mp3"
	"github.com/charle692/hockeyLight/ssdp"
	rpio "github.com/stianeikeland/go-rpio"
)

// Domain - The domain of the api
const Domain = "https://statsapi.web.nhl.com"

var db = initDatabase()

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

// TeamData - Contains the team name
type TeamData struct {
	Name string `json:"name"`
}

func main() {
	gameStarted := false
	goalChan, winningTeam := make(chan bool), make(chan string)
	newTeamSelected := make(chan bool)

	f := initLogFile()
	defer f.Close()

	db.LogMode(true)
	defer db.Close()

	pin := initializeGPIOPin()
	gameStartedChan := waitForGameToStart(newTeamSelected, &gameStarted)
	ssdp.Start("my:hockey-light", "Hockey Light SSDP")
	go startHTTPServer(pin, newTeamSelected)

	for {
		select {
		case game := <-gameStartedChan:
			gameStarted = true
			if !strings.Contains(game.Status.DetailedState, "In Progress") {
				playHornAndTurnOnLight(pin)
			}
			go listenForGoals(game.Link, goalChan, winningTeam, newTeamSelected, &gameStarted)
		case <-goalChan:
			playHornAndTurnOnLight(pin)
		case team := <-winningTeam:
			if team == getSelectedTeamName() {
				gameStarted = false
				playHornAndTurnOnLight(pin)
			}
		}
	}
}

func listenForGoals(link string, goalChan chan bool, winningTeam chan string, newTeamSelected chan bool, gameStarted *bool) {
	ticker := time.NewTicker(time.Second * 2)
	awayGoals, homeGoals := 0, 0
	selectedTeam := getSelectedTeamName()
	firstPull := true

	for range ticker.C {
		select {
		case <-newTeamSelected:
			*gameStarted = false
			ticker.Stop()
			return
		default:
			awayTeam, homeTeam, gameState := retrieveGameData(link)

			if firstPull {
				firstPull = false
				awayGoals = awayTeam.Goals
				homeGoals = homeTeam.Goals
			} else {
				if awayTeam.Team.Name == selectedTeam && awayTeam.Goals > awayGoals {
					log.Printf("The %s have scored!\n", getSelectedTeamName())
					awayGoals = awayTeam.Goals
					goalChan <- true
				}

				if homeTeam.Team.Name == selectedTeam && homeTeam.Goals > homeGoals {
					log.Printf("The %s have scored!\n", getSelectedTeamName())
					homeGoals = homeTeam.Goals
					goalChan <- true
				}

				if gameState == "Final" {
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
	}
}

func retrieveGameData(link string) (Away, Home, string) {
	feedData := &FeedData{}
	resp, err := http.Get(Domain + link)

	if err != nil {
		log.Printf("An error while getting live game data: %s\n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, feedData); err != nil {
		log.Printf("An error while unmarshalling live game data: %s\n", err)
	}

	return feedData.LiveData.LineScore.Teams.Away,
		feedData.LiveData.LineScore.Teams.Home,
		feedData.GameData.GameStatus.AbstractGameState
}

func playHornAndTurnOnLight(pin *rpio.Pin) {
	delay, _ := strconv.Atoi(getDelay().Value)
	time.Sleep(time.Second * time.Duration(delay))
	go turnOnLight(pin)
	go mp3.Play(hornFilePath())
}

func hornFilePath() string {
	fileName := "/home/pi/horns/"
	fileName += strings.Replace(strings.ToLower(getSelectedTeamName()), " ", "_", -1)
	fileName += ".mp3"
	return fileName
}
