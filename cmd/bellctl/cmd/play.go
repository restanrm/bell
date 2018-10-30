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
	Short: "Play sound on the host that run the server command",
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

		q := address.Query()
		if tagOption {
			logrus.Warning("Tag option is deprecated. No need to use it to play a sound now.")
			q.Add("tag", "")
		}

		if viper.GetString("playOnClient") != "" {
			q.Add("destination", viper.GetString("playOnClient"))
		}
		address.RawQuery = q.Encode()

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

	playCmd.Flags().BoolVarP(&tagOption, "tag", "t", false, "Option to play a sound by its tag")
	playCmd.Flags().StringP("destination", "", "", "Destination to play the sound")
	viper.BindPFlag("playOnClient", playCmd.Flags().Lookup("destination"))

}
