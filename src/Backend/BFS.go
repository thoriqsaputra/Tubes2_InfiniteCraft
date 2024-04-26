package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"github.com/PuerkitoBio/goquery"
)

type Base struct {
	startURL      string
	endURL        string
	visitedURL       map[string]bool
	queue         []string
	pageLinks     *PageLinks
	pathToLink    map[string]string
}

type PageLinks struct {
	mu    sync.Mutex
	links map[string][]string
}

func NewBase(startURL, endURL string) *Base {
	return &Base{
		startURL:      startURL,
		endURL:        endURL,
		visitedURL:       make(map[string]bool),
		queue:         []string{startURL},
		pageLinks:     NewPageLinks(),
		pathToLink:    make(map[string]string),
	}
}

func NewPageLinks() *PageLinks {
	return &PageLinks{
		links: make(map[string][]string),
	}
}

func (pl *PageLinks) Add(page, link string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	if _, exists := pl.links[page]; !exists {
		pl.links[page] = []string{}
	}
	pl.links[page] = append(pl.links[page], link)
}

func (pl *PageLinks) Exists(page, link string) bool {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	for _, l := range pl.links[page] {
		if l == link {
			return true
		}
	}
	return false
}

func (pl *PageLinks) GetLinks(page string) []string {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.links[page]
}


func (wr *Base) Bfs() ([]string, error) {
	timeout := 20 * time.Minute // Set timeout 5 menit
	startTime := time.Now()

	for len(wr.queue) > 0 {
		currentPage := wr.queue[0]
		wr.queue = wr.queue[1:]

		// Check if elapsed time exceeds the timeout
		if time.Since(startTime) > timeout {
			return nil, fmt.Errorf("search exceeded time limit of %v", timeout)
		}

		var path []string
		link := currentPage
		for link != wr.startURL {
			path = append([]string{getTitle(link)}, path...)
			link = wr.pathToLink[link]
		}
			
		if currentPage == wr.endURL {
			return wr.buildPath(), nil
		}

		links, err := wr.fetchLinks(currentPage)
		if err != nil {
			return nil, err
		}

		for _, link := range links {
			if !wr.visitedURL[link] {
				wr.visitedURL[link] = true
				wr.queue = append(wr.queue, link)
				wr.pathToLink[link] = currentPage
				if link == wr.endURL {
					return wr.buildPath(), nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no path found from %s to %s", wr.startURL, wr.endURL)
}

func (wr *Base) fetchLinks(pageURL string) ([]string, error) {
	resp, err := wr.getWithTimeout(pageURL, 30*time.Second) // Set timeout 30 detik
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching page %s failed with status: %d", pageURL, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	linkCh := make(chan string)

	go func() {
		doc.Find("p a[href]").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
				link := "https://en.wikipedia.org" + href
				linkCh <- link
			}
		})
		close(linkCh)
	}()

	for link := range linkCh {
		if !wr.pageLinks.Exists(pageURL, link) {
			wr.pageLinks.Add(pageURL, link)
			links = append(links, link)
		}
	}

	return links, nil
}

func (wr *Base) getWithTimeout(url string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	return client.Get(url)
}

func (wr *Base) buildPath() []string {
	var path []string
	currentPage := wr.endURL

	for currentPage != wr.startURL {
		path = append([]string{currentPage}, path...)
		currentPage = wr.pathToLink[currentPage]
	}

	path = append([]string{wr.startURL}, path...)

	return path
}

func getTitle(url string) string {
	title := strings.TrimPrefix(url, "https://en.wikipedia.org/wiki/")
	index := strings.Index(title, "/")
	if index != -1 {
		title = title[:index]
	}
	return title
}

func (wr *Base) Visit() int {
	return len(wr.visitedURL)
}

func ArticleURL(title string) string {
	formattedTitle := strings.ReplaceAll(title, " ", "_")
	return "https://en.wikipedia.org/wiki/" + formattedTitle
}
