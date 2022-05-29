package bot

import (
	"context"
	"fmt"
	"os"
	"restracker/pkg/bingo"
	"restracker/pkg/db"

	"github.com/shomali11/slacker"
)

type Bot struct {
	db *db.Database
}

func NewBot(database *db.Database) *Bot {
	return &Bot{db: database}
}

func (b *Bot) Run() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bingoMgr := bingo.NewBingoManager(b.db)

	bot.Command("bingo add <item>", &slacker.CommandDefinition{
		Description: "Add a new cell to the bingo board",
		Example:     "bingo add <cell text>",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			text := request.Param("item")
			bingo := bingoMgr.GetOrCreate(botCtx.Event().Channel)
			bingo.AddCell(text)
			response.Reply(bingo.ToString())
		},
	})

	bot.Command("bingo remove <id>", &slacker.CommandDefinition{
		Description: "Remove a cell from the bingo board",
		Example:     "bingo remove <cell id>",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			id := request.IntegerParam("id", 0)
			bingo := bingoMgr.GetOrCreate(botCtx.Event().Channel)
			if !bingo.RemoveCell(id) {
				response.Reply(fmt.Sprintf("Cell %d not found", id))
			}

			response.Reply(bingo.ToString())
		},
	})

	bot.Command("bingo mark <id>", &slacker.CommandDefinition{
		Description: "Mark a cell",
		Example:     "bingo mark <cell id>",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			id := request.IntegerParam("id", 0)
			bingo := bingoMgr.GetOrCreate(botCtx.Event().Channel)
			if !bingo.MarkCell(id) {
				response.Reply(fmt.Sprintf("Cell %d not found", id))
			}

			if bingo.IsCompleted() {
				response.Reply(":bell: :bell: :bell: Bingo! You win! :tada:")
			} else {
				response.Reply(bingo.ToString())
			}
		},
	})

	bot.Command("bingo switch <id1> <id2>", &slacker.CommandDefinition{
		Description: "Switch the cells with the given ids",
		Example:     "bingo switch <id1> <id2>",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			id1 := request.IntegerParam("id1", 0)
			id2 := request.IntegerParam("id2", 0)
			bingo := bingoMgr.GetOrCreate(botCtx.Event().Channel)

			if !bingo.SwitchCells(id1, id2) {
				response.Reply("Invalid cell numbers")
			}

			response.Reply(bingo.ToString())
		},
	})

	bot.Command("bingo reset", &slacker.CommandDefinition{
		Description: "Reset the bingo board",
		Example:     "bingo reset",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			bingo := bingoMgr.GetOrCreate(botCtx.Event().Channel)
			bingo.Reset()
			response.Reply(bingo.ToString())
		},
	})

	bot.Command("bingo list", &slacker.CommandDefinition{
		Description: "Show the bingo items in a list",
		Example:     "bingo list",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			bingo := bingoMgr.GetOrCreate(botCtx.Event().Channel)
			response.Reply(bingo.ToString())
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		panic(err)
	}
}
