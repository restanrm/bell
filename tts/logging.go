package tts

import (
	"time"

	"github.com/restanrm/bell/player"
	"github.com/sirupsen/logrus"
)

type loggingService struct {
	Sayer
}

// NewLoggingService returns a new logging service to log calls to say method
func NewLoggingService(s Sayer) Sayer {
	return &loggingService{s}
}

func (l *loggingService) Say(text string, p player.Player) error {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"methods": "Say",
			"text":    text,
			"took":    time.Since(begin),
		}).Info("logging tts service query")
	}(time.Now())
	return l.Sayer.Say(text, p)
}

func (l *loggingService) GetSay(text string) ([]byte, error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"methods": "GetSay",
			"text":    text,
			"took":    time.Since(begin),
		}).Info("Retrieve text to speech")
	}(time.Now())
	return l.Sayer.GetSay(text)
}
