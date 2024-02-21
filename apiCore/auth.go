package apiCore

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	TokenURL = "https://accounts.spotify.com/api/token"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func Authenticate(client http.Client) (AuthResponse, error) {
	body := url.Values{}
	body.Set("grant_type", "client_credentials")
	encodedBody := body.Encode()

	req, err := http.NewRequest("POST", TokenURL, bytes.NewBuffer([]byte(encodedBody)))
	if err != nil {
		fmt.Println(err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(os.Getenv("CLIENT_ID")+":"+os.Getenv("CLIENT_SECRET"))))

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer response.Body.Close()

	var authResponse AuthResponse
	err = json.NewDecoder(response.Body).Decode(&authResponse)
	if err != nil {
		return AuthResponse{}, err
	}

	return authResponse, err
}
