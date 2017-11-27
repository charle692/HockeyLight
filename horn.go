package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func playHorn() {
	fileName := "/home/pi/horns/"
	fileName += strings.Replace(strings.ToLower(TeamName), " ", "_", -1)
	fileName += ".mp3"

	cmdArgs := []string{fileName}
	_, err := exec.Command("mpg123", cmdArgs...).Output()

	if err != nil {
		fmt.Printf("There was an error running mpg123: %s\n", err)
	}
}
