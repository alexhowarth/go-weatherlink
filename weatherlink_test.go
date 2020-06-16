package weatherlink

import (
	"testing"
)

func TestBuildURL(t *testing.T) {

	conf := &Config{
		Key:    "mykey",
		Secret: "mysecret",
	}
	wl := conf.NewClient()

	p := wl.MakeSignatureParams()
	p["foo"] = "bar"
	p["t"] = "123"

	got := wl.buildURL("/foo", p)
	expect := "https://api.weatherlink.com/v2/foo?api-key=mykey&api-signature=e576785c250d8c8db2e5fc2b7857b4c39ee56958107b978137e10d0fa6c1bc7b&t=123"
	if got != expect {
		t.Fatalf("Expected %v got %v", expect, got)
	}
}

func TestSignatureParams(t *testing.T) {

	conf := &Config{
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
