package main

import (
	. "github.com/vincent-petithory/mpdclient"
	"log"
	"sort"
	"strconv"
)

type songStickerChangeHandler func(SongSticker)

func ListenSongStickerChange(songStickerChange chan SongSticker, handler songStickerChangeHandler) {
	for {
		songSticker, ok := <-songStickerChange
		if ok {
			handler(songSticker)
		} else {
			return
		}
	}
}

func generateBestRatedSongs(mpdc *MPDClient, playlistName string, max int) songStickerChangeHandler {
	f := func(songSticker SongSticker) {
		err := mpdc.PlaylistClear(playlistName)
		if err != nil {
			log.Fatal(err)
		}
		songStickers, err := mpdc.StickerFind(StickerSongType, "/", RatingSticker)
		if err != nil {
			log.Fatal(err)
		}
		sort.Sort(sort.Reverse(songStickers))
		for i, songSticker := range songStickers {
			rating, err := strconv.Atoi(songSticker.Value)
			if err != nil {
				continue
			}
			if rating < 1 {
				continue
			}
			mpdc.PlaylistAdd(playlistName, songSticker.Uri)
			if i >= max {
				break
			}
		}
	}
	return songStickerChangeHandler(f)
}

func generateMostPlayedSongs(mpdc *MPDClient, playlistName string, max int) songStickerChangeHandler {
	f := func(songSticker SongSticker) {
		err := mpdc.PlaylistClear(playlistName)
		if err != nil {
			log.Fatal(err)
		}
		songStickers, err := mpdc.StickerFind(StickerSongType, "/", PlaycountSticker)
		if err != nil {
			log.Fatal(err)
		}
		sort.Sort(sort.Reverse(songStickers))
		for i, songSticker := range songStickers {
			_, err = strconv.Atoi(songSticker.Value)
			if err != nil {
				continue
			}
			mpdc.PlaylistAdd(playlistName, songSticker.Uri)
			if i >= max {
				break
			}
		}
	}
	return songStickerChangeHandler(f)
}
