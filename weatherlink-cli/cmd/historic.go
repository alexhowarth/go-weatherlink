package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type historicTime struct {
	t time.Time
}

var start historicTime
var end historicTime

func (h *historicTime) Set(s string) error {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*&h.t = t
	return nil
}

func (h *historicTime) String() string {
	return ""
}

func (h *historicTime) Type() string {
	return "string"
}

var historicCmd = &cobra.Command{
	Use:   "historic",
	Short: "Historic weather",
	Long:  `Provide a start and end time in RFC3339 format with a span no greater than 24 hours`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.HistoricGeneric(station, start.t, end.t)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		printJSON(resp)
	},
}

func init() {
	historicCmd.Flags().IntVar(&station, "station", 0, "numeric station id")
	historicCmd.Flags().Var(&start, "start", "start date (RFC3339)")
	historicCmd.Flags().Var(&end, "end", "end date (RFC3339)")
	historicCmd.MarkFlagRequired("station")
	historicCmd.MarkFlagRequired("start")
	historicCmd.MarkFlagRequired("end")
	rootCmd.AddCommand(historicCmd)
}
