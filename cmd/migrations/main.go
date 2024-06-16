package main

import (
	"log"

	"github.com/tjons/text-to-trade/pkg/model"
)

func main() {
	db, err := model.Connect()
	if err != nil {
		log.Fatal(err)
	}
	if err = db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(&model.Watchlist{}); err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(&model.Notification{}); err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(&model.ChatMessage{}); err != nil {
		log.Fatal(err)
	}

	log.Default().Printf("Migrated successfully")
}
