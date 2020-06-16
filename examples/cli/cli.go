package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/alexhowarth/go-weatherlink"
)

// Simple command line script that will list stations and get current weather data in json

// Could be use with jq, for example:
// go run examples/cli/cli.go --key mykey --secret mysecret --list-stations | jq -r '.stations[].station_id'
// go run examples/cli/cli.go --key mykey --secret mysecret --station 123 | jq -r '.sensors[].data[].wind_dir'

// To turn this into a binary:
// go build examples/cli/cli.go
// ./cli --help
func main() {

	var key = flag.String("key", "", "api key")
	var secret = flag.String("secret", "", "api secret")
	var station = flag.Int("station", 0, "station id")
	var listStations = flag.Bool("list-stations", false, "list all stations")

	flag.Parse()

	config := &weatherlink.Config{
		Key:    *key,
		Secret: *secret,
	}
	client := config.NewClient()

	if *listStations {
		st, err := client.AllStations()
		if err != nil {
			os.Exit(1)
		}
		printJSON(st)
	}

	if *station > 0 {
		cu, err := client.Current(*station)
		if err != nil {
			os.Exit(1)
		}
		printJSON(cu)
	}

	os.Exit(0)
}

func printJSON(data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(string(b))
}
