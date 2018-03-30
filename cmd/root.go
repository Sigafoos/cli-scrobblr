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
	lfm     *lastfm.API
	cli     *bufio.Reader
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "scrobble",
	Short: "Scrobble tracks and albums to Last.FM",
	Long: `something coming soon
	`,
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
	lfm.SetSessionKey(getSessionToken())
}

func getSessionToken() string {
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
	username, _ := cli.ReadString('\n')
	fmt.Print("last.fm password: ")
	password, _ := cli.ReadString('\n')

	token, err := lfm.Authenticate(strings.TrimSpace(username), strings.TrimSpace(password))
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
