package apiCore

import (
	"fmt"
	"io"
	"net/http"
)

func GetHttpClient() http.Client {
	client := new(http.Client)
	return *client
}

func Fetch(url string, client http.Client, bearer string) string {
	req, _ := http.NewRequest("GET", url, nil)
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
