package main

import "github.com/jinzhu/gorm"

// Team - Contains team data
type Team struct {
	gorm.Model
	Selected bool
	Name     string `json:"name"`
}

func getSelectedTeamName() string {
	return getSelectedTeam().Name
}

func getSelectedTeam() *Team {
	selectedTeam := &Team{}
	db.Where(&Team{Selected: true}).First(selectedTeam)
	return selectedTeam
}

func getTeams() *[]Team {
	teams := &[]Team{}
	db.Order("name asc").Find(teams)
	return teams
}

func getTeamByName(teamName string) *Team {
	team := &Team{}
	db.Where(&Team{Name: teamName}).First(team)
	return team
}
