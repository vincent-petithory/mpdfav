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
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/vincent-petithory/mpdclient"
)

type MPDStickerDBFeed struct {
	db *sql.DB
}

func (sd *MPDStickerDBFeed) Feed(ssCh chan SongSticker) error {
	rows, err := sd.db.Query("SELECT uri, name, value FROM sticker WHERE type='song'")
	if err != nil {
		return err
	}
	defer rows.Close()
	defer close(ssCh)

	for rows.Next() {
		var ss SongSticker
		rows.Scan(&ss.Uri, &ss.Name, &ss.Value)
		ssCh <- ss
	}
	return nil
}

func (sd *MPDStickerDBFeed) Close() error {
	return sd.db.Close()
}

func NewMPDStickerDBFeed(filepath string) (SongStickerFeeder, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &MPDStickerDBFeed{db}, nil
}
