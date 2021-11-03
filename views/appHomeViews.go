package views

import (
	"bytes"
	"embed"
	"encoding/json"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/slack-go/slack"
)

const (
	AddStickieNoteActionID   = "add_note"
	ModalDescriptionBlockID  = "note_description"
	ModalDescriptionActionID = "content"
	ModalColorBlockID        = "note_color"
	ModalColorActionID       = "color"
)

type StickieNote struct {
	Description string
	Color       string
	Timestamp   string
}

var appHomeAssets embed.FS

func AppHomeTabView() slack.HomeTabViewRequest {
	str, err := appHomeAssets.ReadFile("appHomeViewsAssets/AppHomeView.json")
	if err != nil {
		log.Printf("Unable to read view `AppHomeView` : %v", err)
	}

	view := slack.HomeTabViewRequest{}
	json.Unmarshal([]byte(str), &view)

	return view
}

func CreateStickieNoteModal() slack.ModalViewRequest {

	str, err := appHomeAssets.ReadFile("appHomeViewsAssets/CreateStickieNoteModal.json")
	if err != nil {
		log.Printf("Unable to read view `CreateStickieNoteModal` : %v", err)
	}
	view := slack.ModalViewRequest{}
	json.Unmarshal([]byte(str), &view)

	return view
}

func AppHomeCreateStickieNote(note StickieNote) slack.HomeTabViewRequest {

	str, err := appHomeAssets.ReadFile("appHomeViewsAssets/AppHomeView.json")
	if err != nil {
		log.Printf("Unable to read view `AppHomeView` : %v", err)
	}
	view := slack.HomeTabViewRequest{}
	json.Unmarshal(str, &view)

	t, err := template.ParseFS(appHomeAssets, "appHomeViewsAssets/NoteBlock.json")
	if err != nil {
		log.Printf("Unable to read view `AppHomeView` : %v", err)
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, note)
	if err != nil {
		panic(err)
	}

	str, _ = ioutil.ReadAll(&tpl)
	note_view := slack.HomeTabViewRequest{}
	json.Unmarshal(str, &note_view)

	view.Blocks.BlockSet = append(view.Blocks.BlockSet, note_view.Blocks.BlockSet...)

	return view
}
