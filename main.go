package main

import (
	"flag"
	"log"
	"github.com/vincent-petithory/mpdfav"
	"sync"
)

var noRatings = flag.Bool("no-ratings", false, "Disable ratings service")
var noPlaycounts = flag.Bool("no-playcounts", false, "Disable playcounts service")

func startService(host string, port uint, service func(*mpdfav.MPDClient), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		mpdc, err := mpdfav.Connect(host, port)
		defer mpdc.Close()
		if err != nil {
			panic(err)
		}
		service(mpdc)
	}()
}

func main() {
	var wg sync.WaitGroup

	if !*noPlaycounts {
		startService("localhost", 6600, mpdfav.RecordPlayCounts, &wg)
		log.Print("Started Playcounts service... ")
	}
	if !*noRatings {
		startService("localhost", 6600, mpdfav.ListenRatings, &wg)
		log.Print("Started Ratings service... ")
	}

	wg.Wait()
}
