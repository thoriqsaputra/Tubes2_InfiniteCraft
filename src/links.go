package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
    "container/list"

	"github.com/PuerkitoBio/goquery"
)

type Cache struct {
    maxEntries int
    ll         *list.List
    cache      map[interface{}]*list.Element
}

type entry struct {
    key   interface{}
    value interface{}
}

func New(maxEntries int) *Cache {
    return &Cache{
        maxEntries: maxEntries,
        ll:         list.New(),
        cache:      make(map[interface{}]*list.Element),
    }
}

func (c *Cache) Add(key, value interface{}) {
    if c.cache == nil {
        c.cache = make(map[interface{}]*list.Element)
        c.ll = list.New()
    }
    if ee, ok := c.cache[key]; ok {
        c.ll.MoveToFront(ee)
        ee.Value.(*entry).value = value
        return
    }
    ele := c.ll.PushFront(&entry{key, value})
    c.cache[key] = ele
    if c.maxEntries != 0 && c.ll.Len() > c.maxEntries {
        c.RemoveOldest()
    }
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
    if c.cache == nil {
        return
    }
    if ele, hit := c.cache[key]; hit {
        c.ll.MoveToFront(ele)
        return ele.Value.(*entry).value, true
    }
    return
}

func (c *Cache) RemoveOldest() {
    if c.cache == nil {
        return
    }
    ele := c.ll.Back()
    if ele != nil {
        c.removeElement(ele)
    }
}

func (c *Cache) removeElement(e *list.Element) {
    c.ll.Remove(e)
    kv := e.Value.(*entry)
    delete(c.cache, kv.key)
}

var httpClient = &http.Client{}

func fetchPageLinks(pageName string) (map[string]struct{}, error) {
	res, err := httpClient.Get("https://en.wikipedia.org/wiki/" + pageName)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	links := make(map[string]struct{})
	doc.Find("div#bodyContent p a[href^='/wiki/']").Each(func(i int, s *goquery.Selection) {
		// Exclude links within the references section and external links section
		if s.Closest("div.reflist").Length() == 0 && s.Closest("span#External_links").Length() == 0 {
			href, exists := s.Attr("href")
			if exists && !strings.Contains(href, ":") {
				pageTitle := strings.TrimPrefix(href, "/wiki/")
				links[pageTitle] = struct{}{}
			}
		}
	})

	return links, nil
}

var linkCache = New(1000);

func fetchPageLinksCached(pageName string) (map[string]struct{}, error) {
    links, ok := linkCache.Get(pageName)
    if ok {
        return links.(map[string]struct{}), nil
    }

    links, err := fetchPageLinks(pageName)
    if err != nil {
        return nil, err
    }

    linkCache.Add(pageName, links)
    return links.(map[string]struct{}), nil
}

type Node struct {
    Name     string
    Path     []string
    Depth    int
    Children []string
}

func DFS(start, goal string, maxDepth int, visited map[string]bool) *Node {
    stack := []*Node{{start, []string{start}, 0, nil}}

    for len(stack) > 0 {
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if node.Name == goal {
            return node
        }

        if node.Depth < maxDepth {
            if node.Children == nil {
                links, err := fetchPageLinksCached(node.Name)
                if err != nil {
                    log.Printf("Failed to fetch links for page %s: %v", node.Name, err)
                    continue
                }
                node.Children = make([]string, 0, len(links))
                for link := range links {
                    node.Children = append(node.Children, link)
                }
            }

            for _, link := range node.Children {
                if !visited[link] {
                    visited[link] = true
                    newPath := append(node.Path[:], link)
                    newNode := &Node{link, newPath, node.Depth + 1, nil}
                    stack = append(stack, newNode)
                }
            }
            node.Children = node.Children[:0] // Clear the children slice for reuse
        }
    }

    return nil
}

func IDS(start, goal string, maxDepth int) []string {
    for depth := 0; depth <= maxDepth; depth++ {
        visited := make(map[string]bool)
        visited[start] = true
        node := DFS(start, goal, depth, visited)
        if node != nil {
            return node.Path
        }
    }
    return nil
}

func main() {
	path := IDS("Disney_Channel", "Paris", 3)
	if path == nil {
		fmt.Println("No path found")
	} else {
		fmt.Println("Path found:", path)
	}
}
