package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func playHorn() {
	done := make(chan struct{})
	f, _ := os.Open("vancouver_horn.wav")
	s, format, _ := wav.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2))
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))
	<-done
}
