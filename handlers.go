package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
		log.Print("client=", r.RemoteAddr, " URL=", r.URL.Path, " Method=", r.Method, " Params=", r.PostForm)
		h.ServeHTTP(w, r)
	})
}

// soundPlayer allow to play a sound from sounder service
func soundPlayer(vault sound.Sounder) http.HandlerFunc {
	m := new(player.MpvPlayer)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sound := vars["sound"]
		err := vault.PlaySound(sound, m)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			logrus.WithFields(logrus.Fields{
				"name": sound,
			}).Info("Sound has not been found in store")
			return
		}
		logrus.WithFields(logrus.Fields{
			"sound": sound,
		}).Debug("Sound have been found, playing now…")
	}
}

// addNewSound to Sounder service
func addSound(vault sound.Sounder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// this function add new sound file on sound dir path
		r.ParseMultipartForm(int64(1 * 1024 * 1024))
		soundName := r.FormValue("name")
		if soundName == "" {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{}).Error("Missing \"name\" field")
			fmt.Fprintf(w, "Missing \"name\" field")
			return
		}

		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Bad request from client")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Missing \"uploadFile\" field")
			return
		}
		defer file.Close()
		soundFilepath := fmt.Sprintf("/tmp/bell-sound-%v.mp3", time.Now().Unix())
		f, err := os.OpenFile(soundFilepath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to open file to save new sound")
		}
		defer f.Close()
		io.Copy(f, file)

		err = vault.CreateSound(soundName, soundFilepath)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to add new sound")
		}
	}
}

// deleteSound allows to delete sound from library
func deleteSound(vault sound.Sounder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		soundName := vars["sound"]
		if soundName == "" {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{}).Error("Missing sound name from query")
			fmt.Fprintf(w, "Missing sound name from query")
			return
		}
		err := vault.DeleteSound(soundName)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to remove sound from store")
		}
	}
}

// listSounds function
func listSounds(vault sound.Sounder) http.HandlerFunc {
	// this function list all currently available sounds
	return func(w http.ResponseWriter, r *http.Request) {
		sounds := vault.GetSounds()
		w.Header().Add("Content-Type", "application/json")

		b, err := json.Marshal(sounds)
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
