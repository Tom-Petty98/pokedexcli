package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
)

type PokemonStat struct {
	Name string `json:"name"`
}

type PokemonStats struct {
	StatValue int         `json:"base_stat"`
	Stat      PokemonStat `json:"stat"`
}

type PokemonType struct {
	Name string `json:"name"`
}

type PokemonTypes struct {
	Details PokemonType `json:"type"`
}

type Pokemon struct {
	Name           string         `json:"name"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Stats          []PokemonStats `json:"stats"`
	Types          []PokemonTypes `json:"types"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

type LocationAreaDetails struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

func commandExplore(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing location_area example usage: explore <location_area>")
	}

	locationName := args[0]
	fmt.Printf("Exploring %s... \n", locationName)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", locationName)

	value, exists := cache.Get(url)
	if !exists {
		res, _ := http.Get(url)
		if res.StatusCode == 404 {
			return fmt.Errorf("Could not find location area with name: %s", locationName)
		} else if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, res.Body)
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		value = body
		cache.Add(url, value)
	}

	locationAreaDetails := LocationAreaDetails{}
	err := json.Unmarshal(value, &locationAreaDetails)
	if err != nil {
		return err
	}

	fmt.Print("Found Pokemon: \n")
	count := len(locationAreaDetails.PokemonEncounters)
	for i := 0; i < count; i++ {
		fmt.Printf("%s\n", locationAreaDetails.PokemonEncounters[i].Pokemon.Name)
	}
	return nil
}

var pokedex = make(map[string]Pokemon)

func commandCatch(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing pokemon_name example usage: explore <pokemon_name>")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s... \n", pokemonName)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)

	value, exists := cache.Get(url)
	if !exists {
		res, _ := http.Get(url)
		if res.StatusCode == 404 {
			return fmt.Errorf("Could not find pokemon with name: %s", pokemonName)
		} else if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, res.Body)
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		value = body
		cache.Add(url, value)
	}

	pokemonDetails := Pokemon{}
	err := json.Unmarshal(value, &pokemonDetails)
	if err != nil {
		return err
	}

	attemptCatchWeightedRoll(&pokemonDetails)

	return nil
}

// func attemptCatchLinear(pokemon *Pokemon) {
// 	const maxDifficulty = 325 // googled hardest to catch and input each into api arceus was highest
// 	catchChance := 1 - float64(pokemon.BaseExperience)/float64(maxDifficulty)
// 	roll := rand.Float64()

// 	if roll < catchChance {
// 		fmt.Printf("%s was caught!\n", pokemon.Name)
// 		pokedex[pokemon.Name] = *pokemon
// 	} else {
// 		fmt.Printf("%s escaped!\n", pokemon.Name)
// 	}
// }

func attemptCatchWeightedRoll(pokemon *Pokemon) {
	roll := rand.IntN(pokemon.BaseExperience + 1)
	playerStrength := 30 // could use some sort of algorithm like num of pokemon caught as player experience could add modifiers

	if roll < playerStrength {
		fmt.Printf("%s was caught and has been added to your pokedex!\n", pokemon.Name)
		pokedex[pokemon.Name] = *pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
}

func commandPokedex(args []string) error {
	fmt.Print("Your Pokedex:\n")
	for key := range pokedex {
		fmt.Printf(" - %s\n", key)

	}
	fmt.Print("Use the inspect <pokemon_name> command to view pokemon details\n")
	return nil
}

func commandInspect(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing pokemon_name example usage: explore <pokemon_name>")
	}

	pokemonName := args[0]

	if _, exists := pokedex[pokemonName]; !exists {
		return fmt.Errorf("You have not found %s", pokemonName)
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)

	value, exists := cache.Get(url)
	if !exists {
		res, _ := http.Get(url)
		if res.StatusCode == 404 {
			return fmt.Errorf("Could not find pokemon with name: %s", pokemonName)
		} else if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, res.Body)
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		value = body
		cache.Add(url, value)
	}

	pokemonDetails := Pokemon{}
	err := json.Unmarshal(value, &pokemonDetails)
	if err != nil {
		return err
	}

	printPokemonDetails(&pokemonDetails)

	return nil
}

func printPokemonDetails(pokemon *Pokemon) {
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", pokemon.Name, pokemon.Weight, pokemon.Height)

	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.StatValue)
	}

	fmt.Print("Types:\n")
	for _, t := range pokemon.Types {
		fmt.Printf("  -%s\n", t.Details.Name)
	}
}
