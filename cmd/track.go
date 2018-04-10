package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sigafoos/lastfm"
	"github.com/spf13/cobra"
)

var (
	album  string
	artist string
)

func init() {
	rootCmd.AddCommand(trackScrobble)

	trackScrobble.Flags().StringVarP(&album, "album", "l", "", "The album of the track; if omitted, will be prompted for one")
	trackScrobble.Flags().StringVarP(&artist, "artist", "a", "", "The artist of the track; if omitted, will be prompted for one")
	trackScrobble.Flags().StringVarP(&title, "title", "t", "", "The title of the track; if omitted, will be prompted for one")
}

var trackScrobble = &cobra.Command{
	Use:   "track",
	Short: "Scrobble a single tracks on a track",
	Long:  "Scrobble a single tracks on a track",
	Run:   scrobbleTrack,
}

func scrobbleTrack(cmd *cobra.Command, args []string) {
	lfm.SetVerbose(verbose)

	if title == "" {
		fmt.Print("track title: ")
		title, _ = cli.ReadString('\n')
		title = strings.TrimSpace(title)
	}

	if artist == "" {
		fmt.Print("track artist: ")
		artist, _ = cli.ReadString('\n')
		artist = strings.TrimSpace(artist)
	}

	if album == "" {
		fmt.Print("track album: ")
		album, _ = cli.ReadString('\n')
		album = strings.TrimSpace(album)
	}

	track := lastfm.Track{
		API:    lfm,
		Name:   title,
		Artist: lastfm.Artist{Name: artist},
		Album:  album,
	}

	err := track.Scrobble()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
