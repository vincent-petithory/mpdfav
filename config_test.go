package mpdfav

import (
	"bytes"
	"testing"
)

func TestConfig(t *testing.T) {
	c := DefaultConfig()
	expectedConfig := Config{false, "Most Played", 30, true, "Best Songs", 20, "192.168.0.255", 6601, "guessme"}

	// var buf bytes.Buffer
	var jsonBlob = []byte(`{
		"PlaycountsEnabled": false,
		"MostPlayedPlaylistName": "Most Played",
		"MostPlayedPlaylistLimit": 30,
		"RatingsEnabled": true,
		"BestRatedPlaylistName": "Best Songs",
		"BestRatedPlaylistLimit": 20,
		"MPDHost": "192.168.0.255",
		"MPDPort": 6601,
		"MPDPassword": "guessme"
		}`)

	_, err := c.Write(jsonBlob)
	if err != nil {
		t.Fatal(err)
	}
	if *c != expectedConfig {
		t.Fatalf("Expected %v, got %v\n", expectedConfig, *c)
	}
}

func TestIncompleteConfig(t *testing.T) {
	c := DefaultConfig()
	expectedConfig := DefaultConfig()
	expectedConfig.PlaycountsEnabled = false
	expectedConfig.MPDHost = "192.168.0.255"

	var jsonBlob = []byte(`{
		"PlaycountsEnabled": false,
		"MPDHost": "192.168.0.255"
		}`)

	_, err := c.Write(jsonBlob)
	if err != nil {
		t.Fatal(err)
	}
	if *c != *expectedConfig {
		t.Fatalf("Expected %v, got %v\n", expectedConfig, *c)
	}
}

func TestFullCycleReadWrite(t *testing.T) {
	var jsonBlob = []byte(`{
		"MostPlayedPlaylistName": "Listen to this",
		"MPDHost": "192.168.0.255"
		}`)

	c1 := DefaultConfig()
	c1.MPDPassword = "random"
	c1.Write(jsonBlob)

	var buf bytes.Buffer
	c1.WriteTo(&buf)

	c2 := DefaultConfig()
	c2.ReadFrom(&buf)
	if *c1 != *c2 {
		t.Fatalf("%v and %v do not match\n", c1, c2)
	}
}
