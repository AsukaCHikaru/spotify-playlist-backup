package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	auth "spotify-playlist-backup/apiCore"

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

	client := getHttpClient()
	authResponse, _ := auth.Authenticate(client)
	fmt.Println(fetch(endpoint, client, authResponse.AccessToken))
}

func getHttpClient() http.Client {
	client := new(http.Client)
	return *client
}

func fetch(url string, client http.Client, bearer string) string {
	req, _ := http.NewRequest("GET", getEndpoint(), nil)
	req.Header.Set("Authorization", "Bearer "+bearer)

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(responseData)
}
