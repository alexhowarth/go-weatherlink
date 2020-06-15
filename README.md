# WeatherLink v2 API in Go

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
```

More examples in example/example.go

## TODO

The following are not currently implemented:

 * Nodes
 * SensorActivity

This is work in progress. Let me know if something breaks.

## API

Documentation for the API can be found [here](https://weatherlink.github.io/v2-api/)
