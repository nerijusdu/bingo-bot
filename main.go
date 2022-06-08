package main

import (
	"bingobot/pkg/auth"
	"bingobot/pkg/db"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := db.NewDatabase()
	defer db.Close()

	auth.StartAllBots(db)
	auth.StartAuthServer(db)
}
