package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/Sigafoos/lastfm"
)

var (
	search      = flag.Bool("s", false, "use the search flow")
	verbose     = flag.Bool("v", false, "print mbid, curl calls, etc")
	requireMBID = flag.Bool("requirembid", true, "require a MusicBrainz ID")
	mbid        = flag.String("mbid", "", "a MusicBrainz ID to scrobble")
)

// APIKey and APISecret are the key and shared secret for a last.fm app. If you're building this from source
// you'll need your own, and to build it with a compile time linker flag:
//
//   go build -ldflags "-X main.APIKey=<your key> -X main.APISecret=<your secret>" -o scrobble main.go
var (
	APIKey    string
	APISecret string
)

type Scrobbler struct {
	cli *bufio.Reader
	lfm *lastfm.LastFM
}

func NewScrobbler(key, secret string, verbose bool) *Scrobbler {
	return &Scrobbler{
		cli: bufio.NewReader(os.Stdin),
		lfm: lastfm.New(key, secret, verbose),
	}
}

func (s *Scrobbler) GetSessionToken() string {
	usr, _ := user.Current()
	filename := fmt.Sprintf("%s/.lastfm", usr.HomeDir)
	b, err := ioutil.ReadFile(filename)

	if err == nil {
		return string(b)
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Println("Generating session token (your password will NOT be stored)...\n")
	fmt.Print("last.fm username: ")
	username, _ := s.cli.ReadString('\n')
	fmt.Print("last.fm password: ")
	password, _ := s.cli.ReadString('\n')

	token, err := s.lfm.Authenticate(strings.TrimSpace(username), strings.TrimSpace(password))
	if err != nil {
		fmt.Println(err)
		f.Close() // won't trigger defer on os.Exit()
		os.Remove(filename)
		os.Exit(1)
	}
	f.WriteString(token)
	f.Sync()

	return token
}

func main() {
	if APIKey == "" || APISecret == "" {
		fmt.Println("fatal error: API key/secret not found. was there a compile-time error?")
		os.Exit(1)
	}

	flag.Parse()
	if !*search && *mbid == "" {
		*search = true
	}

	scrobbler := NewScrobbler(APIKey, APISecret, *verbose)

	token := scrobbler.GetSessionToken()

	scrobbler.lfm.SessionKey(token)

	var album lastfm.Album

	if *search {
		fmt.Print("album title: ")
		title, _ := scrobbler.cli.ReadString('\n')
		title = strings.TrimSpace(title)

		albums, err := scrobbler.lfm.AlbumSearch(title, *requireMBID)
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
		text, _ := scrobbler.cli.ReadString('\n')
		text = strings.TrimSpace(text)
		i, err := strconv.Atoi(text)

		if err != nil || i < 0 || i > 5 || i > len(albums) {
			fmt.Println("invalid entry")
			os.Exit(1)
		} else if i == 0 {
			os.Exit(0)
		}

		album = albums[i-1]
	} else {
		album = lastfm.Album{MusicBrainzID: *mbid}
	}

	if *verbose {
		fmt.Printf("MBID: %s\n", album.MusicBrainzID)
	}

	album, err := scrobbler.lfm.AlbumInfo(album)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Scrobbling %v tracks from %s - \"%s\"\n\n", len(album.TrackList.Tracks), album.Artist, album.Name)
	err = scrobbler.lfm.ScrobbleAlbum(album)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
