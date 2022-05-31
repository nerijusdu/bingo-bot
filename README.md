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
- `bingo add <item>` - adds item to bingo board
- `bingo remove <id>` - removes an item from bingo board
- `bingo mark <id>` - marks an item as completed
- `bingo switch <id1> <id2>` - sitches item <id1> with <id2>
- `bingo reset` - reset the bingo board

## Setup
- Follow the steps to setup slack app [from Slacker docs](https://github.com/shomali11/slacker#preparing-your-slack-app)
- Create `.env` file with `SLACK_BOT_TOKEN` and `SLACK_APP_TOKEN`
- `go run .`