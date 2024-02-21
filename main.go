package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"spotify-playlist-backup/apiCore"

	"github.com/joho/godotenv"
)

var endpoint = "https://api.spotify.com/v1/users/{id}/playlists"

func getEndpoint() string {
	return strings.Replace(endpoint, "{id}", os.Getenv("USER_ID"), 1)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := apiCore.GetHttpClient()
	authResponse, _ := apiCore.Authenticate(client)
	fmt.Println(apiCore.Fetch(getEndpoint(), client, authResponse.AccessToken))
}
