package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PuuidResponse struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// GetRiotPuuid fetches the puuid based on game name and tagline from the Riot API.
func GetRiotPuuid(tagLine string, gameName string) (PuuidResponse, error) {
	url := "https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/" + gameName + "/" + tagLine

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return PuuidResponse{}, err
	}
	req.Header.Add("X-Riot-Token", "RGAPI-733e5712-4df5-4799-9492-a69a98e3500d") // "generic key",  /! to hide later

	res, err := client.Do(req)
	if err != nil {
		return PuuidResponse{}, err
	}
	defer res.Body.Close()

	// Check for successful response status code (e.g., 200)
	if res.StatusCode != http.StatusOK {
		return PuuidResponse{}, fmt.Errorf("error: unexpected status code %d", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	var response PuuidResponse
	err = decoder.Decode(&response)
	if err != nil {
		return PuuidResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	return response, nil
}
