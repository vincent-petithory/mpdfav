package main

import (
	//"flag"
	"github.com/vincent-petithory/mpdfav"
	"os"
	"sync"
)

func main() {
	mpdc, err := mpdfav.Connect("localhost", 6600)
	defer mpdc.Close()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		mpdfav.RecordPlayCounts(mpdc)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		mpdfav.ListenRatings(mpdc)
	}()

	wg.Wait()
}
