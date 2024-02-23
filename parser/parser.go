package parser

import (
	"encoding/json"
)

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

func ParseUserPlaylists(response string) (UserPlaylistsResponse, error) {
	var list UserPlaylistsResponse
	var err = json.Unmarshal([]byte(response), &list)
	if err != nil {
		return list, err
	}

	return list, nil
}

func ParsePlaylistItems(response string, playlistName string) (Playlist, error) {
	var playlistItems PlaylistItemsResponse
	var playlist Playlist
	var err = json.Unmarshal([]byte(response), &playlistItems)
	if err != nil {
		return playlist, err
	}
	playlist.Name = playlistName
	for i := range playlistItems.Items {
		playlist.Songs = append(playlist.Songs, playlistItems.Items[i].Track)
	}
	return playlist, nil
}
