package http

import (
	"html/template"
	"net/http"

	"github.com/restanrm/bell/player"
	"github.com/restanrm/bell/tts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TtsPostHandler handle request to play tts
func TtsPostHandler(sender Sender) http.HandlerFunc {
	var t tts.Sayer
	t = tts.NewTTS(
		viper.GetBool("flite"),
		viper.GetString("polly.accessKey"),
		viper.GetString("polly.secretKey"),
	)
	t = tts.NewLoggingService(t)
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
		if dest, ok := r.URL.Query()["destination"]; ok {
			logrus.WithFields(logrus.Fields{
				"destination": dest[0],
				"sound":       text,
			}).Infof("Sending text to speech order to registered client")
			SayOnClient(sender, dest[0], text)
		} else {
			t.Say(text, m)
		}
	}
}

// TtsPostHandler handle request to play tts
func TtsGetPostHandler() http.HandlerFunc {
	var t tts.Sayer
	t = tts.NewTTS(
		viper.GetBool("flite"),
		viper.GetString("polly.accessKey"),
		viper.GetString("polly.secretKey"),
	)
	t = tts.NewLoggingService(t)
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
		content, err := t.GetSay(text)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			logrus.WithFields(logrus.Fields{"text": string(content)}).Error("Couldn't retrieve text to speech")
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
