package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func playHorn(hornDone chan bool) {
	fileName := "./horns/"
	fileName += strings.Replace(strings.ToLower(TeamName), " ", "_", -1)
	fileName += ".wav"
	fmt.Printf("%s", fileName)

	cmdArgs := []string{fileName}
	_, err := exec.Command("aplay", cmdArgs...).Output()

	if err != nil {
		fmt.Printf("There was an error running aplay: %s\n", err)
	}

	hornDone <- true
}
