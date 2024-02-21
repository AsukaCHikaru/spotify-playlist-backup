package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"spotify-playlist-backup/apiCore"

	"github.com/joho/godotenv"
)

var endpoint = "https://api.spotify.com/v1/users/{id}/playlists?limit=50"

func getEndpoint() string {
	return strings.Replace(endpoint, "{id}", os.Getenv("USER_ID"), 1)
}

type PlaylistItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"href"`
}

type Playlists struct {
	Items []PlaylistItem `json:"items"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := apiCore.GetHttpClient()
	authResponse, _ := apiCore.Authenticate(client)
	playlistResponse := apiCore.Fetch(getEndpoint(), client, authResponse.AccessToken)
	var list Playlists
	err = json.Unmarshal([]byte(playlistResponse), &list)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
