package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func playHorn() {
	fileName := "./horns/"
	fileName += strings.Replace(strings.ToLower(TeamName), " ", "_", -1)
	fileName += ".wav"

	cmdArgs := []string{fileName}
	fmt.Printf("Playing horn: %s\n", fileName)
	_, err := exec.Command("aplay", cmdArgs...).Output()

	if err != nil {
		fmt.Printf("There was an error running aplay: %s\n", err)
	}
}
