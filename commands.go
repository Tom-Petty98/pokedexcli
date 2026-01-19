package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
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
		"explore": {
			name:        "explore",
			description: "Displays all pokemon in a location area. Usage: explore <area_name>",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon. Usage: catch <pokemon_name>",
			callback:    commandCatch,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays all the pokemon you have caught",
			callback:    commandPokedex,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays the details of a pokemon you have caught. Usage: inspect <pokemon_name>",
			callback:    commandInspect,
		},
	}
}

func commandHelp(args []string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key, command := range commands {
		fmt.Printf("%s: %s\n", key, command.description)
	}
	return nil
}

func commandExit(args []string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}
