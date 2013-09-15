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
