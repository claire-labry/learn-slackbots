package drivers

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func ConnectToSlackViaSocketMode() (*socketmode.Client, error) {
	// validates if there's a SLACK_APP_TOKEN when test_slack.env is run; and sets SLACK_APP_TOKEN to appToken
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		return nil, errors.New("SLACK_APP_TOKEN must be set!")
	}
	// if the SLACK_APP_TOKEN is not prefixed with xapp- then it throws this error
	if !strings.HasPrefix(appToken, "xapp-") {
		return nil, errors.New("SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}
	// validates if there's a SLACK_BOT_TOKEN string when test_slack.env is run; and sets SLACK_BOT_TOKEN to botToken
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		return nil, errors.New("SLACK_BOT_TOKEN must be set!")
	}
	// if the SLACK_BOT_TOKEN string is not prefixed with xoxb- then it throws this error
	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(appToken),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	return client, nil
}
