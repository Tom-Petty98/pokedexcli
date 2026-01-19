package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationData struct {
	Count   int            `json:"count"`
	Next    string         `json:"next"`
	Prev    *string        `json:"prev"`
	Results []LocationArea `json:"results"`
}

var nextPage int = 0
var prevPage *int

func commandMap(args []string) error {
	if prevPage != nil {
		nextPage++
	}
	return outputLocationData()
}

func commandMapb(args []string) error {
	if nextPage == 0 {
		fmt.Print("you're on the first page\n")
		return nil
	}
	nextPage--
	return outputLocationData()
}

func outputLocationData() error {
	prevPage = &nextPage
	offset := nextPage * 20
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=20", offset)

	value, exists := cache.Get(url)
	if !exists {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		value = body
		cache.Add(url, value)
	}

	locations := LocationData{}
	err := json.Unmarshal(value, &locations)
	if err != nil {
		return err
	}

	for i := 0; i < 20; i++ {
		fmt.Printf("%s\n", locations.Results[i].Name)
	}
	return nil
}
