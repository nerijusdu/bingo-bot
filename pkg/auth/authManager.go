package auth

import (
	"os"
	"restracker/pkg/bot"
	"restracker/pkg/db"
)

func StartAllBots(db *db.Database) error {
	authRepo := NewAuthRepository(db)
	tokens, err := authRepo.GetAllTokens()
	if err != nil {
		return err
	}

	for _, tokenData := range tokens {
		go bot.NewBot(db, tokenData.AccessToken, tokenData.Team.ID).Run()
	}

	defaultToken := os.Getenv("SLACK_BOT_TOKEN")
	if defaultToken != "" {
		go bot.NewBot(db, defaultToken, "localToken").Run()
	}

	return nil
}
