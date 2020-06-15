package main

import (
	"fmt"
	"time"

	"github.com/alexhowarth/go-weatherlink"
)

func main() {

	//configure
	config := weatherlink.Config{
		Key:    "mykey",
		Secret: "mysecret",
	}

	// build a client from the configuration
	wl := config.NewClient()

	// get all stations
	stations, err := wl.AllStations()
	if err != nil {
		fmt.Println(err)
		return
	}

	// iterate over stations
	for _, station := range stations.Stations {

		fmt.Printf("Found station ID %v (%v)\n", station.StationID, station.StationName)

		// get current weather conditions for this station
		current, err := wl.Current(station.StationID)
		if err != nil {
			fmt.Println(err)
			return
		}

		// iterate over the sensors
		for _, sensor := range current.Sensors {
			// for each sensor, get some data
			for _, data := range sensor.Data {
				fmt.Printf("Wind Direction: %v\n", data.WindDir)
				fmt.Printf("Wind Speed: %v\n", data.WindSpeed)
				fmt.Printf("Last updated: %v\n", time.Unix(int64(data.Ts), 0))
			}
		}

		// create start and end dates
		start := time.Now().Add(-time.Hour * 1)
		end := time.Now().Add(-time.Minute * 30)

		// get historic for this station
		h, err := wl.Historic(station.StationID, start, end)
		if err != nil {
			fmt.Println(err)
			return
		}
		// iterate over the sensor data for this historic data
		fmt.Printf("Historic data...")
		for _, sensor := range h.Sensors {
			// for each sensor, get some data
			for _, data := range sensor.Data {
				fmt.Printf("Date: %v\n", time.Unix(int64(data.Ts), 0))
				fmt.Printf("Prevailing wind Direction: %v\n", data.WindDirOfPrevail)
				fmt.Printf("Wind Speed high: %v\n", data.WindSpeedHi)
			}
		}

	}

}
