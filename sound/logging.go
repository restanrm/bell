package sound

import (
	"time"

	"github.com/restanrm/bell/player"
	"github.com/sirupsen/logrus"
)

type loggingSound struct {
	Sounder
}

var _ Sounder = &loggingSound{}

// NewLoggingSound is the logging implementation of a sounder
func NewLoggingSound(s Sounder) Sounder {
	return &loggingSound{s}
}

func (l *loggingSound) CreateSound(name, filepath string, tags ...string) error {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"method":   "CreateSound",
			"name":     name,
			"filepath": filepath,
			"tags":     tags,
			"took":     time.Since(begin),
		}).Info("sound service query")
	}(time.Now())
	return l.Sounder.CreateSound(name, filepath)
}

func (l *loggingSound) UpdateSound(sound Sound) error {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"method": "UpdateSound",
			"sound":  sound,
			"took":   time.Since(begin),
		}).Info("sound service query")
	}(time.Now())
	return l.Sounder.UpdateSound(sound)
}

func (l *loggingSound) DeleteSound(name string) error {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"method": "DeleteSound",
			"name":   name,
			"took":   time.Since(begin),
		}).Info("sound service query")
	}(time.Now())
	return l.Sounder.DeleteSound(name)
}

func (l *loggingSound) PlaySound(name string, player player.Player) error {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"method": "PlaySound",
			"name":   name,
			"took":   time.Since(begin),
		}).Info("sound service query")
	}(time.Now())
	return l.Sounder.PlaySound(name, player)
}

func (l *loggingSound) PlaySoundByTag(tag string, player player.Player) error {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"method": "PlaySoundByTag",
			"tag":    tag,
			"took":   time.Since(begin),
		}).Info("sound service query")
	}(time.Now())
	return l.PlaySoundByTag(tag, player)
}

func (l *loggingSound) GetSound(name string) ([]byte, error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"method": "GetSound",
			"name":   name,
			"took":   time.Since(begin),
		}).Info("sound service query")
	}(time.Now())
	return l.Sounder.GetSound(name)
}

func (l *loggingSound) GetSounds() []Sound {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"method": "GetSounds",
			"took":   time.Since(begin),
		}).Info("sound service query")
	}(time.Now())
	return l.Sounder.GetSounds()
}
