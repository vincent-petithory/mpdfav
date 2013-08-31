package main

import (
	"errors"
	"fmt"
	. "github.com/vincent-petithory/mpdclient"
	"log"
	"strconv"
)

const (
	ratingSticker  = "rating"
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

	value, err := mpdc.StickerGet(
		StickerSongType,
		(*songInfo)["file"],
		ratingSticker,
	)
	if err != nil {
		return -1, err
	}
	if len(value) == 0 {
		value = "0"
	}
	intval, err := strconv.Atoi(value)
	if err != nil {
		return -1, err
	}

	intval += val

	err = mpdc.StickerSet(
		StickerSongType,
		(*songInfo)["file"],
		ratingSticker,
		strconv.Itoa(intval),
	)
	return intval, err
}

func ListenRatings(mpdc *MPDClient) {
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
				log.Println(">>> message event")
				msgs, err := mpdc.ReadMessages()
				if err != nil {
					log.Println(err)
				} else {
					for _, msg := range msgs {
						msgsCh <- msg
					}
				}
			case "player":
				log.Println(">>> player event")
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
					if rating, err := rateSong(songInfo, channelMessage.Message, mpdc); err != nil {
						log.Println(err)
					} else {
						clientsSentRating = append(clientsSentRating, thisClientId)
						log.Println(fmt.Sprintf("Ratings: %s rating=%d", (*songInfo)["Title"], rating))
					}
				} else {
					log.Println(err)
				}
			} else {
				log.Println(fmt.Sprintf("Client %s already rated", thisClientId))
			}
		case statusInfo := <-playerCh:
			if currentSongId != statusInfo["songid"] {
				log.Println("Ratings: song changed to", statusInfo["songid"])
				clientsSentRating = make([]string, 0)
			}
		}
	}
}
