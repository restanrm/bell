// Package sound describe what a sound service is and propose an
// inMemory implematation of such service
package sound

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"github.com/restanrm/bell/player"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Sounder is the interface to render sound service
type Sounder interface {
	CreateSound(name, filepath string, tags ...string) error
	UpdateSound(sound Sound) error
	DeleteSound(name string) error
	PlaySound(name string, player player.Player) error
	PlaySoundByTag(tag string, player player.Player) error
	GetSound(name string) ([]byte, error)
	GetSounds() []Sound
}

// Sound is the struct to represent a sound
type Sound struct {
	Name     string   `json:"name"`
	filePath string   `json:"file_name"`
	Tags     []string `json:"tags,omitempty"`
}

type inMemorySounds struct {
	configFile string
	m          map[string]Sound
	sync.RWMutex
}

var (
	// ErrNoTagMatch is the error return when no tag can be found for the sound
	ErrNoTagMatch        = func(t string) error { return fmt.Errorf("No sound found matching the tag %v", t) }
	ErrSoundAlreadyExist = errors.New("Sound already exist")
	ErrSoundNotFound     = errors.New("Sound not found")
)

// New create a new instance of a sounder
func New(filepath string) *inMemorySounds {
	err := dirExist(filepath)
	if err != nil {
		logrus.Fatal(errors.Wrapf(err, "Failed to find or create directory to store configuration file"))
		os.Exit(-1)
	}
	ims := &inMemorySounds{
		configFile: filepath,
		m:          make(map[string]Sound),
	}
	ss, err := load(filepath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Warn("An error happened when loading sounds database")
	} else {
		ims.Lock()
		defer ims.Unlock()
		for _, sound := range ss {
			ims.m[sound.Name] = sound
		}
	}
	return ims
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
			return errors.Wrapf(err, "Failed to create directory to store configuration files")
		}
	}
	return nil
}

type ssto struct {
	Name     string   `json:"name"`
	FileName string   `json:"file_name"`
	Tags     []string `json:"tags"`
}

func load(fp string) ([]Sound, error) {
	fd, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	var ss []ssto
	err = json.Unmarshal(data, &ss)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json data to Sounds type. err: %v", err)
	}

	var output []Sound
	for _, s := range ss {
		output = append(output, Sound{
			Name:     s.Name,
			filePath: filepath.Join(viper.GetString("soundDir"), s.FileName),
			Tags:     s.Tags,
		})
	}
	return output, nil
}

// Load some sounds into collection
func Load(file string) Sounder {
	sounds := &inMemorySounds{
		m: make(map[string]Sound),
	}
	sounds.Lock()
	defer sounds.Unlock()

	return sounds
}

func (s inMemorySounds) save() error {
	var ss []ssto
	for k, v := range s.m {
		ss = append(ss, ssto{Name: k, FileName: filepath.Base(v.filePath), Tags: v.Tags})
	}

	f, err := os.OpenFile(s.configFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "Failed to open configuration file to write new content")
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(ss)
	if err != nil {
		return errors.Wrapf(err, "Failed to encode store as json struct")
	}
	return nil
}

// CreateSound a new sound in a collections. The file is already on the disk
func (s inMemorySounds) CreateSound(name, filepath string, tags ...string) error {
	s.Lock()
	defer s.Unlock()

	ss := Sound{
		Name:     name,
		filePath: filepath,
		Tags:     tags,
	}
	s.m[name] = ss

	err := s.save()
	if err != nil {
		return errors.Wrapf(err, "Failed to save the current state of the sound library")
	}
	return nil
}

// UpdateSound a sound in a collection
func (s inMemorySounds) UpdateSound(sound Sound) error {
	s.Lock()
	defer s.Unlock()
	s.m[sound.Name] = sound
	err := s.save()
	if err != nil {
		return errors.Wrapf(err, "Failed to save the current state of the sound library")
	}
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
		return s.PlaySoundByTag(name, player)
	}
	go func() {
		player.PlayFilepath(ss.filePath)
	}()
	return nil
}

func (s inMemorySounds) PlaySoundByTag(tag string, player player.Player) error {
	s.RLock()
	defer s.RUnlock()
	contains := func(list []string, el string) bool {
		for _, a := range list {
			if a == el {
				return true
			}
		}
		return false
	}
	// build list of playable
	var playable []Sound
	for _, v := range s.m {
		if contains(v.Tags, tag) {
			playable = append(playable, v)
		}
	}
	if len(playable) == 0 {
		return ErrNoTagMatch(tag)
	}
	return s.PlaySound(playable[rand.Int()%len(playable)].Name, player)
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

func (s inMemorySounds) GetSound(name string) (ret []byte, err error) {
	s.RLock()
	defer s.RUnlock()
	ss, ok := s.m[name]
	if !ok {
		return s.getSoundByTag(name)
	}
	return ioutil.ReadFile(ss.filePath)
}

func (s inMemorySounds) getSoundByTag(tag string) (ret []byte, err error) {
	s.RLock()
	defer s.RUnlock()
	contains := func(list []string, el string) bool {
		for _, a := range list {
			if a == el {
				return true
			}
		}
		return false
	}
	// build list of playable
	var playable []Sound
	for _, v := range s.m {
		if contains(v.Tags, tag) {
			playable = append(playable, v)
		}
	}
	if len(playable) == 0 {
		return nil, ErrSoundNotFound
	}
	return ioutil.ReadFile(playable[rand.Int()%len(playable)].filePath)
}
