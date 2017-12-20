package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func viperConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/data/")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to read config file")
		os.Exit(-1)
	}
}

func main() {
	viperConfig()
	var sounds Sounds

	sounds.Load(viper.GetString("storefile"))

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	//api.HandleFunc("/play", soundPlayer)
	api.HandleFunc("/", listSounds(sounds))
	api.HandleFunc("/play/{sound:[-a-zA-Z]+}", soundPlayer(sounds))
	api.HandleFunc("/tts", ttsPostHandler()).Methods("POST")
	api.HandleFunc("/tts", ttsGetHandler()).Methods("GET")

	logrus.Info("Listening on address: ", viper.GetString("listen"))
	log.Fatal(http.ListenAndServe(viper.GetString("listen"), webLogger(r)))
}
