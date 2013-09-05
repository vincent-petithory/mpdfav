package main

import (
	"errors"
	"fmt"
	. "github.com/vincent-petithory/mpdclient"
	. "github.com/vincent-petithory/mpdfav"
	"log"
	"strconv"
)

const (
	RatingSticker  = "rating"
	RatingsChannel = "ratings"
)

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
		return -1, errors.New(fmt.Sprintf("Invalid rating code: %s", rateMsg))
	}

	newval, err := AdjustIntStickerBy(mpdc, RatingSticker, (*songInfo)["file"], val)
	if err != nil {
		return -1, err
	}
	return newval, err
}

func ListenRatings(mpdc *MPDClient, channels []chan SongSticker) {
	err := mpdc.Subscribe(RatingsChannel)
	if err != nil {
		panic(err)
	}

	statusInfo, err := mpdc.Status()
	if err != nil {
		panic(err)
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
					log.Println(err)
				} else {
					for _, msg := range msgs {
						msgsCh <- msg
					}
				}
			case "player":
				statusInfo, err := mpdc.Status()
				if err != nil {
					log.Println(err)
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
					} else {
						log.Println(err)
					}
				} else {
					log.Println(err)
				}
			} else {
				log.Printf("Client %s already rated\n", thisClientId)
			}
		case statusInfo := <-playerCh:
			if currentSongId != statusInfo["songid"] {
				clientsSentRating = make([]string, 0)
			}
		}
	}
}
