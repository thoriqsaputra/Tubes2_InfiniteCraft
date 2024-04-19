package main

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

type Node struct {
    Name  string
    Path  []string
    Depth int
}

func fetchPageLinks(pageName string) ([]string, error) {
    res, err := http.Get("https://en.wikipedia.org/wiki/" + pageName)
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

func fetchPageLinksCached(pageName string) ([]string, error) {
    if links, ok := linkCache[pageName]; ok {
        return links, nil
    }

    links, err := fetchPageLinks(pageName)
    if err != nil {
        return nil, err
    }

    linkCache[pageName] = links
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
            links, _ := fetchPageLinksCached(node.Name)
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
    path := DFS("Basketball", "Joko_Widodo", 2)
    if path == nil {
        fmt.Println("No path found")
    } else {
        fmt.Println("Path found:", path)
    }
}