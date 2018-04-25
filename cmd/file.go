package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sigafoos/lastfm"
	"github.com/spf13/cobra"
)

var (
	filepath string
)

func init() {
	rootCmd.AddCommand(fileScrobble)

	fileScrobble.Flags().StringVarP(&filepath, "file", "f", "", "The path to the file to be scrobbled (required)")
	fileScrobble.MarkFlagRequired("file")
}

var fileScrobble = &cobra.Command{
	Use:   "file",
	Short: "Scrobble tracks from a file",
	Long: `Scrobble tracks from a file. Each line is one track, with tab-delimited fields.

The fields are, in order: track name, track album, track artist`,
	Run: scrobbleFile,
}

func scrobbleFile(cmd *cobra.Command, args []string) {
	lfm.SetVerbose(verbose)

	fp, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var lines []string
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// not deferring because os.Exit() won't trigger it
	fp.Close()

	for _, line := range lines {
		values := strings.Split(line, "	")
		track := lastfm.Track{
			API:    lfm,
			Name:   values[0],
			Artist: lastfm.Artist{Name: values[2]},
			Album:  values[1],
		}

		fmt.Printf("Scrobbling %s: \"%s\" (%s)\n", values[2], values[0], values[1])
		err := track.Scrobble()
		if err != nil {
			fmt.Printf("*** ERROR with %s: %s\n", values[0], err)
		}
		time.Sleep(1 * time.Second)
	}
}
