# Slack bingo bot
Slack bot that I made while learning Go.

## Usage
- Add items to your bingo board using the commands bellow.
  For this to work correctly there must be a specific number of items - 1, 4, 9, 16, 25. This is needed to make a square bingo board.
- Mark the items as completed when  needed
- Use `bingo` command to visualize the bingo board

## Commands
- `bingo` - show bingo board
- `bingo list` - lists bingo items
- `bingo add <item>` - adds item to bingo board. You can add multiple items at once by splitting them with `;`
- `bingo remove <id>` - removes an item from bingo board
- `bingo mark <id>` - marks an item as completed
- `bingo switch <id1> <id2>` - sitches item <id1> with <id2>
- `bingo reset` - reset the bingo board

## Setup
- Follow the steps to setup slack app [from Slacker docs](https://github.com/shomali11/slacker#preparing-your-slack-app)
- Create `.env` file with variables: `SLACK_BOT_TOKEN` and `SLACK_APP_TOKEN`
  - `SLACK_CLIENT_ID`
  - `SLACK_CLIENT_SECRET`
  - `SLACK_REDIRECT_URL`
  - `SLACK_APP_TOKEN`
  - `SLACK_BOT_TOKEN` - optional, when bot is connected to single client and doesn't need oauth flow
  - `HOST` - optional, host address (127.0.0.1 for windows to prevent firewall popups)
  - `PORT` - optional, port for auth server, default - 3050
- `go run .`

## Docker info
- Mount a volume to `/usr/src/app/data` to save database file (this project is using SQLite)

## Limitations
Currently, only single workspace can be added per app. Need to rewrite everything without slack-go package or to not use sockets.

## TODO
- Home page with "Add to slack" button
- Bingo board name
  - Render name in the image