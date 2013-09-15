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
