package main

import (
	"restracker/pkg/auth"
	"restracker/pkg/db"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := db.NewDatabase()
	defer db.Close()

	auth.StartAllBots(db)
	auth.StartAuthServer(db)
}
