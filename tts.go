package main

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/restanrm/golang-tts"
	"github.com/sirupsen/logrus"
)

type tts struct {
	polly *golang_tts.TTS
	flite bool
}

type Sayer interface {
	Say(string)
}

func NewTTS(flite bool, accessKey, secretKey string) *tts {
	polly := golang_tts.New(accessKey, secretKey)
	polly.Format(golang_tts.MP3)
	polly.Voice(golang_tts.Amy)
	return &tts{polly: polly, flite: flite}
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
	var filename string
	defer os.Remove(filename)
	defer os.Remove(fd.Name())

	if t.flite {
		filename = fd.Name() + ".wav"
		// create wav file from flite
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
	} else {
		filename = fd.Name() + ".mp3"
		// retrieve mp3 content of speech
		out, err := t.polly.Speech(text)
		if err != nil {
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
					"text":  text,
				}).Error("Failed to query polly to transform text to MP3")
			}
			return err
		}

		// write speech to file
		err = ioutil.WriteFile(filename, out, 0644)
		if err != nil {
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":    err,
					"filename": filename,
				}).Error("Failed to write temporary file")
			}
			return err
		}
	}

	// play speech with player
	player.PlayFilepath(filename)
	return nil
}
