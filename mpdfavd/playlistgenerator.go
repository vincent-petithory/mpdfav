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

func generatePlaylist(mpdc *MPDClient, stickerName string, playlistName string, max int, descending bool) {
		err := mpdc.PlaylistClear(playlistName)
		if err != nil {
			log.Fatal(err)
		}
		songStickers, err := mpdc.StickerFind(StickerSongType, "/", stickerName)
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

func generateBestRatedSongs(mpdc *MPDClient, playlistName string, max int) songStickerChangeHandler {
	f := func(songSticker SongSticker) {
		generatePlaylist(mpdc, songSticker.Name, playlistName, max, true)
	}
	return songStickerChangeHandler(f)
}

func generateMostPlayedSongs(mpdc *MPDClient, playlistName string, max int) songStickerChangeHandler {
	f := func(songSticker SongSticker) {
		generatePlaylist(mpdc, songSticker.Name, playlistName, max, true)
	}
	return songStickerChangeHandler(f)
}
