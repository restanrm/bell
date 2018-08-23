package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
		sounds, err := list()
		if err != nil {
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
			// get max column sizes for sounds and tags
			ta, tb := "Sound", "Tags"
			mls, mlt := getMaxColumnSizes(sounds, ta, tb)
			printHeader(mls, mlt, ta, tb)
			for _, sound := range sounds {
				printLine(mls, mlt, sound.Name, strings.Join(sound.Tags, ","))
			}
			printFooter(mls, mlt)
		}
	},
}

func printHeader(mls, mlt int, ta, tb string) {
	printDashLine(mls, mlt)
	printLine(mls, mlt, ta, tb)
	printDashLine(mls, mlt)
}

func printFooter(mls, mlt int) {
	printDashLine(mls, mlt)
}

func printDashLine(mls, mlt int) {
	fmt.Printf("+-%[1]*[3]v-+-%[2]*[4]v-+\n", mls, mlt, strings.Repeat("-", mls), strings.Repeat("-", mlt))
}

func printLine(mls, mlt int, sound, tags string) {
	fmt.Printf("| %-[1]*[3]v | %-[2]*[4]v |\n", mls, mlt, sound, tags)
}

func getMaxColumnSizes(sounds []sound.Sound, ta, tb string) (int, int) {
	mls, mlt := 0, 0
	sounds = append(sounds, sound.Sound{Name: ta, Tags: []string{tb}})
	for _, sound := range sounds {
		if mls < len(sound.Name) {
			mls = len(sound.Name)
		}
		t := strings.Join(sound.Tags, ",")
		if mlt < len(t) {
			mlt = len(t)
		}
	}
	return mls, mlt
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&tagOption, "tag", "t", false, "Option to enable tag mode (list or play by tag)")
}

func list() (sounds []sound.Sound, err error) {
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
	json.NewDecoder(resp.Body).Decode(&sounds)
	if len(sounds) == 0 {
		err = errors.New("No sounds found")
	}
	return
}
