package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var sensors []int

var sensorsCmd = &cobra.Command{
	Use:   "sensors",
	Short: "Displays sensor data",
	Run: func(cmd *cobra.Command, args []string) {
		var resp interface{}
		var err error
		if len(sensors) > 0 {
			resp, err = client.SensorsGeneric(sensors)
		} else {
			resp, err = client.AllSensorsGeneric()
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		printJSON(resp)
	},
}

func init() {
	sensorsCmd.Flags().IntSliceVar(&sensors, "id", []int{}, "sensor ids")
	rootCmd.AddCommand(sensorsCmd)
}
