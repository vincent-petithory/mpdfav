/* Copyright (C) 2013 Vincent Petithory <vincent.petithory@gmail.com>
 *
 * This file is part of mpdfav.
 *
 * mpdfav is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mpdfav is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with mpdfav.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

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
