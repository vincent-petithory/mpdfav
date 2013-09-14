package main

import (
	"encoding/json"
	. "github.com/vincent-petithory/mpdclient"
	"io"
	"io/ioutil"
)

type JsonFeed struct {
	songStickers SongStickerList
}

func (feed *JsonFeed) Feed(ssCh chan SongSticker) error {
	defer close(ssCh)
	for _, ss := range feed.songStickers {
		ssCh <- ss
	}
	return nil
}

func (feed *JsonFeed) Close() error {
	return nil
}

func NewJsonFeed(r io.Reader) (SongStickerFeeder, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var ssl SongStickerList
	err = json.Unmarshal(data, &ssl)
	if err != nil {
		return nil, err
	}
	return &JsonFeed{ssl}, nil
}
