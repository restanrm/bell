package tts

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/restanrm/bell/player"
	"github.com/restanrm/golang-tts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type tts struct {
	polly *golang_tts.TTS
	flite bool
}

var _ Sayer = &tts{}

// Sayer is the interface to transform text to sound
type Sayer interface {
	Say(string, player.Player) error
	GetSay(string) ([]byte, error)
}

// NewTTS is the function that returns a *tts object. This object implement the
// Sayer interface
func NewTTS(flite bool, accessKey, secretKey string) *tts {
	polly := golang_tts.New(accessKey, secretKey)
	polly.Format(golang_tts.MP3)
	polly.Voice(viper.GetString("polly.voice"))
	return &tts{polly: polly, flite: flite}
}

func dirExist(filename string) error {
	dir := filepath.Dir(filename)
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			// error is not known
			return errors.Wrapf(err, "Couldn't get stat informations on path: %v", dir)
		}
		// path doesn't exist, creating it
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.Wrapf(err, "Failed to create directory to store TTS files")
		}
	}
	return nil
}

// createAudioFlite creates an audiofile with the system command "flite"
func (t *tts) createAudioFlite(text, filename string) error {
	err := dirExist(filename)
	if err != nil {
		return err
	}
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
	err := dirExist(filename)
	if err != nil {
		return err
	}
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

func (t *tts) getSayFilename(text string) (string, error) {
	var err error
	filepath := filepath.Join(viper.GetString("TTSDir"), getHash(text))

	fliteFilename := filepath + ".wav"
	pollyFilename := filepath + ".mp3"

	// cached section
	// if polly file exist, play it
	if exist(pollyFilename) {
		return pollyFilename, nil
	}
	// if flite file exist
	if exist(fliteFilename) {
		// if flite it disabled
		if !t.flite {
			// try to create pollyfile
			err := t.createAudioPolly(text, pollyFilename)
			if err != nil {
				// if fail, play flitefile
				return fliteFilename, nil
			}
			// play new pollyfile
			return pollyFilename, nil
		}
		// flite is enabled and file exist, play it
		return fliteFilename, nil
	}

	// no cache, creation of the file
	if t.flite {
		err = t.createAudioFlite(text, fliteFilename)
		if err != nil {
			return "", errors.Wrapf(err, "Failed to create any sound")
		}
		return fliteFilename, nil
	} else {
		err = t.createAudioPolly(text, pollyFilename)
		if err != nil {
			err = t.createAudioFlite(text, fliteFilename)
			if err != nil {
				return "", errors.Wrapf(err, "Failed to create any sound")
			}
			return fliteFilename, nil
		}
		return pollyFilename, nil
	}
}

// Say create a tempfile based on the choosen technology of TTS and order the player to
// play it on speaker
func (t *tts) Say(text string, p player.Player) error {
	filepath, err := t.getSayFilename(text)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve sound file")
	}
	p.PlayFilepath(filepath)
	return nil
}

func (t *tts) GetSay(text string) ([]byte, error) {
	filepath, err := t.getSayFilename(text)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to retrieve sound file")
	}
	return ioutil.ReadFile(filepath)
}
