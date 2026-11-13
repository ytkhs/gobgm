package main

import (
	"fmt"
	"os"
	"os/exec"
)

func init() {
	RegisterExitProcess()

	if _, err := exec.LookPath(PlayCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {

	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			ShowHelp()
			return
		}
	}

	if len(os.Args) == 1 || os.Args[1][0:1] == "-" {
		fmt.Fprintln(os.Stderr, "please input term")
		ShowHelp()
		os.Exit(1)
	}

	query, rate, shuffle, list := ProcessArgs()

	params := DefaultITunesRequestParams()
	params.Term = query

	result := <-SearchMusic(params)

	if len(result.Results) == 0 {
		fmt.Fprintf(os.Stderr, "no results for %q\n", query)
		os.Exit(1)
	}

	if shuffle {
		result.Results.Shuffle()
	}

	if list {
		InfoAll(result.Results)
		return
	}

	for i, music := range result.Results {
		Info(music, i+1, len(result.Results))
		Play(<-Download(music.PreviewURL), rate)
	}

}

func InfoAll(musics []Music) {
	for i, music := range musics {
		Info(music, i+1, len(musics))
	}
}

func ProcessArgs() (query string, rate string, shuffle bool, list bool) {

	rate = "1"
	shuffle = false
	list = false

	query = os.Args[1]

	hasOption := false

	for i := 2; i < len(os.Args); i++ {

		v := os.Args[i]

		if !hasOption {
			if v[0:1] == "-" {
				hasOption = true
			} else {
				query = query + " " + v
				continue
			}
		}

		switch v {
		case "--rate", "-r":
			if i+1 < len(os.Args) {
				rate = os.Args[i+1]
				i++
			}
		case "--shuffle":
			shuffle = true
		case "--list", "-l":
			list = true
		default:
			fmt.Fprintf(os.Stderr, "unknown option %s (see --help)\n", v)
			os.Exit(1)
		}

	}

	return
}
