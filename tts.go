package main

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type tts struct {
}

type Sayer interface {
	Say(string)
}

func (t *tts) Say(text string, player Player) error {
	fd, err := ioutil.TempFile("/tmp", "")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to create a temporary file")
		return err
	}
	fd.Close()
	filename := fd.Name() + ".wav"
	defer os.Remove(fd.Name())
	defer os.Remove(filename)

	cmd := exec.Command("flite",
		"-t", text,
		"-o", filename,
		"-voice", "awb",
	)

	out, err := cmd.Output()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"output":   out,
			"tempFile": filename,
		}).Error("Failed to read file")
		return err
	}

	player.PlayFilepath(filename)
	return nil
}
