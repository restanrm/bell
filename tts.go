package main

import (
	"crypto/md5"
	"fmt"
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

// createAudioFlite creates an audiofile with the system command "flite"
func (t *tts) createAudioFlite(text, filename string) error {
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
			"text":     text,
		}).Error("Failed to create speech from flite")
		return err
	}
	return nil
}

// createAudioPolly creates an audiofile with the polly API
func (t *tts) createAudioPolly(text, filename string) error {
	out, err := t.polly.Speech(text)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"text":     text,
			"filename": filename,
		}).Error("Failed to query polly to transform text to MP3")
		return err
	}

	// write speech to file
	err = ioutil.WriteFile(filename, out, 0644)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"filename": filename,
		}).Error("Failed to write temporary file")
		return err
	}
	return nil
}

func getHash(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}

// Say create a tempfile based on the choosen technology of TTS and order the player to
// play it on speaker
func (t *tts) Say(text string, player Player) error {
	var err error
	filepath := "/tmp/" + getHash(text)

	fliteFilename := filepath + ".wav"
	pollyFilename := filepath + ".mp3"

	// cached section
	// if polly file exist, play it
	if exist(pollyFilename) {
		player.PlayFilepath(pollyFilename)
		return nil
	}
	// if flite file exist
	if exist(fliteFilename) {
		// if flite it disabled
		if !t.flite {
			// try to create pollyfile
			err = t.createAudioPolly(text, pollyFilename)
			if err != nil {
				// if fail, play flitefile
				player.PlayFilepath(fliteFilename)
				return nil
			}
			// play new pollyfile
			player.PlayFilepath(pollyFilename)
			return nil
		}
		// flite is enabled and file exist, play it
		player.PlayFilepath(fliteFilename)
		return nil
	}

	// no cache, creation of the file
	if t.flite {
		err = t.createAudioFlite(text, fliteFilename)
		if err != nil {
			return err
		}
		player.PlayFilepath(fliteFilename)
		return nil
	} else {
		err = t.createAudioPolly(text, pollyFilename)
		if err != nil {
			err = t.createAudioFlite(text, fliteFilename)
			if err != nil {
				return err
			}
			player.PlayFilepath(fliteFilename)
			return nil
		}
		player.PlayFilepath(pollyFilename)
		return nil
	}

	return nil
}
