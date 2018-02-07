package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	// ListPath is the path to list sounds
	ListPath = "/api/v1/"
	// PlayPath is the path to play sounds
	PlayPath = "/api/v1/play/"
	// TtsPath is the path used to push content to read
	TtsPath = "/api/v1/tts"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bellctl",
	Short: "This allow to control a \"bell\" server with simple commands",
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
}
