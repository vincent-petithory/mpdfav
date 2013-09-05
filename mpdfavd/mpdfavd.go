package main

import (
	"flag"
	. "github.com/vincent-petithory/mpdclient"
	"log"
	"sync"
)

var noRatings = flag.Bool("no-ratings", false, "Disable ratings service")
var noPlaycounts = flag.Bool("no-playcounts", false, "Disable playcounts service")

func startMpdService(mpdc *MPDClient, service func(*MPDClient, []chan SongSticker), songStickerChangeHandlers []songStickerChangeHandler, wg *sync.WaitGroup) {
	wg.Add(len(songStickerChangeHandlers))
	channels := make([]chan SongSticker, len(songStickerChangeHandlers))
	for i, songStickerChangeHandler := range songStickerChangeHandlers {
		ch := make(chan SongSticker)
		channels[i] = ch
		handler := songStickerChangeHandler
		go func() {
			defer wg.Done()
			ListenSongStickerChange(ch, handler)
		}()
	}
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			for i, _ := range channels {
				close(channels[i])
			}
		}()
		service(mpdc, channels)
	}()
}

func main() {
	var wg sync.WaitGroup

	mpdc, err := Connect("localhost", 6600)
	if err != nil {
		panic(err)
	}
	defer mpdc.Close()

	if !*noPlaycounts {
		startMpdService(mpdc, RecordPlayCounts, []songStickerChangeHandler{generateMostPlayedSongs(mpdc, "Most Played", 50)}, &wg)
		log.Println("Started Playcounts service... ")
	}
	if !*noRatings {
		startMpdService(mpdc, ListenRatings, []songStickerChangeHandler{generateBestRatedSongs(mpdc, "Best Rated", 50)}, &wg)
		log.Println("Started Ratings service... ")
	}

	wg.Wait()
}
