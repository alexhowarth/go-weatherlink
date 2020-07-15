package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var stations []int

var stationsCmd = &cobra.Command{
	Use:   "stations",
	Short: "List stations",
	Run: func(cmd *cobra.Command, args []string) {
		var resp interface{}
		var err error
		if len(stations) > 0 {
			resp, err = client.StationsGeneric(stations)
		} else {
			resp, err = client.AllStationsGeneric()
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		printJSON(resp)
	},
}

func init() {
	stationsCmd.Flags().IntSliceVar(&stations, "id", []int{}, "station ids")
	rootCmd.AddCommand(stationsCmd)
}
