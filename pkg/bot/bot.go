package bot

import (
	"bingobot/pkg/bingo"
	"bingobot/pkg/db"
	"bingobot/pkg/visual"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

type Bot struct {
	db     *db.Database
	token  string
	teamId string
}

func NewBot(database *db.Database, accessToken string, teamId string) *Bot {
	return &Bot{
		db:     database,
		token:  accessToken,
		teamId: teamId,
	}
}

func (b *Bot) Run() {
	bot := slacker.NewClient(b.token, os.Getenv("SLACK_APP_TOKEN"))
	bingoMgr := bingo.NewBingoManager(b.db)

	bot.Command("bingo add <item>", &slacker.CommandDefinition{
		Description: "Add a new cell to the bingo board",
		Example:     "bingo add <cell text>",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			items := strings.Split(request.Param("item"), ";")
			bingo := bingoMgr.GetOrCreate(botCtx.Event().Channel)
			if bingo == nil {
				response.Reply("Failed to get bingo board")
				return
			}

			for _, item := range items {
				text := strings.TrimSpace(item)
				if text == "" {
					continue
				}

				if _, err := bingo.AddCell(text); err != nil {
					response.Reply(fmt.Sprintf("Failed to add item: %s", err))
					return
				}
			}

			response.Reply(bingo.ToString())
		},
	})

	bot.Command("bingo remove <id>", &slacker.CommandDefinition{
		Description: "Remove a cell from the bingo board",
		Example:     "bingo remove <cell id>",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			id := request.IntegerParam("id", 0)
			bingo := bingoMgr.GetOrCreate(botCtx.Event().Channel)
			if bingo == nil {
				response.Reply("Failed to get bingo board")
				return
			}

			if !bingo.RemoveCell(id) {
				response.Reply(fmt.Sprintf("Cell %d not found", id))
				return
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
			if bingo == nil {
				response.Reply("Failed to get bingo board")
				return
			}

			if !bingo.MarkCell(id) {
				response.Reply(fmt.Sprintf("Cell %d not found", id))
				return
			}

			if bingo.IsCompleted() {
				response.Reply(":bell: :bell: :bell: Bingo! You win! :tada: :tada: :tada:")
			}

			err := sendBingoBoard(bingo, botCtx)
			if err != nil {
				response.Reply(err.Error())
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
			if bingo == nil {
				response.Reply("Failed to get bingo board")
				return
			}

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
			if bingo == nil {
				response.Reply("Failed to get bingo board")
				return
			}

			bingo.Reset()
			err := sendBingoBoard(bingo, botCtx)
			if err != nil {
				response.Reply(err.Error())
			}
		},
	})

	bot.Command("bingo list", &slacker.CommandDefinition{
		Description: "Show the bingo items in a list",
		Example:     "bingo list",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			bingo := bingoMgr.GetOrCreate(botCtx.Event().Channel)
			if bingo == nil {
				response.Reply("Failed to get bingo board")
				return
			}

			response.Reply(bingo.ToString())
		},
	})

	bot.Command("bingo", &slacker.CommandDefinition{
		Description: "Show the bingo board",
		Example:     "bingo",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			channel := botCtx.Event().Channel
			bingo := bingoMgr.GetOrCreate(channel)
			if bingo == nil {
				response.Reply("Failed to get bingo board")
				return
			}

			if len(bingo.Cells) == 0 {
				response.Reply("No bingo items")
				return
			}

			err := sendBingoBoard(bingo, botCtx)
			if err != nil {
				response.Reply(err.Error())
			}
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Starting bingo bot instance fot team " + b.teamId)
	err := bot.Listen(ctx)
	if err != nil {
		fmt.Println("failed to start bot", err)
	}
}

func sendBingoBoard(bingo *bingo.Bingo, botCtx slacker.BotContext) error {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		visual.GenerateImage(bingo, w)
	}()

	_, err := botCtx.Client().UploadFile(slack.FileUploadParameters{
		Filetype: "png",
		Filename: "bingo.png",
		Title:    "Bingo",
		Channels: []string{botCtx.Event().Channel},
		Reader:   r,
	})

	return err
}
