# WeatherLink v2 API in Go

This package provides a Go library for the Davis WeatherLink v2 API.

## Installation

~~~~
go get github.com/alexhowarth/go-weatherlink
~~~~

## Usage

The following example (change the Key and Secret in the file) outputs data for all of your weather stations.

~~~~
go run example/example.go
~~~~

## TODO

The following are not currently implemented:

 * Nodes
 * SensorActivity

This is work in progress. Let me know if something breaks.

## API

Documentation for the API can be found [here](https://weatherlink.github.io/v2-api/)
