package main

import (
	"bingobot/pkg/auth"
	"bingobot/pkg/db"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	defer fmt.Println("I literally died")
	db := db.NewDatabase()
	defer db.Close()

	auth.StartAllBots(db)
	auth.StartAuthServer(db)
}
