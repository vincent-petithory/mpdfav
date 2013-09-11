package main

import (
	. "github.com/vincent-petithory/mpdclient"
	"log"
	"sort"
	"strconv"
)

type songStickerChangeHandler func(SongSticker)

func ListenSongStickerChange(songStickerChange chan SongSticker, handler songStickerChangeHandler) {
	for songSticker := range songStickerChange {
		handler(songSticker)
	}
}

func generatePlaylist(mpdc *MPDClient, stickerName string, playlistName string, max uint, descending bool) {
	log.Printf("playlist generator: regenerating %s\n", playlistName)
	err := mpdc.PlaylistClear(playlistName)
	if err != nil {
		log.Panic(err)
	}
	songStickers, err := mpdc.StickerFind(StickerSongType, "/", stickerName)
	if err != nil {
		log.Panic(err)
	}
	sort.Sort(sort.Reverse(songStickers))
	for i, songSticker := range songStickers {
		_, err = strconv.Atoi(songSticker.Value)
		if err != nil {
			continue
		}
		mpdc.PlaylistAdd(playlistName, songSticker.Uri)
		if uint(i) >= max {
			break
		}
	}
}

func generateBestRatedSongs(mpdc *MPDClient, playlistName string, max uint) songStickerChangeHandler {
	f := func(songSticker SongSticker) {
		generatePlaylist(mpdc, songSticker.Name, playlistName, max, true)
	}
	return songStickerChangeHandler(f)
}

func generateMostPlayedSongs(mpdc *MPDClient, playlistName string, max uint) songStickerChangeHandler {
	f := func(songSticker SongSticker) {
		generatePlaylist(mpdc, songSticker.Name, playlistName, max, true)
	}
	return songStickerChangeHandler(f)
}
