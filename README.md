# Go WeatherLink v2 Client

[![GoDoc](https://godoc.org/github.com/alexhowarth/go-weatherlink?status.svg)](https://pkg.go.dev/github.com/alexhowarth/go-weatherlink)
[![Go Report Card](https://goreportcard.com/badge/alexhowarth/go-weatherlink)](https://goreportcard.com/report/github.com/alexhowarth/go-weatherlink)

This package provides a Go library for the Davis WeatherLink v2 API.

## Installation

~~~~
go get github.com/alexhowarth/go-weatherlink
~~~~

## Usage

```go
// configure
config := weatherlink.Config{
        Key:    "mykey",
        Secret: "mysecret",
}

// build a client
wl := config.NewClient()

// all weather stations associated with your key
st, err := wl.AllStations()
if err != nil {
        // handle error
}

// current conditions for a station
cu, err := wl.Current(123)
if err != nil {
        // handle error
}

// historic for a station
start := time.Now().Add(-time.Hour * 2)
end := time.Now()

h, err := wl.Historic(123, start, end)
if err != nil {
        // handle error
}

for _, v := range h.Sensors {
        for _, d := range v.Data {
                fmt.Printf("Time: %v Temp: %v\n", time.Unix(d.Ts, 0), d.TempOut)
	}
}
```

## Command line tool

This package contains the command line tool `weatherlink-cli`. To install and use it:
```bash
$ go install weatherlink-cli
$ weatherlink-cli --help
```

To extract certain data from the output, you might use [jq](https://stedolan.github.io/jq/):

```bash
$ weatherlink-cli historic --key mykey --secret mysecret --station 2970 --start="2020-07-08T00:00:00Z" --end="2020-07-08T01:00:00Z"| jq -r '.sensors[].data[] | "timestamp: \(.ts) temp_out: \(.temp_out) bar: \(.bar)"'

timestamp: 1594167300 temp_out: 76.8 bar: 30.018
timestamp: 1594167600 temp_out: 76.8 bar: 30.014
```

## TODO

The following are not currently implemented:

 * Nodes
 * SensorActivity

This is work in progress. Let me know if something breaks or if your sensor type is not supported.

## API

Documentation for the API can be found at [https://weatherlink.github.io/v2-api/](https://weatherlink.github.io/v2-api/)
