package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"spotify-playlist-backup/apiCore"

	"github.com/joho/godotenv"
)

var endpoint = "https://api.spotify.com/v1/users/{id}/playlists?limit=2"

func getEndpoint() string {
	return strings.Replace(endpoint, "{id}", os.Getenv("USER_ID"), 1)
}

type UserPlaylistsResponse struct {
	Items []struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Tracks struct {
			Url string `json:"href"`
		} `json:"tracks"`
	} `json:"items"`
}

type PlaylistItemsResponse struct {
	Items []struct {
		Track Track `json:"track"`
	} `json:"items"`
}

type Track struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}
type Artist struct {
	Name string `json:"name"`
}
type Playlist struct {
	Name  string
	Songs []Track
}
type Snapshot struct {
	Playlists []Playlist
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
	var list UserPlaylistsResponse
	err = json.Unmarshal([]byte(playlistResponse), &list)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result := Snapshot{UpdatedAt: time.Now().Format("2006-01-02 15:04")}

	for i := range list.Items {
		item := list.Items[i]
		playlistItemsResponse := apiCore.Fetch(item.Tracks.Url, client, authResponse.AccessToken)
		var playlistItems PlaylistItemsResponse
		err = json.Unmarshal([]byte(playlistItemsResponse), &playlistItems)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var playlist Playlist
		playlist.Name = item.Name
		for i := range playlistItems.Items {
			playlist.Songs = append(playlist.Songs, playlistItems.Items[i].Track)
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
