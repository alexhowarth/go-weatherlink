package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alexhowarth/go-weatherlink"
	"github.com/spf13/cobra"
)

var key string
var secret string
var station int
var client *weatherlink.Client
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "weatherlink-cli",
	Short: "Command line tool for the Davis WeatherLink v2 API",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&key, "key", "", "api key")
	rootCmd.PersistentFlags().StringVar(&secret, "secret", "", "api secret")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "display verbose output")
	rootCmd.MarkPersistentFlagRequired("key")
	rootCmd.MarkPersistentFlagRequired("secret")
}

func initConfig() {
	config := &weatherlink.Config{
		Key:    key,
		Secret: secret,
	}
	client = config.NewClient()
}

func printJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(string(b))
}
