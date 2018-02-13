package cmd

import (
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Error("Failed to play nothing as a sound")
			return
		}
		sound := args[0]
		address, err := url.Parse(viper.GetString("bell.address") + PlayPath + sound)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"server address": viper.GetString("bell.address"),
				"method":         "PlayCmd.Run",
			}).Error("Failed to build url")
			return
		}
		resp, err := http.Get(address.String())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to contact bell server")
			return
		}
		if resp.StatusCode > 299 {
			logrus.WithFields(logrus.Fields{
				"sound":       sound,
				"status_code": resp.StatusCode,
			}).Info("Failed to play the sound")
		}
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

}