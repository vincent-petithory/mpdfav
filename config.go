package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

type config struct {
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

func defaultConfig() *config {
	c := config{true, "Most Played", 50, true, "Best Rated", 50, "localhost", 6600, ""}
	return &c
}

func (c *config) ReadFrom(r io.Reader) (n int64, err error) {
	data, err := ioutil.ReadAll(r)
	n = int64(len(data))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, c)
	return
}

func (c *config) Read(p []byte) (n int, err error) {
	data, err := json.Marshal(c)
	p = data
	return len(p), err
}

func (c *config) Write(p []byte) (n int, err error) {
	err = json.Unmarshal(p, c)
	return len(p), err
}

func (c *config) WriteTo(w io.Writer) (n int64, err error) {
	data, err := json.Marshal(c)
	if err != nil {
		return int64(len(data)), err
	}
	// On writing to a Writer, it makes sense to
	// use a pretty-printed
	var buf bytes.Buffer
	err = json.Indent(&buf, data, "", "  ")
	if err != nil {
		return 0, err
	}
	buf.WriteByte('\n')
	n1, err := w.Write(buf.Bytes())
	return int64(n1), err
}
