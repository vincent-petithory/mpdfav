# mpdfav -- Song playcounts and ratings out-of-the-box for MPD.

mpdfav provides auto-generated playlists from your music habits. Most played and Best rated playlists are currently available.

While playcounts are of course automatic, ratings are up to you.

Ratings in mpdfav are *not* the usual 0-5 stars. Instead, you vote "like" or "dislike" at the moment you hear a song, which will raise or lower the rating of that song. The more you listen and rate, the better mpdfav gets to know you!
This approach has the advantage to be lightweight: you don't need to browse your music library and rate all those songs, and do that again later when your music tastes changes.
Also, when your tastes change, you won't need to lower the 5 stars songs you liked previously best. The new songs you'll listen to and rate over time will seamlessly get their deserved rank.

`mpdfav` makes use of MPD stickers to store its data.

# Installation

1. [Install Go](http://golang.org/doc/install) if you didn't yet.
2. Install mpdfav and its dependancies:

        go get github.com/vincent-petithory/mpdfav/mpdfavd
        # if you also want the import utility (see below)
        go get github.com/vincent-petithory/mpdfav/mpdfav-import

# Usage

Start mpdfavd:

    mpdfavd

If this is your first time running, it will create a default config in `$HOME/.mpdfav.json`:

    {
      "PlaycountsEnabled": true,
      "MostPlayedPlaylistName": "Most Played",
      "MostPlayedPlaylistLimit": 80,
      "RatingsEnabled": true,
      "BestRatedPlaylistName": "Best Rated",
      "BestRatedPlaylistLimit": 50,
      "MPDHost": "localhost",
      "MPDPort": 6600,
      "MPDPassword": ""
    }

Adapt to your needs and re-run `mpdfavd`. Done.

The last thing you need to know is how to rate a song. For that you need `mpc`, the command-line `mpd` client.

    # rate good the current song
    mpc sendmessage ratings like
    # or
    mpc sendmessage ratings +

    # rate bad a song
    mpc sendmessage ratings dislike
    # or
    mpc sendmessage ratings -

Setting global key-bindings to 2 of those commands is recommended.

## Import from your previous music player

In case you could dump playcounts / ratings from your previous music player, you can import them in mpdfav, with `mpdfav-import`.
`mpdfav-import` uses the same config file than `mpdfavd`.
For now it knows how to import from the formats below.

* Another MPD sticker DB.
* JSON; see below.
* csv.

In all file formats, the URIs of the songs are relative to your MPD music directory.

### JSON format

Here's an example of a file format:

    [
      {
        "Uri": "track-01.ogg",
        "Name": "playcount",
        "Value": "1"
      },
      {
        "Uri": "path/to/yet/another/song.mp3",
        "Name": "playcount",
        "Value": "4"
      },
      {
        "Uri": "path/to/another/track-02.ogg",
        "Name": "rating",
        "Value": "2"
      }
    ]

### CSV format

This looks like the following:

        track-01.ogg,playcount,1
        path/to/yet/another/song.mp3,playcount,4
        path/to/another/track-02.ogg,rating,2

# Things to improve

## mpdfavd

* Limit ratings to one rating per song per client, if possible. (Currently no way to identify clients who sent a rating, since we use mpd channels)

