// Package cmd is a command line interface for scrobbling last.fm tracks.
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/Sigafoos/lastfm"

	"github.com/spf13/cobra"
)

// APIKey and APISecret are the key and shared secret for a last.fm app. If you're building this from source
// you'll need your own, and to build it with a compile time linker flag:
//
//   go build -ldflags "-X github.com/Sigafoos/scrobble/cmd.APIKey=<your key> -X github.com/Sigafoos/scrobble/cmd.APISecret=<your secret>" -o scrobble github.com/Sigafoos/scrobble/cmd.go
var (
	APIKey    string
	APISecret string
)

var (
	lfm *lastfm.API
	cli *bufio.Reader
)

// command line options
var (
	verbose     bool
	requireMBID bool
	mbid        string
	title       string
	album       string
	artist      string
)

var rootCmd = &cobra.Command{
	Use:   "scrobble",
	Short: "Scrobble tracks and albums to Last.FM",
	Long:  `Scrobble tracks and albums to Last.FM`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

// Execute runs the root command of the app.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	if APIKey == "" || APISecret == "" {
		fmt.Println("fatal error: API key/secret not found. was there a compile-time error?")
		os.Exit(1)
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	cli = bufio.NewReader(os.Stdin)
	lfm = lastfm.New(APIKey, APISecret)

	token, err := getSessionToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lfm.SetSessionKey(token)
}

func getSessionToken() (token string, err error) {
	usr, _ := user.Current()
	filename := fmt.Sprintf("%s/.lastfm", usr.HomeDir)
	b, err := ioutil.ReadFile(filename)

	// if the file exists, read it and return
	if err == nil {
		token = string(b)
		if token != "" {
			return
		}
		// it was empty, so remove the file and re-create
		_ = os.Remove(filename)
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	// if we can't chmod it don't error out, but let the user know
	if err = os.Chmod(filename, os.FileMode(0600)); err != nil {
		fmt.Printf("warning: %s could not be given restricted permissions\n", filename)
	}

	fmt.Println("Generating session token (your password will NOT be stored)...\n")
	fmt.Print("last.fm username: ")
	username, _ := cli.ReadString('\n')
	fmt.Print("last.fm password: ")
	password, _ := cli.ReadString('\n')

	token, err = lfm.Authenticate(strings.TrimSpace(username), strings.TrimSpace(password))
	if err != nil {
		fmt.Println(err)
		f.Close() // won't trigger defer on os.Exit()
		os.Remove(filename)
		os.Exit(1)
	} else if token == "" {
		fmt.Println("token returned from last.fm was empty")
		f.Close() // won't trigger defer on os.Exit()
		os.Remove(filename)
		os.Exit(1)
	}

	_, err = f.WriteString(token)
	if err != nil {
		return "", err
	}
	err = f.Sync()
	if err != nil {
		return "", err
	}

	return
}
