package bell

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/restanrm/bell/player"
	"github.com/restanrm/bell/sound"
	"github.com/spf13/viper"
)

// func midLogger(fn http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Print("client=", r.RemoteAddr, " URL=", r.URL.Path)
// 		fn(w, r)
// 	}
// }

func WebLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Print("client=", r.RemoteAddr, " URL=", r.URL.Path, " Method=", r.Method, " Params=", r.PostForm)
		h.ServeHTTP(w, r)
	})
}

var rxSound = regexp.MustCompile(`^[-a-zA-Z]+$`)

// SoundPlayer allow to play a sound from sounder service
func SoundPlayer(vault sound.Sounder) http.HandlerFunc {
	m := new(player.MpvPlayer)
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sound := vars["sound"]
		// validate sound name to regex
		if !rxSound.MatchString(sound) {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{"soundname": sound}).Warn("Client made a request with wrong input file name. It doesn't match the regexp")
			fmt.Fprintf(w, "Bad sound name. It doesn't match the regex '^[-a-ZA-Z]+$'")
			return
		}
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
			fmt.Fprintf(w, "Bad sound name. It doesn't match the regex '^[-a-ZA-Z]+$'")
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

		soundFilepath := fmt.Sprintf("/tmp/bell-sound-%v.mp3", uuid.New().String())
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
			fmt.Fprintf(w, "Bad sound name. It doesn't match the regex '^[-a-ZA-Z]+$'")
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
			fmt.Fprintf(w, "Bad sound name. It doesn't match the regex '^[-a-ZA-Z]+$'")
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

// TtsPostHandler handle request to play tts
func TtsPostHandler() http.HandlerFunc {
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

// TtsGetHandler handle request for TextToSpeech
func TtsGetHandler() http.HandlerFunc {
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
