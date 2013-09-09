package main

import (
	"flag"
	. "github.com/vincent-petithory/mpdclient"
	"log"
	"sync"
)

var noRatings = flag.Bool("no-ratings", false, "Disable ratings service")
var noPlaycounts = flag.Bool("no-playcounts", false, "Disable playcounts service")

func startMpdService(mpdc *MPDClient, service func(*MPDClient, []chan SongSticker, chan bool), songStickerChangeHandlers []songStickerChangeHandler, wg *sync.WaitGroup, gate *Gate) {
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
			for i, _ := range channels {
				close(channels[i])
			}
			// Notify all services to shutdown
			gate.Open()
			wg.Done()
		}()
		service(mpdc, channels, gate.Waiter())
	}()
}

func startMpdServices() {
	var wg sync.WaitGroup

	mpdc, err := Connect("localhost", 6600)
	if err != nil {
		panic(err)
	}
	defer mpdc.Close()

	gate := NewGate()

	if !*noPlaycounts {
		startMpdService(mpdc, RecordPlayCounts, []songStickerChangeHandler{generateMostPlayedSongs(mpdc, "Most Played", 50)}, &wg, &gate)
		log.Println("Started Playcounts service... ")
	}
	if !*noRatings {
		startMpdService(mpdc, ListenRatings, []songStickerChangeHandler{generateBestRatedSongs(mpdc, "Best Rated", 50)}, &wg, &gate)
		log.Println("Started Ratings service... ")
	}
	wg.Wait()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	for {
		startMpdServices()
		log.Println("Restarting...")
	}
}
