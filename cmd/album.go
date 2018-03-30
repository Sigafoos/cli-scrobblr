package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Sigafoos/lastfm"
	"github.com/spf13/cobra"
)

var (
	requireMBID bool
	mbid        string
)

func init() {
	rootCmd.AddCommand(albumScrobble)

	rootCmd.Flags().BoolVarP(&requireMBID, "requireMBID", "r", true, "Require results to have a MusicBrainz ID")
	rootCmd.Flags().StringVarP(&mbid, "mbid", "m", "", "A MusicBrainz ID to use instead of searching (note that not all mbids are in Last.FM)")
}

var albumScrobble = &cobra.Command{
	Use:   "album",
	Short: "Scrobble all tracks on an album",
	Long: `Search for an album, either by title or MusicBrainz ID.

Then retrieve the tracks from Last.FM and scrobble them all in a bulk command.`,
	Run: scrobbleAlbum,
}

func scrobbleAlbum(cmd *cobra.Command, args []string) {
	lfm.SetVerbose(verbose)

	var album lastfm.Album

	if mbid != "" {
		album = lastfm.Album{
			API:           lfm,
			MusicBrainzID: mbid,
		}
	} else {
		fmt.Print("album title: ")
		title, _ := cli.ReadString('\n')
		title = strings.TrimSpace(title)

		albums, err := lfm.AlbumSearch(title, requireMBID)
		if err != nil {
			fmt.Println(err)
		}

		if len(albums) == 0 {
			fmt.Printf("\nNo results found for \"%s\"\n", title)
			os.Exit(0)
		}

		fmt.Println()
		for i := 0; i < 5 && i < len(albums); i++ {
			album := albums[i]
			fmt.Printf("[%v] \"%s\" - %s\n", i+1, album.Name, album.Artist)
		}

		fmt.Print("\nChoose a result (0 to exit): ")
		text, _ := cli.ReadString('\n')
		text = strings.TrimSpace(text)
		i, err := strconv.Atoi(text)

		if err != nil || i < 0 || i > 5 || i > len(albums) {
			fmt.Println("invalid entry")
			os.Exit(1)
		} else if i == 0 {
			os.Exit(0)
		}

		album = albums[i-1]
	}

	if verbose {
		fmt.Printf("MBID: %s\n", album.MusicBrainzID)
	}

	err := album.GetInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Scrobbling %v tracks from %s - \"%s\"\n", len(album.TrackList.Tracks), album.Artist, album.Name)
	err = album.Scrobble()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
