package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/restanrm/bell/sound"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
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
}

func exitIfNotSetted(key string) {
	s := viper.GetString(key)
	if s == "" {
		fmt.Printf("required variable %q is not setted\n", key)
		os.Exit(1)
	}
}

func main() {
	if !viper.GetBool("flite") {
		exitIfNotSetted("polly.accessKey")
		exitIfNotSetted("polly.secretKey")
	}

	sounds := sound.Load(viper.GetString("storefile"))

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/", listSounds(sounds))
	api.HandleFunc("/play/{sound:[-a-zA-Z]+}", soundPlayer(sounds))
	api.HandleFunc("/sounds", addSound(sounds)).Methods("POST")
	api.HandleFunc("/sounds", listSounds(sounds)).Methods("GET")
	api.HandleFunc("/sounds/{sound:[-a-zA-Z]+}", deleteSound(sounds)).Methods("DELETE")

	api.HandleFunc("/tts", ttsPostHandler()).Methods("POST")
	api.HandleFunc("/tts", ttsGetHandler()).Methods("GET")

	// Delegate all paths to front/dist
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./front/dist")))

	logrus.Info("Listening on address: ", viper.GetString("listen"))
	log.Fatal(http.ListenAndServe(viper.GetString("listen"), webLogger(r)))
}
