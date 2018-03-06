package player

import (
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type MpvPlayer struct {
}

type Player interface {
	Play(filepath string) error
	PlayFilepath(string) error
}

func (mp *MpvPlayer) Play(path string) error {
	return mp.play(filepath.Join(viper.GetString("soundDir"), path))
}

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
