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
	"encoding/csv"
	"fmt"
	. "github.com/vincent-petithory/mpdclient"
	"io"
)

// CsvFeed implements importing song stickers from a csv file.
// The csv file must use ',' as a separator, and contain 3 fields for each record.
// With the records being uri, sticker name, sticker value, in that order.
// The uri is, as usual, relative to the MPD music directory.
type CsvFeed struct {
	r *csv.Reader
}

func (feed *CsvFeed) Feed(ssCh chan SongSticker) error {
	const NFields = 3
	defer close(ssCh)
	for {
		record, err := feed.r.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n := len(record)
		if n != NFields {
			return fmt.Errorf("Expected %d fields, got %d for field: \"%v\"", NFields, n, record)
		}
		songSticker := SongSticker{Uri: record[0], Name: record[1], Value: record[2]}
		ssCh <- songSticker
	}
}

func (feed *CsvFeed) Close() error {
	return nil
}

func NewCsvFeed(r io.Reader) (SongStickerFeeder, error) {
	// default csv decoder, using ',' as separator
	return &CsvFeed{r: csv.NewReader(r)}, nil
}
