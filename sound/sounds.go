package sound

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sort"
	"sync"

	"github.com/restanrm/bell/player"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Sounder interface {
	CreateSound(name, filepath string) error
	UpdateSound(sound Sound) error
	DeleteSound(name string) error
	PlaySound(name string, player player.Player) error
	GetSounds() []Sound
}

type Sound struct {
	Name     string `json:"name"`
	filePath string `json:"file_name"`
}

type inMemorySounds struct {
	m map[string]Sound
	sync.RWMutex
}

var (
	ErrSoundAlreadyExist = errors.New("Sound already exist")
	ErrSoundNotFound     = errors.New("Sound not found")
)

// Load some sounds into collection
func Load(file string) Sounder {
	// func (ss *Sounds) Load(file string) {
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

	type ssto struct {
		Name     string `json:"name"`
		FileName string `json:"file_name"`
	}

	var ss []ssto
	err = json.Unmarshal(data, &ss)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to unmarshal json data to Sounds type")
		return nil
	}

	sounds := &inMemorySounds{
		m: make(map[string]Sound),
	}
	sounds.Lock()
	defer sounds.Unlock()
	for _, s := range ss {
		sounds.m[s.Name] = Sound{
			Name:     s.Name,
			filePath: viper.GetString("soundDir") + "/" + s.FileName,
		}
	}

	return sounds
}

// CreateSound a new sound in a collections. The file is already on the disk
func (s inMemorySounds) CreateSound(name, filepath string) error {
	s.Lock()
	defer s.Unlock()

	ss := Sound{
		Name:     name,
		filePath: filepath,
	}
	s.m[name] = ss
	return nil
}

// UpdateSound a sound in a collection
func (s inMemorySounds) UpdateSound(sound Sound) error {
	s.Lock()
	defer s.Unlock()

	s.m[sound.Name] = sound
	return nil
}

// DeleteSound remove sound from a collection
func (s inMemorySounds) DeleteSound(name string) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.m[name]; !ok {
		return ErrSoundNotFound
	}
	delete(s.m, name)

	return nil
}

// PlaySound is playing a sound from a sound collection
func (s inMemorySounds) PlaySound(name string, player player.Player) error {
	s.RLock()
	defer s.RUnlock()
	ss, ok := s.m[name]
	if !ok {
		return ErrSoundNotFound
	}
	player.PlayFilepath(ss.filePath)
	return nil
}

// GetSounds return a list of sounds for inMemoryImplementation of the service
func (s inMemorySounds) GetSounds() []Sound {
	s.RLock()
	defer s.RUnlock()
	var out []Sound
	for _, v := range s.m {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
