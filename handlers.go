package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/restanrm/bell/player"
	"github.com/restanrm/bell/sound"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func midLogger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("client=", r.RemoteAddr, " URL=", r.URL.Path)
		fn(w, r)
	}
}

func webLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Print("client=", r.RemoteAddr, " URL=", r.URL.Path, " Params=", r.PostForm)
		h.ServeHTTP(w, r)
	})
}

func soundPlayer(vault sound.Sounds) http.HandlerFunc {
	m := new(player.MpvPlayer)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sound := vars["sound"]
		for _, s := range vault {
			if s.Name == sound {
				logrus.WithFields(logrus.Fields{
					"sound": s,
				}).Debug("Sound have been found, playing now…")

				m.Play(s.FileName)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		logrus.WithFields(logrus.Fields{
			"Sound": sound,
		}).Info("Sound has not been found on sound store")
	}
}

func addNewSound(w http.ResponseWriter, r *http.Request) {
	// this function add new sound file on sound dir path
	fmt.Fprintf(w, "Not implemented yet")
}

func listSounds(vault sound.Sounds) http.HandlerFunc {
	// this function list all currently available sounds
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		b, err := json.Marshal(vault)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":  err,
				"sounds": vault,
			}).Error("Failed to encode sounds to json")
			http.Error(w, "Failed to encode json object", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"b":     string(b),
			}).Error("Failed to write sound vault to web connection")
			http.Error(w, "Failed to return content to user", http.StatusInternalServerError)
			return
		}
	}
}

// part for TextToSpeech

func ttsPostHandler() http.HandlerFunc {
	var tts = NewTTS(
		viper.GetBool("flite"),
		viper.GetString("polly.accessKey"),
		viper.GetString("polly.secretKey"),
	)
	var m = &player.MpvPlayer{}
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		texts, ok := r.PostForm["text"]
		if !ok {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		var text = "Please give some content to text variable via POST form"
		if len(texts) >= 1 {
			text = texts[0]
		}
		tts.Say(text, m)
	}
}

func ttsGetHandler() http.HandlerFunc {
	pattern := `
<!doctype html>
	<head></head>
	<body>
		<div>
			<form method="POST">
				<label for="text">Text to say</label>
				<input type="text" name="text" id="text" size="75"/>
				<input type="submit" value="Send" />
			</form>
		</div>
	</body>
</html>
	`
	tmpl, err := template.New("ttsPost").Parse(pattern)
	//tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to load template")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		err = tmpl.Execute(w, nil)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to write template to client")
		}
	}
}
