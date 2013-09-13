package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/vincent-petithory/mpdclient"
	. "github.com/vincent-petithory/mpdfav"
	"log"
	"os"
	"strconv"
	"sync"
)

const (
	FORMAT_MPD_STICKER_DB = "stickerdb"
	FORMAT_JSON           = "json"
	FORMAT_CSV            = "csv"
)

var showHelp = flag.Bool("help", false, "Displays this help message.")
var format = flag.String("format", FORMAT_MPD_STICKER_DB, fmt.Sprintf("Format of the data FILE. Valid values are: \"%s\", \"%s\", \"%s\".", FORMAT_MPD_STICKER_DB, FORMAT_JSON, FORMAT_CSV))

var (
	defaultConfigFile string = os.ExpandEnv("$HOME/.mpdfav.json")
	configFile        string
	conf              *Config
)

func ImportSongSticker(mpdc *MPDClient, ss SongSticker) error {
	val, err := strconv.Atoi(ss.Value)
	if err != nil {
		return err
	}
	_, err = AdjustIntStickerBy(mpdc, ss.Name, ss.Uri, val)
	return err
}

type SongStickerFeeder interface {
	Feed(ssCh chan SongSticker) error
	Close() error
}

type MPDStickerDBFeed struct {
	db *sql.DB
}

func (sd *MPDStickerDBFeed) Feed(ssCh chan SongSticker) error {
	rows, err := sd.db.Query("SELECT uri, name, value FROM sticker WHERE type='song'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	defer close(ssCh)

	for rows.Next() {
		var ss SongSticker
		rows.Scan(&ss.Uri, &ss.Name, &ss.Value)
		ssCh <- ss
	}
	return nil
}

func (sd *MPDStickerDBFeed) Close() error {
	return sd.db.Close()
}

func NewMPDStickerDBFeed(filepath string) (SongStickerFeeder, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &MPDStickerDBFeed{db}, nil
}

func PrintHelp() {
	fmt.Fprintf(os.Stderr, `Usage: %s [OPTION] [FILE]
Imports MPD sticker data from a source FILE of songs' sticker-like data.

Options:
`, os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}

func init() {
	flag.StringVar(&configFile, "config-file", defaultConfigFile, fmt.Sprintf("Use this config file instead of %s", defaultConfigFile))
}

func main() {
	conf = DefaultConfig()
	flag.Parse()
	if *showHelp {
		PrintHelp()
		os.Exit(0)
	}

	f, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	} else {
		n, err := conf.ReadFrom(f)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			log.Fatalf("No data could be read from %s\n", configFile)
		}
		f.Close()
	}

	filepath := flag.Arg(0)
	if filepath == "" {
		PrintHelp()
		os.Exit(1)
	}

	// Create mpd client
	var mpdc *MPDClient
	if conf.MPDPassword != "" {
		mpdc, err = ConnectAuth(conf.MPDHost, conf.MPDPort, conf.MPDPassword)
	} else {
		mpdc, err = Connect(conf.MPDHost, conf.MPDPort)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer mpdc.Close()

	// Create a feeder
	var feeder SongStickerFeeder
	switch *format {
	case FORMAT_MPD_STICKER_DB:
		feeder, err = NewMPDStickerDBFeed(filepath)
		if err != nil {
			log.Fatal(err)
		}
	case FORMAT_JSON:
		log.Fatalf("Not implemented\n")
	case FORMAT_CSV:
		log.Fatalf("Not implemented\n")
	default:
		PrintHelp()
		log.Fatalf("Invalid format %s\n", *format)

	}

	ssCh := make(chan SongSticker)
	defer feeder.Close()
	go feeder.Feed(ssCh)

	var wg sync.WaitGroup
	for ss := range ssCh {
		wg.Add(1)
		go func(ss SongSticker) {
			defer wg.Done()
			err := ImportSongSticker(mpdc, ss)
			switch err {
			case nil:
				fmt.Printf("[ OK ] %s: imported sticker « %s »\n", ss.Uri, ss.Name)
			case err.(*MPDError):
				fmt.Printf("[SKIP] %s: %s\n", ss.Uri, err)
			default:
				fmt.Printf("[ERR ] %v", err)
			}
		}(ss)
	}
	wg.Wait()
}
