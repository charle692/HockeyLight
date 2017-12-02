package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
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

func retrieveSchedule(gameChan chan Game, waitingForGameToStart *bool, gameStarted *bool) {
	ticker := time.NewTicker(time.Minute * 1)

	for range ticker.C {
		if !*gameStarted && !*waitingForGameToStart {
			getTodaysSchedule(gameChan)
		}
	}
}

func getTodaysSchedule(gameChan chan Game) {
	schedule := &Schedule{}
	time := time.Now()
	year, month, day := time.Date()
	date := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	resp, err := http.Get(Domain + SchedulePath + "?startDate=" + date + "&endDate=" + date)

	if err != nil {
		fmt.Printf("An error occured while retrieving today's game data")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("An error occured while retrieving today's game data")
	}

	gjson.Unmarshal(body, &schedule)
	resp.Body.Close()
	isTeamPlayingToday(schedule, gameChan)
}

// Modify this to determine if today's game has already passed
func isTeamPlayingToday(schedule *Schedule, gameChan chan Game) {
	for i := 0; i < len(schedule.Dates); i++ {
		date := schedule.Dates[i]

		for x := 0; x < len(date.Games); x++ {
			game := date.Games[x]
			if (game.Teams.Home.Team.Name == getSelectedTeamName() || game.Teams.Away.Team.Name == getSelectedTeamName()) && game.Status.DetailedState != "Final" {
				gameChan <- game
			}
		}
	}
}

func waitForGameToStart(game Game, gameStartedChan chan string, waitingForGameToStart *bool) {
	setWaitingForGameToStart(waitingForGameToStart, true)
	defer setWaitingForGameToStart(waitingForGameToStart, false)
	startDate, _ := time.Parse(time.RFC3339, game.GameDate)
	startDateInUnix := startDate.Unix()
	currentTime := time.Now().Unix()

	if startDateInUnix-currentTime > 0 {
		timeUntilGameStarts := time.Duration(startDateInUnix - currentTime)
		time.Sleep((timeUntilGameStarts / 60) * time.Minute)
		gameStartedChan <- game.Link
		return
	}

	gameStartedChan <- game.Link
}

func setWaitingForGameToStart(waitingForGameToStart *bool, status bool) {
	*waitingForGameToStart = status
}
