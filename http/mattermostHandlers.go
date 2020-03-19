package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/restanrm/bell/tts"

	"github.com/restanrm/bell/player"
	"github.com/restanrm/bell/sound"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// Ephemeral response ar only seen by querier
	Ephemeral = "ephemeral"
	// InChannel responses are seen by all actors of the channel
	InChannel = "in_channel"
)

// SlashCommandResponse is the type of response object used in mattermost
type SlashCommandResponse struct {
	Text         string `json:"text,omitempty"`
	Type         string `json:"response_type,omitempty"`
	Username     string `json:"username,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	GotoLocation string `json:"goto_location,omitempty"`
	// Attachments struct not implemented
	// Type not implemented yet
	// props not implemented yet
}

type listSender interface {
	Lister
	Sender
}

// MattermostHandler handle mattermost /bell commands.
// it allows to list and play sounds, and do some TTS.
func MattermostHandler(vault sound.Sounder, listSender listSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rToken := r.FormValue("token")
		logrus.Debug("rToken: ", rToken)
		if token := viper.GetString("mattermost.token"); token != "" {
			if rToken != token {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		logrus.Debugf("mattermost request token is valid")

		// command := r.FormValue("command") // complete command ex: /bell
		text := r.FormValue("text") // complete list of arguments
		responseURL := r.FormValue("response_url")

		// parse command and build response to send back to caller
		response := parseCommand(vault, listSender, text)

		jres, err := json.Marshal(response)
		if err != nil {
			logrus.WithError(err).Error("Failed to encode response to json")
			return
		}

		logrus.Debug("response message: ", string(jres))
		_, err = http.Post(responseURL, "application/json", bytes.NewBuffer(jres))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"response_url": responseURL,
			}).Error("Failed to send response to mattermost server.")
		}

	}
}

func parseCommand(vault sound.Sounder, listSender listSender, text string) (response SlashCommandResponse) {

	response = SlashCommandResponse{
		Type: Ephemeral,
	}

	arguments := strings.Fields(text)

	if len(arguments) <= 0 {
		response.Text = "No subcommand specified, please use the following (list|play|say)"
		return
	}

	command, arguments := arguments[0], arguments[1:]
	switch command {
	case "list":
		switch {
		case len(arguments) == 1:
			if arguments[0] == "clients" || arguments[0] == "client" {
				clients := listSender.List()
				response.Text = formatClients(clients)
			}
		default:
			sounds := vault.GetSounds()
			response.Text = formatSounds(sounds)
		}
	case "play":
		m := new(player.MpvPlayer)
		switch {
		case len(arguments) <= 0:
			response.Text = "Cannot guess what sound to play"
		case len(arguments) == 1:
			err := vault.PlaySound(arguments[0], m)
			if err != nil {
				response.Text = fmt.Sprintf("Failed to play the sound: %v", err)
				return
			}
			response.Text = fmt.Sprintf(":musical_note: %q is playing :musical_note:", arguments[0])
			response.Type = InChannel
		case len(arguments) == 2 && arguments[0] == "-d":
			response.Text = fmt.Sprintf("destination or sound hasn't been specified. See documentation.")
		case len(arguments) > 2 && arguments[0] == "-d":
			// send playing sound to destination
			// arguments[0] should be "-d"
			// arguments[1] should be the destination "supervision"
			// arguments[2] should be the play ""
			destination := arguments[1]
			sound := arguments[2]
			err := PlayOnClient(listSender, destination, sound)
			if err != nil {
				response.Text = fmt.Sprintf("Failed to play %v on destination %v: %v",
					sound,
					destination,
					err,
				)
				return
			}
			response.Text = fmt.Sprintf(":musical_note: sound %q is playing on destination %q :musical_note:", sound, destination)
			response.Type = InChannel
		}
	case "say":
		var m player.Player
		var t tts.Sayer
		m = new(player.MpvPlayer)
		t = tts.NewTTS(
			viper.GetBool("flite"),
			viper.GetString("polly.accessKey"),
			viper.GetString("polly.secretKey"),
		)
		text := strings.Join(arguments, " ")
		err := t.Say(text, m)
		if err != nil {
			response.Text = fmt.Sprintf(":broken_heart: something went wrong: %s", err)
			return
		}
		response.Text = text
		response.Type = InChannel

	default:
		response.Text = "This command isn't supported"
	}
	return response
}

func formatSounds(sounds []sound.Sound) (out string) {
	if len(sounds) <= 0 {
		out += fmt.Sprintf("No sounds found")
		return
	}
	out += fmt.Sprintf("|%s|%s|\n", "Sound", "Tags")
	out += fmt.Sprintf("|:--|:--|\n")
	for _, sound := range sounds {
		out += fmt.Sprintf("|%s|%s|\n", sound.Name, strings.Join(sound.Tags, ","))
	}
	return out
}

func formatClients(clients []string) (out string) {
	if len(clients) <= 0 {
		out += fmt.Sprintf("No clients have been registered yet")
		return
	}
	out += fmt.Sprintf("|%s|\n", "Client")
	out += fmt.Sprintf("|:--|\n")
	for _, client := range clients {
		out += fmt.Sprintf("|%s|\n", client)
	}
	return out
}
