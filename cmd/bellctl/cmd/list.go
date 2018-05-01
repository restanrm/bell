package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/restanrm/bell/sound"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List available sounds to play",
	Long:    ``,
	Aliases: []string{`ls`},
	Run: func(cmd *cobra.Command, args []string) {
		address, err := url.Parse(viper.GetString("bell.address") + ListPath)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"server address": viper.GetString("bell.address"),
				"method":         "listCmd.Run",
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
		var sounds []sound.Sound
		json.NewDecoder(resp.Body).Decode(&sounds)
		if len(sounds) == 0 {
			logrus.Info("No sounds found")
			return
		}

		// extract list of tags and uniq them
		var tags = make(map[string]struct{})
		for _, sound := range sounds {
			if tagOption {
				for _, t := range sound.Tags {
					tags[t] = struct{}{}
				}
			}
		}

		if tagOption {
			fmt.Printf("List of tags\n")
			for k := range tags {
				fmt.Printf("  - %v\n", k)
			}
		} else {
			fmt.Println("List of sounds")
			for _, sound := range sounds {
				fmt.Printf("  - %v\n", sound.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&tagOption, "tag", "t", false, "Option to enable tag mode (list or play by tag)")
}
