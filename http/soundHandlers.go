package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/restanrm/bell/player"
	"github.com/restanrm/bell/sound"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var rxSound = regexp.MustCompile(`^[-a-zA-Z0-9]+$`)

// WebLogger return log about the http queries
func WebLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		defer func(begin time.Time) {
			logrus.WithFields(logrus.Fields{
				"client": r.RemoteAddr,
				"URL":    r.URL.Path,
				"Method": r.Method,
				"Params": r.PostForm,
				"took":   time.Since(begin),
			}).Debug("HTTP informations")
		}(time.Now())
		h.ServeHTTP(w, r)
	})
}

// SoundPlayer allow to play a sound from sounder service
func SoundPlayer(vault sound.Sounder, sender Sender) http.HandlerFunc {
	m := new(player.MpvPlayer)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sound := vars["sound"]
		// validate sound name to regex
		if !rxSound.MatchString(sound) {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{"soundname": sound}).Warn("Client made a request with wrong sound or tagname. It doesn't match the regexp")
			fmt.Fprintf(w, "Bad sound or tag name. It doesn't match the regex %q", rxSound.String())
			return
		}
		var err error
		if dest, ok := r.URL.Query()["destination"]; ok {
			logrus.WithFields(logrus.Fields{
				"destination": dest[0],
				"sound":       sound,
			}).Infof("Sending play order to registerd client")
			PlayOnClient(sender, dest[0], sound)
		} else {
			err = vault.PlaySound(sound, m)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				logrus.WithFields(logrus.Fields{
					"name": sound,
				}).Info("Sound or tag has not been found in store")
				return
			}
		}
	}
}

// AddNewSound to Sounder service
func AddSound(vault sound.Sounder) http.HandlerFunc {
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
		// validate sound name to regex
		if !rxSound.MatchString(soundName) {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{"soundname": soundName}).Warn("Client made a request with wrong input file name. It doesn't match the regexp")
			fmt.Fprintf(w, "Bad sound name. It doesn't match the regex %q", rxSound.String())
			return
		}

		// retrieve tags and validate them
		var soundTags []string
		values := r.URL.Query()
		if tags, ok := values["tag"]; ok {
			soundTags = tags
		}
		// validate tag name to regex also
		for _, t := range soundTags {
			if !rxSound.MatchString(t) {
				logrus.WithFields(logrus.Fields{"tag": t}).Warn("Client made a request with wrong tag name. It doesn't match the regexp")
				fmt.Fprintf(w, "Bad tag name. It doesn't match the regex %q", rxSound.String())
				return
			}
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

		soundFilepath := filepath.Join(viper.GetString("soundDir"), fmt.Sprintf("%v.mp3", soundName))
		err = dirExist(soundFilepath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithError(err).Errorf("Directory to store sound doesn't exist")
			return
		}
		f, err := os.OpenFile(soundFilepath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to open file to save new sound")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		err = vault.CreateSound(soundName, soundFilepath, soundTags...)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to add new sound")
		}
	}
}

func dirExist(filename string) error {
	dir := filepath.Dir(filename)
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			// error is not known
			return errors.Wrapf(err, "Couldn't get stat informations on path: %v", dir)
		}
		// path doesn't exist, creating it
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.Wrapf(err, "Failed to create directory to store sound files")
		}
	}
	return nil
}

// DeleteSound allows to delete sound from library
func DeleteSound(vault sound.Sounder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		soundName := vars["sound"]
		if soundName == "" {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{}).Error("Missing sound name from query")
			fmt.Fprintf(w, "Missing sound name from query")
			return
		}
		// validate sound name to regex
		if !rxSound.MatchString(soundName) {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{"soundname": soundName}).Warn("Client made a request with wrong input file name. It doesn't match the regexp")
			fmt.Fprintf(w, "Bad sound name. It doesn't match the regex %q", rxSound.String())
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

// ListSounds function
func ListSounds(vault sound.Sounder) http.HandlerFunc {
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

// GetSound allows to retrieve asound from the server
func GetSound(vault sound.Sounder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		soundName := vars["sound"]
		if soundName == "" {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{}).Error("Missing sound name from query")
			fmt.Fprintf(w, "Missing sound name from query")
			return
		}
		// validate sound name to regex
		if !rxSound.MatchString(soundName) {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{"soundname": soundName}).Warn("Client made a request with wrong input file name. It doesn't match the regexp")
			fmt.Fprintf(w, "Bad sound name. It doesn't match the regex %q", rxSound.String())
			return
		}
		content, err := vault.GetSound(soundName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			logrus.WithFields(logrus.Fields{"soundName": soundName}).Error("Couldn't find sound name from query")
			http.Error(w, "Failed to find the requested file", http.StatusNotFound)
			return
		}
		w.Header().Add("ContentType", "audio/mpeg3")
		_, err = w.Write(content)
		if err != nil {
			logrus.WithField("err", err).Error("Couldn' write file content the responseWriter")
			return
		}
	}
}

// part for TextToSpeech
