package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"spotify-playlist-backup/pkg/fetch"
	"spotify-playlist-backup/pkg/parser"

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

	playlistMeta, err := getUserPlaylistMeta()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	playlists, err := getPlaylists(playlistMeta)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result := Snapshot{
		UpdatedAt: time.Now().Format("2006-01-02 15:04"),
		Playlists: playlists,
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

func getUserPlaylistMeta() (parser.UserPlaylistsResponse, error) {
	playlistResponse, err := fetch.Fetch(getEndpoint())
	if err != nil {
		fmt.Println(err.Error())
		return parser.UserPlaylistsResponse{}, err
	}

	playlistMeta, err := parser.ParseUserPlaylists(playlistResponse)
	if err != nil {
		return parser.UserPlaylistsResponse{}, err
	}
	return playlistMeta, nil
}

func getPlaylists(meta parser.UserPlaylistsResponse) ([]parser.Playlist, error) {
	var playlists []parser.Playlist
	for i := range meta.Items {
		item := meta.Items[i]
		playlistItemsResponse, err := fetch.Fetch(item.Tracks.Url)
		if err != nil {
			return []parser.Playlist{}, err
		}
		playlist, err := parser.ParsePlaylistItems(playlistItemsResponse, item.Name)
		if err != nil {
			return []parser.Playlist{}, err
		}

		playlists = append(playlists, playlist)
	}
	return playlists, nil
}
