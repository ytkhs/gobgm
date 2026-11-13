package main

import "fmt"

func ShowHelp() {
	fmt.Print(`usage: gobgm <term...> [-l|--list] [--shuffle] [-r|--rate num]

  term           search query
  -l, --list     show results without playing
  -r, --rate     playback rate (e.g. 0.5, 2)
      --shuffle  play in random order
  -h, --help     show this help

Press Ctrl-C to stop.
`)
}
