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
	. "github.com/vincent-petithory/mpdfav"
	"log"
	"strconv"
	"time"
)

const (
	songPlayedThresholdSeconds = 10
	tickMillis                 = 900
	PlaycountSticker           = "playcount"
)

type songStatusInfo struct {
	StatusInfo Info
	SongInfo   Info
}

func incSongPlayCount(songInfo *Info, mpdc *MPDClient) (int, error) {
	newval, err := AdjustIntStickerBy(mpdc, PlaycountSticker, (*songInfo)["file"], 1)
	if err != nil {
		return -1, err
	}
	return newval, err
}

func considerSongPlayed(statusInfo *Info, limit int) bool {
	current, total := statusInfo.Progress()
	if total == 0 || current == 0 {
		return false
	}
	return (total - current) < limit
}

func checkSongChange(si *songStatusInfo, mpdc *MPDClient) (bool, error) {
	info, err := mpdc.Status()
	if err != nil {
		return false, err
	}
	defer func() {
		si.StatusInfo = *info
	}()

	if (*info)["songid"] != si.StatusInfo["songid"] {
		if played := considerSongPlayed(&si.StatusInfo, songPlayedThresholdSeconds); played {
			return true, nil
		}
	}
	return false, nil
}

func processStateUpdate(si *songStatusInfo, mpdc *MPDClient, channels []chan SongSticker) error {
	changed, err := checkSongChange(si, mpdc)
	if err != nil {
		return err
	}
	if changed {
		playcount, err := incSongPlayCount(&si.SongInfo, mpdc)
		if err != nil {
			return err
		}

		songSticker := SongSticker{si.SongInfo["file"], PlaycountSticker, strconv.Itoa(playcount)}
		for _, channel := range channels {
			c := channel
			go func() {
				c <- songSticker
			}()
		}
		log.Printf("playcounts: %s playcount=%d\n", si.SongInfo["Title"], playcount)
	}
	// We store the current song after processing,
	// since that should be the next song playing already.
	songInfo, err := mpdc.CurrentSong()
	if err != nil {
		return err
	}
	si.SongInfo = *songInfo
	if err != nil {
		return err
	}
	return nil
}

func RecordPlayCounts(mpdc *MPDClient, channels []chan SongSticker, quit chan bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic in RecordPlayCounts: %s\n", err)
		}
	}()
	statusInfo, err := mpdc.Status()
	if err != nil {
		log.Panic(err)
	}
	songInfo, err := mpdc.CurrentSong()
	if err != nil {
		log.Panic(err)
	}

	si := songStatusInfo{}
	si.StatusInfo = *statusInfo
	si.SongInfo = *songInfo

	idleSub := mpdc.Idle("player")
	pollCh := time.Tick(tickMillis * time.Millisecond)
	ignorePoll := si.StatusInfo["state"] != "play"

	for {
		select {
		case <-pollCh:
			if !ignorePoll {
				err = processStateUpdate(&si, mpdc, channels)
				if err != nil {
					log.Panic(err)
				}
			}
		case <-idleSub.Ch:
			err := processStateUpdate(&si, mpdc, channels)
			if err != nil {
				log.Panic(err)
			}

			// Suspend poll goroutine if player is stopped or paused
			ignorePoll = si.StatusInfo["state"] != "play"
		case <-quit:
			return
		}
	}
}
