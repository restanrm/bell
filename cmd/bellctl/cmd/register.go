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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/restanrm/bell/connstore"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register allows to connect to websocket of bell server. It will receive play orders and run them with `mpv`.",
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
		if address.Scheme == "https" {
			address.Scheme = "wss"
		} else {
			address.Scheme = "ws"
		}
		c, r, err := websocket.DefaultDialer.Dial(address.String(), nil)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":    err,
				"response": fmt.Sprintf("%#v", r),
			}).Error("Failed create websocket connection to server")
			return
		}
		defer c.Close()

		// send name to server
		err = sendName(c, viper.GetString("register.name"))
		if err != nil {
			logrus.WithError(err).Errorf("Failed to send name to destination")
			return
		}

		// read the name that is actually used by the server
		name, err := readName(c)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to read name from the server")
			return
		}
		logrus.Infof("Client registered as %q", name)

		// start code to listen to play order or closing message
		done := make(chan struct{})
		go readOrder(c, done)

		// register on local closing request
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, os.Kill)
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

type ReadMessager interface {
	ReadMessage() (messageType int, p []byte, err error)
}

func readOrder(c ReadMessager, done chan struct{}) {
	defer close(done)
	dir, err := ioutil.TempDir("/tmp", "bellPlayer")
	if err != nil {
		logrus.Errorf("Failed to create temp dir to store the sounds")
		os.Exit(-1)
	}
	defer os.RemoveAll(dir)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				// normal close
				logrus.Info("Closing the websocket")
			} else {
				// not expected error
				logrus.WithError(err).Errorf("Failed to received some message")
			}
			return
		}
		s := &connstore.PlayerRequest{}
		json.Unmarshal(message, s)
		go func() {
			switch s.Type {
			case "error":
				logrus.Error(s.Data)
			case "tts":
				logrus.WithField("text", s.Data).Info("Received TTS order")
				err = getTTSAndPlay(dir, s.Data)
				if err != nil {
					logrus.WithError(err).Errorf("Failed to retrieve and tts %v", s.Data)
				}
			case "sound":
				logrus.WithField("sound", s.Data).Info("Received play sound order")
				err = getAndPlay(dir, s.Data)
				if err != nil {
					logrus.WithError(err).Errorf("Failed to play the sound: %v", s.Data)
				}
			}
		}()
	}
}

func getAndPlay(dir, sound string) error {
	fp := filepath.Join(dir, fmt.Sprintf("%v.mp3", sound))
	err := get(sound, fp)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve sound %v", sound)
	}
	return play(fp)
}

func getTTSAndPlay(dir, text string) error {
	// compute hash of the text to have a filename
	fp, err := getTTS(dir, text)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve the sound from the bell server")
	}
	return play(fp)
}

func play(fp string) error {
	cmd := exec.Command(
		"mpv",
		fp,
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Failed to run the command: %q", strings.Join(cmd.Args, " "))
	}
	return nil
}

type ReadJSONer interface {
	ReadJSON(interface{}) error
}

type WriteJSONer interface {
	WriteJSON(interface{}) error
}

func sendName(c WriteJSONer, name string) error {
	err := c.WriteJSON(connstore.RegisterRequest{Name: name})
	if err != nil {
		return errors.Wrapf(err, "Failed to send name to server")
	}
	return nil
}

func readName(c ReadJSONer) (name string, err error) {
	// received used name on server side and timer duration to send ping
	resp := &connstore.RegisterResponse{}
	err = c.ReadJSON(resp)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read response")
	}
	return resp.Name, nil
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
