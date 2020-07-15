package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Current weather",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.CurrentGeneric(station)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		printJSON(resp)
	},
}

func init() {
	currentCmd.Flags().IntVar(&station, "station", 0, "numeric station id")
	currentCmd.MarkFlagRequired("station")
	rootCmd.AddCommand(currentCmd)
}
