// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// 	"time"
// )

// func main() {
// 	args := os.Args[1:]
// 	if len(args) < 1 {
// 		fmt.Println("Usage: go run main.go <option> [<start_article_title> <target_article_title>]")
// 		fmt.Println("Option: 1 for your main function, 2 for your friend's main function")
// 		return
// 	}

// 	option := args[0]
// 	startArticle := args[1]
//     targetArticle := args[2]

// 	switch option {
// 	case "1":
// 		// Hardcode the start and end articles
// 	startURL := ArticleURL(startArticle)
// 	endURL := ArticleURL(targetArticle)

// 	var path []string
// 	var duration time.Duration
// 	var err error
// 	var links int


// 	bfsInstance := NewBase(startURL, endURL)
// 	startTime := time.Now()
// 	path, err = bfsInstance.Bfs()
// 	links = bfsInstance.Visit()
// 	duration = time.Since(startTime)

	
// 	if err != nil {//
// 		log.Fatalf("Error finding path: %v", err)
// 	} else {
// 		fmt.Println()
// 		fmt.Printf("Jumlah artikel yang diperiksa: %d\n", links)
// 		fmt.Println("Jumlah artikel yang dilalui: ", len(path)-1)
// 		fmt.Println("Path route:")
// 		for i := 0; i < len(path)-1; i++ {
// 			parts := strings.Split(path[i], "/")
// 			rute := parts[len(parts)-1]
			
// 			if i < len(path)-2 {
// 				fmt.Printf("%s -> ", rute)
// 			} else {
// 				parts2 := strings.Split(path[i+1], "/")
// 				rute2 := parts2[len(parts2)-1]
// 				fmt.Printf("%s -> %s", rute, rute2)
// 			}
// 		}
// 		fmt.Println()
// 		fmt.Printf("Time Taken (ms): %v\n", duration.Milliseconds())
// 	}

// 	case "2":
// 		start := time.Now()
// 		path := IDS(startArticle,targetArticle, 3)

// 		elapsed := time.Since(start)
// 		fmt.Println("Time taken:", elapsed.Milliseconds(), "ms")

// 		// Jumlah artikel yang diperiksa
// 		fmt.Println("Jumlah artikel yang diperiksa:", articlesChecked)

// 		// Jumlah artikel yang dilalui
// 		if path != nil {
// 			fmt.Println("Jumlah artikel yang dilalui:", len(path)-1)
// 		} else {
// 			fmt.Println("Jumlah artikel yang dilalui: 0")
// 		}

// 		// Rute
// 		if path != nil {
// 			fmt.Println("Rute:", strings.Join(path, " -> "))
// 		} else {
// 			fmt.Println("Rute: No path found")
// 		}

// 	default:
// 		fmt.Println("Invalid option. Please choose 1 for your main function, 2 for your friend's main function.")
// 	}
// }