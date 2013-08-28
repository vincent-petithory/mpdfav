package main

import (
	"flag"
	"log"
	"github.com/vincent-petithory/mpdfav"
	"sync"
)

var noRatings = flag.Bool("no-ratings", false, "Disable ratings service")
var noPlaycounts = flag.Bool("no-playcounts", false, "Disable playcounts service")

func startMpdService(mpdc *MPDClient, service func(*mpdfav.MPDClient), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		service(mpdc)
	}()
}

func main() {
	var wg sync.WaitGroup

	mpdc, err := mpdfav.Connect("localhost", 6600)
	if err != nil {
		panic(err)
	}
	defer mpdc.Close()

	if !*noPlaycounts {
		startMpdService(mpdc, mpdfav.RecordPlayCounts, &wg)
		log.Print("Started Playcounts service... ")
	}
	if !*noRatings {
		startMpdService(mpdc, mpdfav.ListenRatings, &wg)
		log.Print("Started Ratings service... ")
	}

	wg.Wait()
}
