package mpdfav

import (
	. "github.com/vincent-petithory/mpdclient"
	"strconv"
)

func AdjustIntStickerBy(mpdc *MPDClient, stickerName string, uri string, by int) (int, error) {
	value, err := mpdc.StickerGet(
		StickerSongType,
		uri,
		stickerName,
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
	newval := intval + by
	err = mpdc.StickerSet(
		StickerSongType,
		uri,
		stickerName,
		strconv.Itoa(newval),
	)
	return newval, err
}
