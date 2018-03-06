// Copyright © 2018 Adrien Raffin-Caboisse
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
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/restanrm/bell"
	"github.com/restanrm/bell/sound"
	_ "github.com/restanrm/bell/statik"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bell [OPTIONS]",
	Short: "Bell is the command used to render the bell service",
	Long: `Bell command can run a bell server or only the front interface, or both.
By default, both the front and the API are run on the same server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.GetBool("flite") {
			exitIfNotSetted("polly.accessKey")
			exitIfNotSetted("polly.secretKey")
		}

		r := mux.NewRouter()
		switch len(args) {
		case 0:
			prepareFront(r)
			prepareAPI(r)
			serve(r)
		case 1:
			switch args[0] {
			case "-server":
				prepareAPI(r)
				serve(r)
			case "-front":
				prepareFront(r)
				serve(r)
			default:
				logrus.Error("Options can only be [-server|-front]")
				return
			}
		default:
			logrus.Error("Too many arguments")
			return
		}

	},
}

func prepareAPI(r *mux.Router) {
	sounds := sound.New(viper.GetString("storefile"))

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/", bell.ListSounds(sounds))
	api.HandleFunc("/play/{sound:[-a-zA-Z]+}", bell.SoundPlayer(sounds))
	api.HandleFunc("/sounds", bell.AddSound(sounds)).Methods("POST")
	api.HandleFunc("/sounds", bell.ListSounds(sounds)).Methods("GET")
	api.HandleFunc("/sounds/{sound:[-a-zA-Z]+}", bell.DeleteSound(sounds)).Methods("DELETE")
	api.HandleFunc("/sounds/{sound:[-a-zA-Z]+}", bell.GetSound(sounds)).Methods("GET")

	api.HandleFunc("/tts", bell.TtsPostHandler()).Methods("POST")
	api.HandleFunc("/tts", bell.TtsGetHandler()).Methods("GET")

}
func prepareFront(r *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	if viper.GetBool("embed.front") {
		r.PathPrefix("/").Handler(http.FileServer(statikFS))
	} else {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("front/dist")))
	}
}
func serve(r *mux.Router) {
	logrus.Info("Listening on address: ", viper.GetString("listen"))
	log.Fatal(http.ListenAndServe(viper.GetString("listen"), cors.Default().Handler(bell.WebLogger(r))))
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.BindEnv("storefile", "STORE_FILE")
	viper.SetDefault("storefile", "store.json")
	viper.BindEnv("soundDir", "SOUND_DIR")
	viper.SetDefault("soundDir", "sounds")
	viper.BindEnv("listen", "LISTEN_ADDR")
	viper.SetDefault("listen", ":10101")
	viper.BindEnv("polly.accessKey", "POLLY_ACCESS_KEY")
	viper.BindEnv("polly.secretKey", "POLLY_SECRET_KEY")
	viper.BindEnv("flite", "FLITE")
	viper.SetDefault("flite", true)
	viper.BindEnv("polly.voice", "POLLY_VOICE")
	viper.SetDefault("polly.voice", "Amy")
	viper.BindEnv("embed.front", "EMBED_FRONT")
	viper.SetDefault("embed.front", true)
	viper.AutomaticEnv() // read in environment variables that match
}

func exitIfNotSetted(key string) {
	s := viper.GetString(key)
	if s == "" {
		fmt.Printf("required variable %q is not setted\n", key)
		os.Exit(1)
	}
}
