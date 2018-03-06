// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "retrieve sound and store it locally",
	Long:  `Retrieve sound from bell server and store it locally`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Error("We need a sound to retrieve")
			return
		}
		sound := args[0]
		address, err := url.Parse(viper.GetString("bell.address") + GetSoundPath + sound)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"server address": viper.GetString("bell.address"),
				"method":         "getCmd.Run",
			}).Error("Failed to build url")
			return
		}

		resp, err := http.Get(address.String())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("Failed to retrieve sound content")
			return
		}
		defer resp.Body.Close()

		output := cmd.Flag("output").Value.String()
		var w io.WriteCloser
		if output == "-" {
			w = os.Stdout
		} else {
			w, err = os.Create(output)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("Failed to open destination file")
				return
			}
		}
		defer w.Close()
		n, err := io.Copy(w, resp.Body)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":        err,
				"bytesWritten": n,
				"destination":  output,
			}).Error("Failed to copy from web to destination")
		}
		logrus.WithFields(logrus.Fields{"bytesWritten": n}).Debug("Written output to destination")
		err = w.Close()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to close the file")
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("output", "o", "-", "Filepath of where to save the file content")
}
