package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of previous 20 location areas",
			callback:    commandMapb,
		},
	}
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key, command := range commands {
		fmt.Printf("%s: %s\n", key, command.description)
	}
	return nil
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

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

func commandMap() error {
	if prevPage != nil {
		nextPage++
	}
	return outputLocationData()
}

func commandMapb() error {
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
