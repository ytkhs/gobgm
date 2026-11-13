package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var (
	downloadedFiles []string
)

const (
	ItunesEndPoint = "https://itunes.apple.com/search/"
	PlayCmd        = "afplay"
)

type ITunesRequestParams struct {
	Term    string
	Country string
	Lang    string
	Media   string
	Entity  string
	Limit   string
}

func DefaultITunesRequestParams() (d ITunesRequestParams) {
	d = ITunesRequestParams{
		Term:    "",
		Country: "JP",
		Lang:    "ja_jp",
		Media:   "music",
		Entity:  "song",
		Limit:   "200",
	}
	return
}

type Music struct {
	ArtistName     string `json:"artistName"`
	TrackName      string `json:"trackName"`
	PreviewURL     string `json:"previewUrl"`
	CollectionName string `json:"collectionName"`
	TrackViewURL   string `json:"trackViewUrl"`
	ReleaseDate    string `json:"releaseDate"`
	Genre          string `json:"primaryGenreName"`
}

type ITunesResponse struct {
	Count   int    `json:"resultCount"`
	Results Musics `json:"results"`
}

type Musics []Music

func (self *Musics) Shuffle() {

	musics := *self

	rand.Shuffle(len(musics), func(i, j int) {
		musics[i], musics[j] = musics[j], musics[i]
	})
}

func SearchMusic(v ITunesRequestParams) <-chan ITunesResponse {

	resultChan := make(chan ITunesResponse)

	params := url.Values{}
	params.Add("term", v.Term)
	params.Add("country", v.Country)
	params.Add("lang", v.Lang)
	params.Add("media", v.Media)
	params.Add("entity", v.Entity)
	params.Add("limit", v.Limit)

	go func() {
		fmt.Println("* provided courtesy of iTunes *")
		fmt.Println()
		response, err := http.Get(ItunesEndPoint + "?" + params.Encode())

		if err != nil {
			log.Fatal(err)
			os.Exit(1)

		} else {
			defer response.Body.Close()

			var data ITunesResponse
			err = json.NewDecoder(response.Body).Decode(&data)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			resultChan <- data
		}

	}()

	return resultChan
}

func Play(fileName string, rate string) {
	defer os.Remove(fileName)

	out, _ := exec.Command(PlayCmd, fileName, "--rate", rate, "-q", "1").CombinedOutput()
	fmt.Print(string(out))
}

func Info(music Music, num int, total int) {
	fmt.Printf("* (%d/%d)\n", num, total)
	fmt.Printf("♪ %s - %s / %s%s\n", music.TrackName, music.ArtistName, music.CollectionName, music.details())
	fmt.Printf("# %s\n", music.TrackViewURL)
	fmt.Println()
}

// details formats as " (2012, J-Pop)", omitting fields the API left empty.
func (m Music) details() string {
	var parts []string
	if len(m.ReleaseDate) >= 4 {
		parts = append(parts, m.ReleaseDate[:4])
	}
	if m.Genre != "" {
		parts = append(parts, m.Genre)
	}
	if len(parts) == 0 {
		return ""
	}
	return " (" + strings.Join(parts, ", ") + ")"
}

func Download(url string) <-chan string {

	fileNameChan := make(chan string)

	go func(url string) {
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		file, err := os.CreateTemp("", "bgm_tmp_")
		if err != nil {
			log.Fatal(err)
		}
		downloadedFiles = append(downloadedFiles, file.Name())
		defer file.Close()

		io.Copy(file, response.Body)

		fileNameChan <- file.Name()

	}(url)

	return fileNameChan

}

func SweepFiles() {
	for _, file := range downloadedFiles {
		os.Remove(file)
	}
}
