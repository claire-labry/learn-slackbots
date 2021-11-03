package main

import (
	"claire-labry/learn-slackbots/controllers"
	"claire-labry/learn-slackbots/drivers"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	// this uses the zerolog module to prettify the output if an error is thrown in terminal
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// if there's a test_slack.env available, this will run the file thanks to godotenv,
	// if not, an error will be outputted by the logger
	err := godotenv.Load("./test_slack.env")
	if err != nil {
		log.Fatal().Msg("Error loading .env file!")
	}
	// listens for the client to connect to slack via socketmode thru drivers module
	client, err := drivers.ConnectToSlackViaSocketMode()
	if err != nil {
		log.Error().
			Str("error", err.Error()).
			Msg("Unable to connect to Slack!")

		os.Exit(1)
	}
	// Builds a Slack App Home in Golang using Socket Mode
	socketmodeHandler := socketmode.NewsSocketmodeHandler(client)

	controllers.NewAppHomeController(socketmodeHandler)

	socketmodeHandler.RunEventLoop()
}
