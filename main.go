package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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
	fmt.Println(fetch(endpoint, client))
}

func getHttpClient() http.Client {
	client := new(http.Client)
	return *client
}

func fetch(url string, client http.Client) string {
	req, _ := http.NewRequest("GET", getEndpoint(), nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("BEARER"))

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
