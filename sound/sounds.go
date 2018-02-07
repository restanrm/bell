package sound

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

type Sound struct {
	Name     string `json:"name"`
	FileName string `json:"file_name"`
}

type Sounds []Sound

func (s *Sounds) Load(file string) {
	fd, err := os.Open(file)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"filepath": file,
		}).Error("Failed to open store for sounds")
		os.Exit(-1)
	}

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"filepath": file,
		}).Error("Failed to read all content of file")
		os.Exit(-1)
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to unmarshal json data to Sounds type")
		return
	}
}
