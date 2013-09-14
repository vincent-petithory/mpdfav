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
