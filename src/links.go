package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
    "container/list"

	"github.com/PuerkitoBio/goquery"
)

type Node struct {
	Name  string
	Path  []string
	Depth int
}

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
		// Exclude links within the references section and external links section
		if s.Closest("div.reflist").Length() == 0 && s.Closest("span#External_links").Length() == 0 {
			href, exists := s.Attr("href")
			if exists && !strings.Contains(href, ":") {
				pageTitle := strings.TrimPrefix(href, "/wiki/")
				links = append(links, pageTitle)
			}
		}
	})

	return links, nil
}

var linkCache = New(1000);

func fetchPageLinksCached(pageName string) ([]string, error) {
    links, ok := linkCache.Get(pageName)
    if ok {
        return links.([]string), nil
    }

    links, err := fetchPageLinks(pageName)
    if err != nil {
        return nil, err
    }

    linkCache.Add(pageName, links)
    return links.([]string), nil
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
					newNode := Node{link, append(node.Path, link), node.Depth + 1}
					stack = append(stack, newNode)
				}
			}
		}
	}

	return nil
}

func IDS(start, goal string, maxDepth int) []string {
	for depth := 0; depth <= maxDepth; depth++ {
		path := DFS(start, goal, depth)
		if path != nil {
			return path
		}
	}
	return nil
}

func main() {
	path := IDS("Basketball", "Joko_Widodo", 3)
	if path == nil {
		fmt.Println("No path found")
	} else {
		fmt.Println("Path found:", path)
	}
}
