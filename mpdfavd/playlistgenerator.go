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
		if uint(i+1) >= max {
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
