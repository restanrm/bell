package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/restanrm/bell/sound"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore command help to put an archive sounds list back into a bell server",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		r := cmd.Flag("rate").Value.String()

		// define rate for requests
		var sleepRate time.Duration
		rate, err := strconv.Atoi(r)
		if err != nil {
			logrus.WithField("error", err).Error("Failed to convert rate option to integer. Please provide an integer.")
			return
		}
		if rate != 0 {
			sleepRate = time.Duration(60/rate) * time.Second
		} else {
			sleepRate = time.Duration(0)
		}

		// extract filename from arguments
		archivePath := args[0]

		fi, err := os.Stat(archivePath)
		if err != nil {
			logrus.WithField("error", err).Errorf("Failed to call \"stat\" on the given argument")
			return
		}
		if fi.IsDir() {
			logrus.Error("We need an archive, not a directory")
			return
		}
		a := strings.Split(archivePath, ".tar")
		if len(a) < 2 {
			logrus.Error("File doesn't seem to be an archive. It's name must contain \".tar\"")
			return
		}
		archive := filepath.Base(a[0])

		// create temp dir to store the files uncompressed
		extractDir, err := ioutil.TempDir("/tmp", "bell-extract")
		if err != nil {
			logrus.WithField("error", err).Error("Failed to create a temporary directory to extract the archive")
			return
		}
		defer os.RemoveAll(extractDir)

		// extract archive to that dir
		tar := exec.Command(
			"tar",
			"-x",
			"-f",
			archivePath,
			"-C",
			extractDir,
		)
		out, err := tar.CombinedOutput()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tarArgs":   strings.Join(tar.Args, " "),
				"tarOutput": out,
				"error":     err,
			}).Error("Failed to extract the archive")
			return
		}
		// load the produced json

		p := extractDir + "/" + archive
		reader, err := os.Open(p + "/sounds.json")
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":    err,
				"filepath": p + "/sounds.json",
			}).Error("Failed to open the list of sound file")
			return
		}
		var sounds []sound.Sound
		err = json.NewDecoder(reader).Decode(&sounds)
		if err != nil {
			logrus.WithField("error", err).Error("Failed to decode \"sounds.json\" as json content")
			return
		}
		// for each sound found, send it to the server and sleep for sleepRate

		for i, sound := range sounds {
			time.Sleep(sleepRate)
			err = add(p+"/sounds/"+sound.Name+".mp3", sound.Name, sound.Tags...)
			if err == nil {
				logrus.Infof("[%v%%] sound uploaded %v", (i+1)*100/len(sounds), sound.Name)
			}
		}
		logrus.Infof("Successfully restored the archive")

	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringP("rate", "r", "0", "Number of request per minute")
}
