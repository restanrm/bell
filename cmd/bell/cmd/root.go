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
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/restanrm/bell/connstore"
	localHttp "github.com/restanrm/bell/http"
	"github.com/restanrm/bell/metrics"
	"github.com/restanrm/bell/sound"
	_ "github.com/restanrm/bell/statik"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
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
		viper.SetDefault("soundDir", filepath.Join(viper.GetString("dataDir"), "sounds"))
		viper.SetDefault("TTSDir", filepath.Join(viper.GetString("dataDir"), "tts"))
		if !viper.GetBool("flite") {
			exitIfNotSetted("polly.accessKey")
			exitIfNotSetted("polly.secretKey")
		}

		r := mux.NewRouter()
		m := r.PathPrefix("/").Subrouter()
		m.Handle("/metrics", promhttp.Handler())

		if serverOptions.api {
			prepareAPI(r)
			serve(r)
		}

		if serverOptions.front {
			prepareFront(r)
			serve(r)
		}

		prepareAPI(r)
		prepareFront(r)
		serve(r)

	},
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

	viper.BindEnv("listen", "LISTEN_ADDR")
	viper.SetDefault("listen", ":10101")
	viper.BindEnv("polly.accessKey", "POLLY_ACCESS_KEY")
	viper.BindEnv("polly.secretKey", "POLLY_SECRET_KEY")
	viper.SetDefault("flite", true)
	viper.BindEnv("polly.voice", "POLLY_VOICE")
	viper.SetDefault("polly.voice", "Amy")
	viper.SetDefault("embed.front", true)
	viper.SetDefault("verbose", false)
	viper.BindEnv("mattermost.token", "MATTERMOST_SLASH_TOKEN")

	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	rootCmd.Flags().BoolVarP(&serverOptions.api, "api", "a", false, "Allows to run the api as standalone service")
	rootCmd.Flags().BoolVarP(&serverOptions.front, "front", "f", false, "Allows to run the front separatly from the backend")

	rootCmd.Flags().StringP("dataDir", "d", "data", "Directory where all sounds and tts sounds are stored. The configuration file should also be located there.")
	viper.BindPFlag("dataDir", rootCmd.Flags().Lookup("dataDir"))

	rootCmd.Flags().StringP("config", "c", "store.json", "Configuration file where description of the sounds are stored")
	viper.BindPFlag("storefile", rootCmd.Flags().Lookup("config"))

	rootCmd.Flags().Bool("disable-websocket-checkorigin", false, "Disable the check of the origin for the websockets")
	viper.BindPFlag("websocket.checkorigin.disabled", rootCmd.Flags().Lookup("disable-websocket-checkorigin"))

	viper.AutomaticEnv() // read in environment variables that match

}

var serverOptions struct {
	front bool
	api   bool
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func exitIfNotSetted(key string) {
	s := viper.GetString(key)
	if s == "" {
		fmt.Printf("required variable %q is not setted\n", key)
		os.Exit(1)
	}
}

func prepareAPI(r *mux.Router) {
	var sounds sound.Sounder
	sounds = sound.New(filepath.Join(viper.GetString("dataDir"), viper.GetString("storefile")))
	sounds = sound.NewLoggingSound(sounds)

	api := r.PathPrefix("/api/v1").Subrouter()

	cs := connstore.New()

	// register metrics endpoint

	api.HandleFunc("/", instProm("root", localHttp.ListSounds(sounds)))
	api.HandleFunc("/play/{sound:[-a-zA-Z0-9]+}", instProm("play", localHttp.SoundPlayer(sounds, cs)))
	api.HandleFunc("/sounds", instProm("add", localHttp.AddSound(sounds))).Methods("POST")
	api.HandleFunc("/sounds", instProm("list", localHttp.ListSounds(sounds))).Methods("GET")
	api.HandleFunc("/sounds/{sound:[-a-zA-Z0-9]+}", instProm("delete", localHttp.DeleteSound(sounds))).Methods("DELETE")
	api.HandleFunc("/sounds/{sound:[-a-zA-Z0-9]+}", instProm("get", localHttp.GetSound(sounds))).Methods("GET")

	api.HandleFunc("/tts", instProm("say", localHttp.TtsPostHandler(cs))).Methods("POST")
	api.HandleFunc("/tts/retrieve", instProm("getsay", localHttp.TtsGetPostHandler())).Methods("POST")
	api.HandleFunc("/tts", instProm("sayform", localHttp.TtsGetHandler())).Methods("GET")

	api.HandleFunc("/mattermost", instProm("mattermost", localHttp.MattermostHandler(sounds, cs))).Methods("POST")

	// websocket handler
	api.HandleFunc("/clients", instProm("connStoreList", localHttp.ListClients(cs))).Methods("Get")
	api.HandleFunc("/clients/register", instProm("connStoreRegister", localHttp.RegisterClients(cs))).Methods("GET")

}

func prepareFront(r *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	f := r.PathPrefix("/").Subrouter()

	if viper.GetBool("embed.front") {
		f.PathPrefix("/").HandlerFunc(instProm("front", http.FileServer(statikFS).ServeHTTP))
	} else {
		f.PathPrefix("/").HandlerFunc(instProm("front", http.FileServer(http.Dir("front/dist")).ServeHTTP))
	}
}

func instProm(label string, f http.HandlerFunc) http.HandlerFunc {
	return promhttp.InstrumentHandlerDuration(
		metrics.HTTPRequestDuration.MustCurryWith(
			prometheus.Labels{"handler": label}),
		promhttp.InstrumentHandlerCounter(
			metrics.HTTPRequestsCount.MustCurryWith(
				prometheus.Labels{"handler": label}),
			f,
		),
	)
}

func serve(r *mux.Router) {
	logrus.Info("Listening on address: ", viper.GetString("listen"))
	log.Fatal(http.ListenAndServe(viper.GetString("listen"), cors.Default().Handler(localHttp.WebLogger(r))))

}
