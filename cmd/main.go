package main

import (
	"fmt"
	"restracker/pkg/bot"
	"restracker/pkg/db"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := db.NewDatabase()
	defer db.Close()

	bot := bot.NewBot(db)

	fmt.Println("Starting application")
	bot.Run()
}
