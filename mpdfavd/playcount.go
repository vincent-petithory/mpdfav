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
	songPlayedPrecision = 4 // Number of uniformly distributed cue points to statisfy
	tickMillis          = 900
	PlaycountSticker    = "playcount"
)

type SongCues []float32

type songStatusInfo struct {
	StatusInfo Info
	SongInfo   Info
	SongCues   SongCues
}

func incSongPlayCount(songInfo *Info, mpdc *MPDClient) (int, error) {
	newval, err := AdjustIntStickerBy(mpdc, PlaycountSticker, (*songInfo)["file"], 1)
	if err != nil {
		return -1, err
	}
	return newval, err
}

func considerSongPlayed(songCues SongCues, precision uint) bool {
	var cueStep float32 = 1.0 / (float32(precision) + 1.0)
	var curCue float32 = cueStep
	for _, cue := range songCues {
		for curCue < cue && cue < curCue+cueStep {
			curCue += cueStep
		}
		if curCue >= 1.0 {
			return true
		}
	}
	return false
}

func processStateUpdate(si *songStatusInfo, mpdc *MPDClient, channels []chan SongSticker) error {
	info, err := mpdc.Status()
	if err != nil {
		return err
	}

	if (*info)["songid"] != si.StatusInfo["songid"] {
		songPlayed := considerSongPlayed(si.SongCues, songPlayedPrecision)
		if songPlayed {
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
		si.SongCues = make(SongCues, 0)
	} else {
		current, total := si.StatusInfo.Progress()
		si.SongCues = append(si.SongCues, float32(current)/float32(total))
	}
	si.StatusInfo = *info
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

	si := songStatusInfo{*statusInfo, *songInfo, make([]float32, 0)}

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
