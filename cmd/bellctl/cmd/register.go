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
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		address, err := url.Parse(viper.GetString("bell.address") + RegisterPath)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"server address": viper.GetString("bell.address"),
				"method":         "RegisterCmd.Run",
			}).Error("Failed to build url")
			return
		}
		address.Scheme = "ws"
		c, _, err := websocket.DefaultDialer.Dial(address.String(), nil)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed create websocket connection to server")
			return
		}
		defer c.Close()

		// This code wait for order of music to play and then play them.
		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				mt, message, err := c.ReadMessage()
				if err != nil {
					if mt == websocket.CloseNormalClosure {
						return
					}
					logrus.WithError(err).Error("Failed to receive some message")
					return
				}
				logrus.Infof("sound to play: %v", string(message))
			}
		}()

		client := struct {
			Name string `json:"name"`
		}{
			Name: viper.GetString("register.name"),
		}
		// register the client
		err = c.WriteJSON(client)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to send message")
			return
		}

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		for {
			select {
			case <-done:
				return
			case <-interrupt:
				logrus.Infof("Interruption received")
				err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing"))
				if err != nil {
					logrus.WithError(err).Errorf("Failed to send normal closing message")
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringP("name", "n", "", "Name used to register the client. This name will be used as destination to play sounds by client")
	cn := "register.name"
	viper.BindPFlag(cn, registerCmd.Flags().Lookup("name"))
	viper.BindEnv(cn, "BELL_REGISTER_NAME")
	// default to hostname if no fail
	hn, err := os.Hostname()
	if err == nil {
		viper.SetDefault(cn, hn)
	}

}
