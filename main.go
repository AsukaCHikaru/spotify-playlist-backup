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

type PlaylistItemList struct {
	Items []TrackItem `json:"items"`
}

type TrackItem struct {
	Track Track `json:"track"`
}

type Track struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

type Playlist2 struct {
	Name  string
	Songs []Track
}

type Artist struct {
	Name string `json:"name"`
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

	for i := range list.Items {
		playlistApiEndpoint := list.Items[i].Url
		playlistItemsResponse := apiCore.Fetch(playlistApiEndpoint+"/tracks", client, authResponse.AccessToken)
		var list2 PlaylistItemList
		err = json.Unmarshal([]byte(playlistItemsResponse), &list2)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var result Playlist2
		result.Name = list.Items[i].Name
		for i := range list2.Items {
			result.Songs = append(result.Songs, list2.Items[i].Track)
		}
		fmt.Println((result))
	}
}
