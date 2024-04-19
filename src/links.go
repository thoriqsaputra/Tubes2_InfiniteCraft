package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "sync"

    "github.com/PuerkitoBio/goquery"
)

type Node struct {
    Name  string
    Path  []string
    Depth int
}

var httpClient = &http.Client{}

func fetchPageLinks(pageName string) ([]string, error) {
    res, err := httpClient.Get("https://en.wikipedia.org/wiki/" + pageName)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        return nil, err
    }

    var links []string
    doc.Find("div#bodyContent p a[href^='/wiki/']").Each(func(i int, s *goquery.Selection) {
        href, exists := s.Attr("href")
        if exists && !strings.Contains(href, ":") {
            pageTitle := strings.TrimPrefix(href, "/wiki/")
            links = append(links, pageTitle)
        }
    })

    return links, nil
}

var linkCache = make(map[string][]string)
var cacheMutex = &sync.Mutex{}

func fetchPageLinksCached(pageName string) ([]string, error) {
    cacheMutex.Lock()
    links, ok := linkCache[pageName]
    cacheMutex.Unlock()
    if ok {
        return links, nil
    }

    links, err := fetchPageLinks(pageName)
    if err != nil {
        return nil, err
    }

    cacheMutex.Lock()
    linkCache[pageName] = links
    cacheMutex.Unlock()
    return links, nil
}

func DFS(start, goal string, maxDepth int) []string {
    stack := []Node{{start, []string{start}, 0}}
    visited := map[string]bool{start: true}

    for len(stack) > 0 {
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if node.Name == goal {
            return node.Path
        }

        if node.Depth < maxDepth {
            links, err := fetchPageLinksCached(node.Name)
            if err != nil {
                log.Printf("Failed to fetch links for page %s: %v", node.Name, err)
                continue
            }
            for _, link := range links {
                if !visited[link] {
                    visited[link] = true
                    stack = append(stack, Node{link, append(node.Path, link), node.Depth + 1})
                }
            }
        }
    }

    return nil
}

func main() {
    path := DFS("Mike_Tyson", "Joko_Widodo", 3)
    if path == nil {
        fmt.Println("No path found")
    } else {
        fmt.Println("Path found:", path)
    }
}