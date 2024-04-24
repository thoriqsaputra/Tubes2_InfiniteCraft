package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: go run main.go <option> [<start_article_title> <target_article_title>]")
		fmt.Println("Option: 1 for your main function, 2 for your friend's main function")
		return
	}

	option := args[0]
	startArticle := args[1]
    targetArticle := args[2]

	switch option {
	case "1":
		// Hardcode the start and end articles
		path, numChecked, duration, err := bfs(startArticle, targetArticle)
		if err != nil {
			log.Fatal(err)
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Number of articles checked: %d\n", numChecked)
		fmt.Printf("Number of articles traversed: %d\n", len(path))
		fmt.Printf("Route: %s\n", strings.Join(path, " -> "))
		fmt.Printf("Time taken: %v\n", duration)

	case "2":
		start := time.Now()
		path := IDS(startArticle,targetArticle, 3)

		elapsed := time.Since(start)
		fmt.Println("Time taken:", elapsed.Milliseconds(), "ms")

		// Jumlah artikel yang diperiksa
		fmt.Println("Jumlah artikel yang diperiksa:", articlesChecked)

		// Jumlah artikel yang dilalui
		if path != nil {
			fmt.Println("Jumlah artikel yang dilalui:", len(path)-1)
		} else {
			fmt.Println("Jumlah artikel yang dilalui: 0")
		}

		// Rute
		if path != nil {
			fmt.Println("Rute:", strings.Join(path, " -> "))
		} else {
			fmt.Println("Rute: No path found")
		}

	default:
		fmt.Println("Invalid option. Please choose 1 for your main function, 2 for your friend's main function.")
	}
}
