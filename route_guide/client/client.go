package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	gatewayURL = "http://localhost:8081/getFeature"
)

type Point struct {
	Latitude  int32 `json:"latitude"`
	Longitude int32 `json:"longitude"`
}

type Feature struct {
	Name     string `json:"name"`
	Location Point  `json:"location"`
	Area     string `json:"area"`
}

func main() {
	point := Point{
		Latitude:  410248224,
		Longitude: -747127767,
	}

	requestBody, err := json.Marshal(&point)

	if err != nil {
		log.Fatalf("could not marshal request body: %v", err)
	}

	response, err := http.Post(gatewayURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		log.Fatalf("could not make request to gateway: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", response.StatusCode)
	}

	var feature Feature
	if err := json.NewDecoder(response.Body).Decode(&feature); err != nil {
		log.Fatalf("could not decode response: %v", err)
	}
	log.Println(feature)

}
