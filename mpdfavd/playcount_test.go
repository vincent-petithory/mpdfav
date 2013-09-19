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

package main

import (
	"testing"
)

func TestConsiderSongPlayed(t *testing.T) {
	type CueTest struct {
		SongCues
		Precision uint
		Satisfy   bool
	}
	cues := make(SongCues, 0)
	for i := 0.0; i < 1.0; i += 0.02 {
		cues = append(cues, float32(i))
	}
	tests := []CueTest{
		CueTest{
			SongCues{
				0.01, 0.1, 0.3, 0.4, 0.5, 0.7, 0.9,
			},
			4, true,
		},
		CueTest{
			SongCues{
				0.01, 0.1, 0.8, 0.9,
			},
			4, false,
		},
		CueTest{
			SongCues{
				0.1, 0.2, 0.3,
			},
			4, false,
		},
		CueTest{
			SongCues{
				0.01, 0.4, 0.8,
			},
			4, false,
		},
		CueTest{
			cues,
			4, true,
		},
	}
	for _, test := range tests {
		satisfied := considerSongPlayed(test.SongCues, test.Precision)
		if test.Satisfy && !satisfied {
			t.Fatal("Expected cues to be satisfied, they weren't")
		}
		if !test.Satisfy && satisfied {
			t.Fatal("Expected cues to not be satisfied, they were")
		}
	}
}
