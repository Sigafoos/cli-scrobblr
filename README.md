# scrobble
A command line interface to last.fm

## Installing
If you'd like to use the app as-is, download the latest [release](https://github.com/Sigafoos/scrobble/releases). You will be prompted for your last.fm credentials on the first run. These _are not stored_: they're exchanged with last.fm for a session token, which is saved in `$HOME/.lastfm`

## Compiling
You'll need a [last.fm api account](https://www.last.fm/api). When compiling/running you'll need to pass them in as flags:

	go build -ldflags "-X github.com/Sigafoos/scrobble/cmd.APIKey=<YourAPIKey> -X github.com/Sigafoos/scrobble/cmd.APISecret=<YourAPISecret>" -o scrobble main.go

An api account is not necessary for using the compiled versions.

## Roadmap
* v0.1 Provide a release so that the link above works
* v0.2 Create an actual roadmap
