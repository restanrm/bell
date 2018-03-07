// Package player is the package that describe what a player is and
// it propose an implementation of this interface with mpv player
package player

import (
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MpvPlayer struct {
}

// Middleware is the type that allows to chain Player objects
type Middleware func(Player) Player

// Player is the interface to allow any implementation of the player service.
type Player interface {
	Play(filepath string) error
	PlayFilepath(string) error
}

// Play is the function used to play sounds
// it concatenate the given path with the filepath of the application
func (mp *MpvPlayer) Play(path string) error {
	return mp.play(filepath.Join(viper.GetString("soundDir"), path))
}

// PlayFilepath is a function to play a file given a filepath
func (mp *MpvPlayer) PlayFilepath(fp string) error {
	return mp.play(fp)
}

func (mp *MpvPlayer) play(fp string) error {
	cmd := exec.Command(
		"mpv",
		"--audio-normalize-downmix=yes",
		fp,
	)

	out, err := cmd.Output()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"output":   string(out),
			"filepath": fp,
		}).Error("Failed to read file")
		return err
	}
	return nil

}
