package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	cache      = make(map[string][]string)
	cacheMutex sync.Mutex
)
var done = make(chan struct{}) // done adalah channel yang akan ditutup ketika semua pekerjaan selesai.

// fetchPageFromCache mengambil halaman dari cache jika tersedia.
func fetchPageFromCache(url string) ([]string, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if links, ok := cache[url]; ok {
		return links, true
	}
	return nil, false
}

// savePageToCache menyimpan halaman ke dalam cache.
func savePageToCache(url string, links []string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	cache[url] = links
}

// fetchPage mengambil konten HTML dari URL menggunakan goquery.
func fetchPage(articleTitle string) ([]string, error) {
	// Construct the full URL for the Wikipedia article
	url := "https://en.wikipedia.org/wiki/" + articleTitle

	if links, ok := fetchPageFromCache(url); ok {
		return links, nil
	}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var links []string

	doc.Find("p a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":") {
			links = append(links, strings.TrimPrefix(link, "/wiki/"))
		}
	})

	savePageToCache(url, links)
	return links, nil
}

// worker adalah fungsi yang akan dijalankan oleh setiap goroutine dalam pool.
func worker(jobs <-chan string, results chan<- []string) {
	for job := range jobs {
		links, err := fetchPage(job)
		if err != nil {
			log.Printf("Error fetching page for %s: %v", job, err)
			continue
		}
		results <- links
	}
}

// bfs performs a breadth-first search from startTitle to targetTitle
// bfs performs a breadth-first search from startTitle to targetTitle
func bfs(startTitle, targetTitle string) ([]string, int, time.Duration, error) {
	queue := [][]string{{startTitle}}
	visited := make(map[string]bool)
	startTime := time.Now()

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Tambahkan goroutine untuk menunggu sinyal bahwa semua pekerjaan telah selesai
	go func() {
		wg.Wait()
		close(done)
	}()

	for len(queue) > 0 {
		mu.Lock()
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]
		mu.Unlock()

		fmt.Println("Checking:", node)

		if node == targetTitle {
			return path, len(visited), time.Since(startTime), nil
		}

		if !visited[node] {
			visited[node] = true
			links, err := fetchPage(node)
			if err != nil {
				return nil, 0, 0, err
			}

			wg.Add(len(links))
			for _, link := range links {
				go func(link string) {
					defer wg.Done()
					newPath := append([]string{}, path...)
					newPath = append(newPath, link)
					mu.Lock()
					queue = append(queue, newPath)
					mu.Unlock()
				}(link)
			}
		}
	}

	return nil, 0, 0, fmt.Errorf("path not found")
}

// fetchPath fetches the path from startTitle to targetTitle
func fetchPath(startTitle, targetTitle string, visited map[string]bool) []string {
	queue := [][]string{{startTitle}}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]

		fmt.Println("Checking:", node)

		if node == targetTitle {
			return path
		}

		if !visited[node] {
			visited[node] = true
			links, err := fetchPage(node)
			if err != nil {
				log.Printf("Error fetching page for %s: %v", node, err)
				continue
			}
			for _, link := range links {
				newPath := append([]string{}, path...)
				newPath = append(newPath, link)
				queue = append(queue, newPath)
			}
		}
	}

	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <start_article_title> <target_article_title>")
		return
	}

	startArticle := args[0]
	targetArticle := args[1]

	path, numChecked, duration, err := bfs(startArticle, targetArticle)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of articles checked: %d\n", numChecked)
	fmt.Printf("Number of articles traversed: %d\n", len(path))
	fmt.Printf("Route: %s\n", strings.Join(path, " -> "))
	fmt.Printf("Time taken: %v\n", duration)
}
