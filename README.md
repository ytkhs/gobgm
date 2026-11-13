# gobgm

gobgm - iTunes preview player

## description

gobgm is based on a ruby tool "[bgm](https://github.com/hitode909/bgm)" and written in golang.
Download music preview files from [iTunes Search API](https://www.apple.com/itunes/affiliates/resources/documentation/itunes-store-web-service-search-api.html), and play it via `afplay` command on OSX.


## installation

```
go install github.com/ytkhs/gobgm@latest
```

## usage

```
usage: gobgm <term...> [-l|--list] [--shuffle] [-r|--rate num]

  term           search query
  -l, --list     show results without playing
  -r, --rate     playback rate (e.g. 0.5, 2)
      --shuffle  play in random order
  -h, --help     show this help

Press Ctrl-C to stop.
```


## example

```
gobgm SEKAI NO OWARI

gobgm SEKAI NO OWARI -r 0.5

gobgm SEKAI NO OWARI --shuffle

gobgm SEKAI NO OWARI --list
```

	
