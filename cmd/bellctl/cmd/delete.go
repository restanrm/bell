package cmd

import (
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{`rm`, `del`},
	Short:   "delete allows to remove sounds from library",
	Long:    `This endpoint allows to remove sounds from library. To reset library to base just reboot the service.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Error("Failed to play nothing as a sound")
			return
		}
		sound := args[0]
		address, err := url.Parse(viper.GetString("bell.address") + DeleteSoundPath + sound)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"server address": viper.GetString("bell.address"),
				"method":         "deletecCmd.Run",
			}).Error("Failed to build url")
			return
		}
		req, err := http.NewRequest(http.MethodDelete, address.String(), nil)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to create request to delete resource")
			return
		}
		// perform request on remote server
		resp, err := http.DefaultClient.Do(req)
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
	rootCmd.AddCommand(deleteCmd)

}
