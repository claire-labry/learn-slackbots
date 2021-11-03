package controllers

import (
	"claire-labry/learn-slackbots/views"
	"log"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// We create a sctucture to let us use dependency injection
type AppHomeController struct {
	EventHandler *socketmode.SocketmodeHandler
}

// event handles the app home that you see when you click on an app and it pops up the app's view
func NewAppHomeController(eventhandler *socketmode.SocketmodeHandler) AppHomeController {
	c := AppHomeController{
		EventHandler: eventhandler,
	}

	// app home (2) when the user hits the event app home i.e. when the user hits 'claire test slackbot'
	c.EventHandler.HandleEventsAPI(
		slackevents.AppHomeOpened,
		c.publishHomeTabView,
	)
	// user triggers the create a stickie note
	c.EventHandler.HandleInteractionBlockAction(
		views.AddStickieNoteActionID,
		c.openCreateStickieNoteModal,
	)

	// create stickie note submitted to app
	c.EventHandler.HandleInteraction(
		slack.InteractionTypeViewSubmission,
		c.createStickieNote,
	)

	return c
}

// the event that calls the app to actually publish the app view when user clicks on it
func (c *AppHomeController) publishHomeTabView(evt *socketmode.Event, clt *socketmode.Client) {
	// we need to cast our socketmode.Event into slackevents.AppHomeOpenedEvent
	evt_api, _ := evt.Data.(slackevents.EventsAPIEvent)
	evt_app_home_opened, _ := evt_api.InnerEvent.Data.(slackevents.AppHomeOpenedEvent)

	// creates the view using the block-kit
	view := views.AppHomeTabView()

	// Publish the view
	// api client from `clt` and posts our view
	_, err := clt.GetApiClient().PublishView(evt_app_home_opened.User, view, "")

	// handles errors
	if err != nil {
		log.Printf("ERROR publishHomeTabView: %v", err)
	}
}

func (c *AppHomeController) openCreateStickieNoteModal(evt *socketmode.Event, clt *socketmode.Client) {
	// we need to cast our socketmode.Event
	interaction := evt.Data.(slack.InteractionCallback)

	// makes sure that the response to the server is received to avoid any errors
	clt.Ack(*evt.Request)

	// create a view using block-kit
	view := views.CreateStickieNoteModal()

	// actually opens the modal
	_, err := clt.GetApiClient().OpenView(interaction.TriggerID, view)

	if err != nil {
		log.Printf("ERROR openCreateStickieNoteModal: %v", err)
	}
}

func (c *AppHomeController) createStickieNote(evt *socketmode.Event, clt *socketmode.Client) {
	view_submission := evt.Data.(slack.InteractionCallback)

	clt.Ack(*evt.Request)

	note := views.StickieNote{
		Description: view_submission.View.State.Values[views.ModalDescriptionBlockID][views.ModalDescriptionActionID].Value,
		Color:       view_submission.View.State.Values[views.ModalColorBlockID][views.ModalColorActionID].SelectedOption.Value,
		Timestamp:   time.Unix(time.Now().Unix(), 0).String(),
	}

	view := views.AppHomeCreateStickieNote(note)

	_, err := clt.GetApiClient().PublishView(view_submission.User.ID, view, "")

	if err != nil {
		log.Printf("ERROR createStickieNote: %v", err)
	}
}
