package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{`upload`},
	Short:   "Add new sounds to library",
	Long:    `Allows to add new sounds to library. Client endpoint not implemented yet`,
	Example: `
  bellctl add --name toto --file ./file.mp3
	`,
	Run: func(cmd *cobra.Command, args []string) {
		soundName := cmd.Flag("name").Value.String()
		soundFile := cmd.Flag("file").Value.String()

		err := add(soundName, soundFile, tagsList...)
		if err != nil {
			return
		}
		logrus.Info("Sound successfully uploaded")

	},
}

var tagsList []string

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("name", "n", "", "sound name to use as a sound identifier")
	addCmd.MarkFlagRequired("name")

	addCmd.Flags().StringP("file", "f", "", "Filepath of the file to send to remote end of the peer")
	addCmd.MarkFlagRequired("file")

	addCmd.Flags().StringSliceVarP(&tagsList, "tags", "t", []string{}, "List of tags to add to the sound")
}

func add(file, name string, tags ...string) (err error) {
	address, err := url.Parse(viper.GetString("bell.address") + AddSoundPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":          err,
			"server address": viper.GetString("bell.address"),
			"method":         "addCmd.Run",
		}).Error("Failed to build url")
		return
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("uploadFile", file) // this name shoud be seen later

	fh, err := os.Open(file)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"filename": file,
		}).Error("Failed to read the given file")
		return
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to copy content of file to request writer")
		return
	}

	bodyWriter.WriteField("name", name)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// add tags to the query if tag list is more than nothing
	if len(tags) > 0 {
		values := address.Query()
		for _, tag := range tags {
			values.Add("tag", tag)
		}
		address.RawQuery = values.Encode()
	}

	resp, err := http.Post(address.String(), contentType, bodyBuf)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to send data to bell server")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		content, _ := ioutil.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"body":        string(content),
			"soundName":   name,
		}).Info("Failed to add new sound to bell server")
		return
	}
	return
}
