package main

import (
	"log"
	"strings"
    "container/list"
    "github.com/gocolly/colly"

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

func fetchPageLinks(pageName string) (map[string]struct{}, error) {
    c := colly.NewCollector()

    links := make(map[string]struct{})

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        if !strings.HasPrefix(link, "/wiki/") {
            return
        }
        if strings.Contains(link, ":") || strings.Contains(link, "Category:") {
            return
        }
        pageTitle := strings.TrimPrefix(link, "/wiki/")
        links[pageTitle] = struct{}{}
    })

    err := c.Visit("https://en.wikipedia.org/wiki/" + pageName)
    if err != nil {
        return nil, err
    }

    return links, nil
}

var linkCache = New(5000);

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

func workers(id int, jobs <-chan string, results chan<- *Node) {
    for pageName := range jobs {
        links, err := fetchPageLinksCached(pageName)
        if err != nil {
            log.Printf("Failed to fetch links for page %s: %v", pageName, err)
            continue
        }
        children := make([]string, 0, len(links))
        for link := range links {
            children = append(children, link)
        }
        results <- &Node{pageName, nil, 0, children}
    }
}


var articlesChecked int

func DFS(start, goal string, maxDepth int, visited map[string]bool) *Node {
    stack := []*Node{{start, []string{start}, 0, nil}}

    // Create a channel for jobs and a channel for results
    jobs := make(chan string, 100)
    results := make(chan *Node, 100)

    // Start a fixed number of worker goroutines
    for w := 1; w <= 10; w++ {
        go workers(w, jobs, results)
    }

    for len(stack) > 0 {
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if node.Name == goal {
            return node
        }

        if node.Depth < maxDepth {
            if node.Children == nil {
                // Send the page name to the jobs channel
                jobs <- node.Name
                // Receive the result from the results channel
                newNode := <-results
                node.Children = newNode.Children
            }

            for _, link := range node.Children {
                if !visited[link] {
                    visited[link] = true
                    articlesChecked++
                    newPath := make([]string, len(node.Path)) // Create a new slice to hold the path
                    copy(newPath, node.Path) // Copy the current node's path into the new slice
                    newPath = append(newPath, link) // Append the link to the new path
                    newNode := &Node{link, newPath, node.Depth + 1, nil}
                    stack = append(stack, newNode)
                }
            }
            node.Children = node.Children[:0] 
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

// func main() {
//     start := time.Now()
//     path := IDS("Basketball", "Mark_Cuban", 3)

//     elapsed := time.Since(start)
//     fmt.Println("Time taken:", elapsed.Milliseconds(), "ms")

//     // Jumlah artikel yang diperiksa
//     fmt.Println("Jumlah artikel yang diperiksa:", articlesChecked)

//     // Jumlah artikel yang dilalui
//     if path != nil {
//         fmt.Println("Jumlah artikel yang dilalui:", len(path)-1)
//     } else {
//         fmt.Println("Jumlah artikel yang dilalui: 0")
//     }

//     // Rute
//     if path != nil {
//         fmt.Println("Rute:", strings.Join(path, " -> "))
//     } else {
//         fmt.Println("Rute: No path found")
//     }
// }
