package cmd

import (
	"github.com/spf13/cobra"
)

var sensorPath string

var sensorcatalogCmd = &cobra.Command{
	Use:   "sensorcatalog",
	Short: "Downloads the sensor catalogue",
	Run: func(cmd *cobra.Command, args []string) {
		client.SensorCatalog(sensorPath)
	},
}

func init() {
	sensorcatalogCmd.Flags().StringVar(&sensorPath, "path", "./sensor-catalog.json", "path to save catalog to")
	rootCmd.AddCommand(sensorcatalogCmd)
}
