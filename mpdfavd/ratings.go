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
	"errors"
	. "github.com/vincent-petithory/mpdclient"
	. "github.com/vincent-petithory/mpdfav"
	"log"
	"strconv"
)

const (
	RatingSticker  = "rating"
	RatingsChannel = "ratings"
)

var ErrInvalidRatingCode = errors.New("ratings: invalid rating code")

func rateSong(songInfo *Info, rateMsg string, mpdc *MPDClient) (int, error) {
	// fail fast if the rateMsg is invalid
	var val int
	switch rateMsg {
	case "+":
		fallthrough
	case "like":
		val = 1
	case "-":
		fallthrough
	case "dislike":
		val = -1
	default:
		val = 0
	}
	if val == 0 {
		return -1, ErrInvalidRatingCode
	}

	newval, err := AdjustIntStickerBy(mpdc, RatingSticker, (*songInfo)["file"], val)
	if err != nil {
		return -1, err
	}
	return newval, err
}

func ListenRatings(mpdc *MPDClient, channels []chan SongSticker, quit chan bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic in ListenRatings: %s\n", err)
		}
	}()
	err := mpdc.Subscribe(RatingsChannel)
	if err != nil {
		log.Panic(err)
	}

	statusInfo, err := mpdc.Status()
	if err != nil {
		log.Panic(err)
	}
	currentSongId := (*statusInfo)["songid"]

	clientsSentRating := make([]string, 0)

	msgsCh := make(chan ChannelMessage)
	playerCh := make(chan Info)

	go func() {
		idleSub := mpdc.Idle("message", "player")
		for {
			subsystem := <-idleSub.Ch
			switch subsystem {
			case "message":
				msgs, err := mpdc.ReadMessages()
				if err != nil {
					log.Panic(err)
				} else {
					for _, msg := range msgs {
						msgsCh <- msg
					}
				}
			case "player":
				statusInfo, err := mpdc.Status()
				if err != nil {
					log.Panic(err)
				} else {
					playerCh <- *statusInfo
				}
			}
		}
	}()

	for {
		select {
		case channelMessage := <-msgsCh:
			log.Println("Ratings: incoming message", channelMessage)
			if channelMessage.Channel != RatingsChannel {
				break
			}

			// FIXME find a way to Uidentify a client submitting a rating
			thisClientId := "0"
			clientExists := false
			for _, clientId := range clientsSentRating {
				if thisClientId == clientId {
					clientExists = true
					break
				}
			}
			if !clientExists {
				songInfo, err := mpdc.CurrentSong()
				if err == nil {
					if rating, err := rateSong(songInfo, channelMessage.Message, mpdc); err == nil {
						clientsSentRating = append(clientsSentRating, thisClientId)
						log.Printf("Ratings: %s rating=%d\n", (*songInfo)["Title"], rating)
						songSticker := SongSticker{(*songInfo)["file"], RatingSticker, strconv.Itoa(rating)}
						for _, channel := range channels {
							c := channel
							go func() {
								c <- songSticker
							}()
						}
					} else if err == ErrInvalidRatingCode {
						log.Println(err)
					} else {
						log.Panic(err)
					}
				} else {
					log.Panic(err)
				}
			} else {
				log.Printf("Client %s already rated\n", thisClientId)
			}
		case statusInfo := <-playerCh:
			if currentSongId != statusInfo["songid"] {
				clientsSentRating = make([]string, 0)
			}
		case <-quit:
			return
		}
	}
}
