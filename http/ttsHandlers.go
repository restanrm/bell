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
func TtsPostHandler() http.HandlerFunc {
	var tts = tts.NewTTS(
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
