package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	// ListPath is the path to list sounds
	ListPath = "/api/v1/sounds"
	// PlayPath is the path to play sounds
	PlayPath = "/api/v1/play/"
	// TtsPath is the path used to push content to read
	TtsPath = "/api/v1/tts"
	// DeleteSoundPath is the path used to delete sounds from the library
	DeleteSoundPath = "/api/v1/sounds/"

	// AddSoundPath is the path used to upload new sounds
	AddSoundPath = "/api/v1/sounds"
	// GetSoundPath is the path used to retrieve sound content
	GetSoundPath = "/api/v1/sounds/"
)

var (
	tagOption bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bellctl",
	Short: "This allow to control a \"bell\" server with simple commands",
	Long: `You can controal a bell server. To choose your bell server use the env variable BELL_ADDRESS.addCmd
Example: 
	export BELL_ADDRESS=http://localhost:10101
	bellctl list 
	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.BindEnv("bell.address", "BELL_ADDRESS")
	viper.SetDefault("bell.address", "http://localhost:10101")
	viper.BindEnv("verbose", "VERBOSE")
	viper.SetDefault("verbose", false)
	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
