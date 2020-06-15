package weatherlink_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexhowarth/go-weatherlink"
)

func TestSignatureParams(t *testing.T) {

	conf := &weatherlink.Config{
		Key:    "mykey",
		Secret: "mysecret",
	}
	wl := conf.NewClient()

	p := wl.MakeSignatureParams()
	p["foo"] = "bar"
	p["t"] = "123"

	got := p.String()
	expect := "api-keymykeyfoobart123"
	if got != expect {
		t.Fatalf("Expected %v got %v", expect, got)
	}
}

func TestBuildURL(t *testing.T) {

	conf := &weatherlink.Config{
		Key:    "mykey",
		Secret: "mysecret",
	}
	wl := conf.NewClient()

	p := wl.MakeSignatureParams()
	p["foo"] = "bar"
	p["t"] = "123"

	got := wl.BuildURL("/foo", p)
	expect := "https://api.weatherlink.com/v2/foo?api-key=mykey&api-signature=e576785c250d8c8db2e5fc2b7857b4c39ee56958107b978137e10d0fa6c1bc7b&t=123"
	if got != expect {
		t.Fatalf("Expected %v got %v", expect, got)
	}
}

// Provide a Transport for testing purposes
type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

func TestCurrent(t *testing.T) {

	conf := &weatherlink.Config{
		Key:    "mykey",
		Secret: "mysecret",
		Client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(t, "current.json"))),
			}, nil
		})}}

	wl := conf.NewClient()

	c, err := wl.Current(2970)
	if err != nil {
		t.Fatal(err)
	}

	{
		expect := 2970
		got := c.StationID
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}
	{
		expect := 1591894200
		got := c.Sensors[0].Data[0].Ts
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}
}

func TestHistoric(t *testing.T) {

	conf := &weatherlink.Config{
		Key:    "mykey",
		Secret: "mysecret",
		Client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(t, "historic.json"))),
			}, nil
		})}}

	wl := conf.NewClient()

	start := time.Now().Add(-time.Hour * 2)
	end := time.Now().Add(-time.Hour * 1)

	c, err := wl.Historic(2970, start, end)
	if err != nil {
		t.Fatal(err)
	}

	{
		expect := 2970
		got := c.StationID
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}
}

func TestStations(t *testing.T) {

	conf := &weatherlink.Config{
		Key:    "mykey",
		Secret: "mysecret",
		Client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(t, "stations.json"))),
			}, nil
		})}}

	wl := conf.NewClient()

	s, err := wl.AllStations()
	if err != nil {
		t.Fatal(err)
	}
	{
		expect := 1
		got := len(s.Stations)
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}
	{
		expect := 2970
		got := s.Stations[0].StationID
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}
	{
		expect := "Foo station"
		got := s.Stations[0].StationName
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}

	s, err = wl.Stations([]int{2970})
	if err != nil {
		t.Fatal(err)
	}

	{
		expect := 1
		got := len(s.Stations)
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}
	{
		expect := 2970
		got := s.Stations[0].StationID
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}
}

func TestSensors(t *testing.T) {

	conf := &weatherlink.Config{
		Key:    "mykey",
		Secret: "mysecret",
		Client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(t, "sensors.json"))),
			}, nil
		})}}

	wl := conf.NewClient()

	s, err := wl.AllSensors()
	if err != nil {
		t.Fatal(err)
	}

	{
		expect := 2
		got := len(s.Sensors)
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}
	{
		expect := "Vantage Vue, Wireless"
		got := s.Sensors[0].ProductName
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}

	s, err = wl.Sensors([]int{2970})
	if err != nil {
		t.Fatal(err)
	}

	{
		expect := 2
		got := len(s.Sensors)
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}
	{
		expect := "Vantage Vue, Wireless"
		got := s.Sensors[0].ProductName
		if got != expect {
			t.Fatalf("Expected %v got %v", expect, got)
		}
	}

}

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
