package main

import (
	"log"
	. "github.com/vincent-petithory/mpdclient"
	"sort"
)

type songMetadata struct {
	Uri      string
	Metadata string
	Value    string
}

type songMetadataChangeHandler func(songMetadata)

func ListenSongMetadataChange(songMetadataChange chan songMetadata, handler songMetadataChangeHandler) {
	for {
		songMetadata, ok := <-songMetadataChange
		if ok {
			handler(songMetadata)
		} else {
			return
		}
	}
}

func generateBestRatedSongs(mpdc *MPDClient, playlistName string, max int) songMetadataChangeHandler {
	f := func(songMetadata songMetadata) {
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
			mpdc.PlaylistAdd(playlistName, songSticker.Uri)
			if i >= max {
				break
			}
		}
	}
	return songMetadataChangeHandler(f)
}

func generateMostPlayedSongs(mpdc *MPDClient, playlistName string, max int) songMetadataChangeHandler {
	f := func(songMetadata songMetadata) {
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
			mpdc.PlaylistAdd(playlistName, songSticker.Uri)
			if i >= max {
				break
			}
		}
	}
	return songMetadataChangeHandler(f)
}
