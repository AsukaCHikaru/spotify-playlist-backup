package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"spotify-playlist-backup/apiCore"
	"spotify-playlist-backup/parser"

	"github.com/joho/godotenv"
)

var endpoint = "https://api.spotify.com/v1/users/{id}/playlists?limit=50"

func getEndpoint() string {
	return strings.Replace(endpoint, "{id}", os.Getenv("USER_ID"), 1)
}

type Snapshot struct {
	Playlists []parser.Playlist
	UpdatedAt string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := apiCore.GetHttpClient()
	authResponse, _ := apiCore.Authenticate(client)
	playlistResponse := apiCore.Fetch(getEndpoint(), client, authResponse.AccessToken)

	playlistMeta, err := parser.ParseUserPlaylists(playlistResponse)
	if err != nil {
		fmt.Println(err.Error())
	}

	result := Snapshot{UpdatedAt: time.Now().Format("2006-01-02 15:04")}

	for i := range playlistMeta.Items {
		item := playlistMeta.Items[i]
		playlistItemsResponse := apiCore.Fetch(item.Tracks.Url, client, authResponse.AccessToken)
		playlist, err := parser.ParsePlaylistItems(playlistItemsResponse, item.Name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result.Playlists = append(result.Playlists, playlist)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = os.WriteFile("output.json", []byte(string(jsonData)), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
