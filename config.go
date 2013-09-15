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
	"encoding/json"
	"io"
	"io/ioutil"
)

type Config struct {
	PlaycountsEnabled       bool
	MostPlayedPlaylistName  string
	MostPlayedPlaylistLimit uint
	RatingsEnabled          bool
	BestRatedPlaylistName   string
	BestRatedPlaylistLimit  uint
	MPDHost                 string
	MPDPort                 uint
	MPDPassword             string
}

func DefaultConfig() *Config {
	c := Config{true, "Most Played", 50, true, "Best Rated", 50, "localhost", 6600, ""}
	return &c
}

func (c *Config) ReadFrom(r io.Reader) (n int64, err error) {
	data, err := ioutil.ReadAll(r)
	n = int64(len(data))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, c)
	return
}

func (c *Config) Read(p []byte) (n int, err error) {
	data, err := json.Marshal(c)
	p = data
	return len(p), err
}

func (c *Config) Write(p []byte) (n int, err error) {
	err = json.Unmarshal(p, c)
	return len(p), err
}

func (c *Config) WriteTo(w io.Writer) (n int64, err error) {
	data, err := json.Marshal(c)
	if err != nil {
		return int64(len(data)), err
	}
	// On writing to a Writer, it makes sense to
	// use a pretty-printed output
	var buf bytes.Buffer
	err = json.Indent(&buf, data, "", "  ")
	if err != nil {
		return 0, err
	}
	buf.WriteByte('\n')
	n1, err := w.Write(buf.Bytes())
	return int64(n1), err
}
