package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup the list of sounds into an archive",
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

		sounds, err := list()
		if err != nil {
			return
		}

		baseDir := time.Now().Format("./bell-backup-20060102-150405")
		err = os.Mkdir(baseDir, 0755)
		if err != nil {
			logrus.WithField("error", err).Error("Failed to create the directory to store the backup")
			return
		}

		// save sound list
		soundFile, err := os.Create(baseDir + "/sounds.json")
		if err != nil {
			logrus.WithField("error", err).Error("Failed to create temporary file to store sounds list")
			return
		}
		enc := json.NewEncoder(soundFile)
		enc.SetIndent("", "  ")
		err = enc.Encode(sounds)
		if err != nil {
			logrus.WithField("error", err).Error("Failed to encode list of sounds as json")
			return
		}

		// save sounds
		soundDir := baseDir + "/sounds"
		err = os.Mkdir(soundDir, 0755)
		if err != nil {
			logrus.WithField("error", err).Error("Failed to create directory to store the downloaded sounds")
		}

		for i, sound := range sounds {
			logrus.Infof("[%v%%] Retrieving sound %v", (i+1)*100/len(sounds), sound.Name)
			get(sound.Name, fmt.Sprintf("%v/%v.mp3", soundDir, sound.Name))
			time.Sleep(sleepRate)
		}

		archiveName := baseDir + ".tar.xz"
		// create archive
		tar := exec.Command(
			"tar",
			"-c",
			"--xz",
			"--remove-files",
			"-f",
			archiveName,
			baseDir,
		)
		out, err := tar.CombinedOutput()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tarArgs":   strings.Join(tar.Args, " "),
				"tarOutput": out,
				"error":     err,
			}).Error("Failed to create the archive")
			return
		}
		logrus.Infof("Successfully created archive %v", archiveName)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringP("rate", "r", "0", "Number of request per minute")

}
