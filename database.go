package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDatabase() *gorm.DB {
	db, err := gorm.Open("sqlite3", "/home/pi/hockey_light.db")
	if err != nil {
		log.Printf("Error connecting to database: %s\n", err)
	}

	db.AutoMigrate(&Team{})
	return db
}
