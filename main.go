package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/asukachikaru/spotify-playlist-backup/pkg/fetch"
	"github.com/asukachikaru/spotify-playlist-backup/pkg/parser"
	"github.com/asukachikaru/spotify-playlist-backup/pkg/write"

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

	write.WriteToJson(result)
	fmt.Println("Backup completed!")
}

func getUserPlaylistMeta() (parser.UserPlaylistsResponse, error) {
	fmt.Println("Fetching users' playlists...")
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
		fmt.Printf("Fetching playlist %v...", item.Name)
		playlistItemsResponse, err := fetch.Fetch(item.Tracks.Url)
		if err != nil {
			return []parser.Playlist{}, err
		}
		playlist, err := parser.ParsePlaylistItems(playlistItemsResponse, item.Name)
		if err != nil {
			return []parser.Playlist{}, err
		}

		fmt.Printf("complete\n")
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}
