package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/vincent-petithory/mpdclient"
	. "github.com/vincent-petithory/mpdfav"
	"log"
	"strconv"
	"sync"
)

func ImportSongSticker(mpdc *MPDClient, ss SongSticker) error {
	val, err := strconv.Atoi(ss.Value)
	if err != nil {
		return err
	}
	_, err = AdjustIntStickerBy(mpdc, ss.Name, ss.Uri, val)
	return err
}

type SongStickerFeeder interface {
	Feed(ssCh chan SongSticker)
}

type MPDStickerDB struct {
	db *sql.DB
}

func (sd *MPDStickerDB) Feed(ssCh chan SongSticker) error {
	rows, err := sd.db.Query("SELECT uri, name, value FROM sticker WHERE type='song'")
	if err != nil {
		log.Fatal(err)
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

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	db, err := sql.Open("sqlite3", "sticker.sql")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	feeder := MPDStickerDB{db}

	mpdc, err := Connect("localhost", 6600)
	if err != nil {
		log.Fatal(err)
	}
	defer mpdc.Close()

	ssCh := make(chan SongSticker)

	go feeder.Feed(ssCh)

	var wg sync.WaitGroup
	for ss := range ssCh {
		wg.Add(1)
		go func(ss SongSticker) {
			defer wg.Done()
			err := ImportSongSticker(mpdc, ss)
			switch err {
			case nil:
				fmt.Printf("[ OK ] %s: imported sticker « %s »\n", ss.Uri, ss.Name)
			case err.(*MPDError):
				fmt.Printf("[SKIP] %s: %s\n", ss.Uri, err)
			default:
				fmt.Printf("[ERR ] %v", err)
			}
		}(ss)
	}
	wg.Wait()
}
