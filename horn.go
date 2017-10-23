package main

import (
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func playHorn(hornDone chan struct{}) {
	done := make(chan struct{})
	f, err := os.Open("horns/" + strings.Replace(strings.ToLower(TeamName), " ", "_", -1) + ".wav")

	if err != nil {
		panic(err)
	}

	s, format, err := wav.Decode(f)

	if err != nil {
		panic(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2))
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(hornDone)
		close(done)
	})))
	<-done
}
