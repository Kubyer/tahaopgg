package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type data struct {
	Value string `json:"value"`
}

type puuidResponse struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

var XRiotToken = "RGAPI-a238c290-1b85-4cbd-8826-e1287fcd8be1"

func getRiotPuuid(tagLine string, gameName string) (puuidResponse, error) {

	url := "https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/" + gameName + "/" + tagLine

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	var response puuidResponse
	if err != nil {
		fmt.Println(err)
		return response, err
	}
	req.Header.Add("X-Riot-Token", XRiotToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return response, err
	}
	defer res.Body.Close()

	// Check for successful response status code (e.g., 200)
	if res.StatusCode != http.StatusOK {
		fmt.Println("Error:", res.StatusCode)
		return response, err
	}
	// Create a decoder from the response body
	decoder := json.NewDecoder(res.Body)

	// Decode the JSON data from the response body
	err = decoder.Decode(&response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return response, err
	}
	fmt.Println(response.Puuid)

	return response, nil
}

func handleFormSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Handle non-POST requests (optional)
		return
	}

	// Create a decoder based on content type (assuming JSON)
	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	// Decode the request body into the struct
	var receivedData data
	fmt.Println("receivedData ", receivedData.Value)

	err := decoder.Decode(&receivedData)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	parts := strings.SplitN(receivedData.Value, "#", 2)
	// Handle cases with no "#" or just one element
	if len(parts) == 0 {
		fmt.Println("String is empty")
		return
	} else if len(parts) == 1 {
		gameName := parts[0]
		fmt.Println("gameName:", gameName)
		return
	}

	gameName := strings.Join(parts[:len(parts)-1], "#")
	tagline := parts[len(parts)-1]

	getRiotPuuid(tagline, gameName)
	return
}

func main() {
	staticDir := http.Dir("static")
	fileHandler := http.FileServer(staticDir)
	http.Handle("/", fileHandler)
	http.HandleFunc("/submit", handleFormSubmit)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
