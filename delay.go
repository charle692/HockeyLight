package main

import "github.com/jinzhu/gorm"

// Delay contains delay before setting the light off
type Delay struct {
	gorm.Model
	Value string `json:"value"`
}

func getDelay() *Delay {
	delay := &Delay{}
	db.First(delay)
	return delay
}
