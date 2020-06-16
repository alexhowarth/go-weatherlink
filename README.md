# Go WeatherLink v2 API client

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

See the examples directory for more.

## TODO

The following are not currently implemented:

 * Nodes
 * SensorActivity

This is work in progress. Let me know if something breaks.

## API

Documentation for the API can be found at [https://weatherlink.github.io/v2-api/](https://weatherlink.github.io/v2-api/)
