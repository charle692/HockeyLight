package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Schedule - Contains game schedule of the day
type Schedule struct {
	Dates Dates `json:"dates"`
}

// Dates - Contains Dates
type Dates []struct {
	Date       string `json:"date"`
	TotalGames int    `json:"totalGames"`
	Games      Games  `json:"games"`
}

// Games - Contains Games of the day
type Games []struct {
	Link     string `json:"link"`
	GameDate string `json:"gameDate"`
	Teams    Teams  `json:"teams"`
	Status   Status `json:"status"`
}

// Game - Contains game data
type Game struct {
	Link     string `json:"link"`
	GameDate string `json:"gameDate"`
	Teams    Teams  `json:"teams"`
	Status   Status `json:"status"`
}

// Status - Contains game status information
type Status struct {
	DetailedState string `json:"detailedState"`
}

// SchedulePath - Path to schedule
const SchedulePath = "/api/v1/schedule"

func waitForGameToStart(newTeamSelected chan bool, gameStarted *bool) <-chan string {
	gameStartedChan := make(chan string)
	go func() {
		for range time.NewTicker(time.Second * 30).C {
			if !*gameStarted {
				retrieveTodaysSchedule(newTeamSelected, gameStartedChan)
			}
		}
	}()
	return gameStartedChan
}

func retrieveTodaysSchedule(newTeamSelected chan bool, gameStartedChan chan string) {
	schedule := &Schedule{}
	date := today()
	resp, err := http.Get(Domain + SchedulePath + "?startDate=" + date + "&endDate=" + date)

	if err != nil {
		fmt.Printf("An error occured while retrieving today's game data")
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("An error occured while retrieving today's game data")
	}

	json.Unmarshal(body, schedule)
	isTeamPlayingToday(schedule, newTeamSelected, gameStartedChan)
}

func isTeamPlayingToday(schedule *Schedule, newTeamSelected chan bool, gameStartedChan chan string) {
	for i := 0; i < len(schedule.Dates); i++ {
		date := schedule.Dates[i]

		for x := 0; x < len(date.Games); x++ {
			game := date.Games[x]
			selectedTeam := getSelectedTeamName()

			if (game.Teams.Home.Team.Name == selectedTeam || game.Teams.Away.Team.Name == selectedTeam) && game.Status.DetailedState != "Final" {
				log.Printf("The %s are playing today!\n", selectedTeam)
				waitUntilGameStarts(game, newTeamSelected, gameStartedChan)
				return
			}
		}
	}
}

func waitUntilGameStarts(game Game, newTeamSelected chan bool, gameStartedChan chan string) {
	startDate, _ := time.Parse(time.RFC3339, game.GameDate)
	startDateInUnix := startDate.Unix()
	currentTime := time.Now().Unix()

	if startDateInUnix-currentTime > 0 {
		timeUntilGameStarts := time.Duration(startDateInUnix - currentTime)

		select {
		case <-newTeamSelected:
			log.Println("New team selected, stopped waiting for game.")
			return
		case <-time.After(time.Minute * (timeUntilGameStarts / 60)):
			log.Println("The game has started!")
			gameStartedChan <- game.Link
			return
		}
	}

	log.Println("The game has started!")
	gameStartedChan <- game.Link
}

func today() string {
	time := time.Now()
	year, month, day := time.Date()
	return strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
}
