package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	pokecache "github.com/Tom-Petty98/pokedexcli/internal"
)

var cache *pokecache.Cache

func main() {
	userInput := bufio.NewScanner(os.Stdin)
	cache = pokecache.NewCache(1 * time.Minute)

	for {
		fmt.Print("Pokedex > ")
		if !userInput.Scan() {
			break
		}

		text := userInput.Text()
		words := cleanInput(text)
		if len(words) > 0 {
			value, exists := commands[words[0]]
			if exists {
				args := words[1:]

				if err := value.callback(args); err != nil {
					fmt.Printf("%s\n", err)
				}

			} else {
				fmt.Print("Unknown command\n")
			}
		}
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(strings.TrimSpace(text))
	return strings.Fields(text)
}
