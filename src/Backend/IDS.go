package main

import (
	"log"
	"strings"
    "container/list"
    "github.com/gocolly/colly"
	"net/url"
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

// Cache Related
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



func fetchPageLinks(pageName string, lang string) (map[string]struct{}, error) {
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

    err := c.Visit("https://" + lang + ".wikipedia.org/wiki/" + pageName)
    if err != nil {
        return nil, err
    }

    return links, nil
}

var linkCache = New(10000);

func fetchPageLinksCached(pageName string, lang string) (map[string]struct{}, error) {
    links, ok := linkCache.Get(pageName)
    if ok {
        return links.(map[string]struct{}), nil
    }

    links, err := fetchPageLinks(pageName, lang)
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

type job struct {
    pageName string
    lang     string
}

func workers(id int, jobs <-chan job, results chan<- *Node) {
    for j := range jobs {
        links, err := fetchPageLinksCached(j.pageName, j.lang)
        if err != nil {
            log.Printf("Failed to fetch links for page %s: %v", j.pageName, err)
            continue
        }
        children := make([]string, 0, len(links))
        for link := range links {
            children = append(children, link)
        }
        results <- &Node{j.pageName, nil, 0, children}
    }
}


var articlesChecked int

func DFS(start, goal string, maxDepth int, visited map[string]bool, startLang, goalLang string) *Node {
    stack := []*Node{{start, []string{start}, 0, nil}}

    // Create a channel for jobs and a channel for results
    jobs := make(chan job, 100)
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
                jobs <- job{node.Name, startLang}
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


func IDS(startURL, goalURL string, maxDepth int) []string {
    start, startLang, err := extractPageNameAndLang(startURL)
    if err != nil {
        log.Printf("Failed to extract page name from URL %s: %v", startURL, err)
        return nil
    }

    goal, goalLang, err := extractPageNameAndLang(goalURL)
    if err != nil {
        log.Printf("Failed to extract page name from URL %s: %v", goalURL, err)
        return nil
    }

    for depth := 0; depth <= maxDepth; depth++ {
        visited := make(map[string]bool)
        visited[start] = true
        node := DFS(start, goal, depth, visited, startLang, goalLang)
        if node != nil {
            return node.Path
        }
    }
    return nil
}

func extractPageNameAndLang(pageURL string) (string, string, error) {
    u, err := url.Parse(pageURL)
    if err != nil {
        return "", "", err
    }
    lang := strings.Split(u.Host, ".")[0]
    return strings.TrimPrefix(u.Path, "/wiki/"), lang, nil
}



// func main() {
//     start := time.Now()
//     path := IDS("Mike_Tyson", "Joko_Widodo", 3)

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