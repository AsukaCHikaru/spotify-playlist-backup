package apiCore

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	TokenURL = "https://accounts.spotify.com/api/token"
)

var (
	client      = &http.Client{}
	accessToken string
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func authenticate() (AuthResponse, error) {
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

func getAccessToken() (string, error) {
	if accessToken == "" {
		var err error
		authResponse, err := authenticate()
		accessToken = authResponse.AccessToken
		if err != nil {
			return "", err
		}
	}
	return accessToken, nil
}

func Fetch(url string) (string, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(responseData), nil
}
