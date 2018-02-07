package cmd

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sayCmd represents the say command
var sayCmd = &cobra.Command{
	Use:   "say",
	Short: "say target use tts to say what you wrote",
	Long: `say use a string argument to play sound to your bell server.append

	ex:
		bellctl say hello world
		bellctl say "Why does the skeleton dances alone ? Because he has nobody."
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var text string
		switch {
		case len(args) < 1:
			logrus.Error("Failed to say nothing. Please say something...")
			return
		case len(args) == 1:
			text = args[0]
		case len(args) > 1:
			text = strings.Join(args, " ")
		}
		address, err := url.Parse(viper.GetString("bell.address") + TtsPath)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"server address": viper.GetString("bell.address"),
				"method":         "PlayCmd.Run",
			}).Error("Failed to build url")
		}
		resp, err := http.PostForm(address.String(), url.Values{"text": {text}})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to contact bell server")
		}
		if resp.StatusCode > 299 {
			logrus.WithFields(logrus.Fields{
				"text":        text,
				"status_code": resp.StatusCode,
			}).Info("Failed to say what you wanted")
		}
	},
}

func init() {
	rootCmd.AddCommand(sayCmd)
}
