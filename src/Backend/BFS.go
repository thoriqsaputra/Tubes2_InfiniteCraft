package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"github.com/PuerkitoBio/goquery"
	"log"
)

// Base merupakan struktur untuk menjalankan algoritma BFS
type Base struct {
	startURL    string         // URL awal
	endURL      string         // URL akhir
	visitedURL  map[string]bool // Menyimpan URL yang sudah dikunjungi
	queue       []string       // Antrian URL yang akan dikunjungi
	pageLinks   *PageLinks     // Struktur untuk menyimpan tautan halaman
	pathToLink  map[string]string // Menyimpan tautan ke halaman sebelumnya dalam jalur terpendek
}

// PageLinks merupakan struktur untuk menyimpan tautan halaman
type PageLinks struct {
	mu    sync.Mutex      // Mutex untuk mengamankan akses ke data tautan
	links map[string][]string // Peta tautan untuk setiap halaman
}

// NewBase membuat instansi baru dari Base dengan URL awal dan akhir
func NewBase(startURL, endURL string) *Base {
	return &Base{
		startURL:    startURL,
		endURL:      endURL,
		visitedURL:  make(map[string]bool),
		queue:       []string{startURL},
		pageLinks:   NewPageLinks(),
		pathToLink:  make(map[string]string),
	}
}

// NewPageLinks membuat instansi baru dari PageLinks
func NewPageLinks() *PageLinks {
	return &PageLinks{
		links: make(map[string][]string),
	}
}

// Add menambahkan tautan baru ke tautan halaman
func (pl *PageLinks) Add(page, link string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	if _, exists := pl.links[page]; !exists {
		pl.links[page] = []string{}
	}
	pl.links[page] = append(pl.links[page], link)
}

// Exists memeriksa apakah tautan sudah ada di tautan halaman
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

// GetLinks mendapatkan semua tautan dari halaman tertentu
func (pl *PageLinks) GetLinks(page string) []string {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.links[page]
}

// Bfs menjalankan algoritma BFS untuk menemukan jalur terpendek antara URL awal dan akhir
func (wr *Base) Bfs() ([]string, error) {
	timeout := 20 * time.Minute // Set timeout 20 menit
	startTime := time.Now()

	for len(wr.queue) > 0 {
		currentPage := wr.queue[0]
		wr.queue = wr.queue[1:]

		// Periksa apakah waktu yang berlalu melebihi batas waktu
		if time.Since(startTime) > timeout {
			return nil, fmt.Errorf("pencarian melebihi batas waktu %v", timeout)
		}

		var path []string
		link := currentPage
		for link != wr.startURL {
			path = append([]string{getTitle(link)}, path...)
			link = wr.pathToLink[link]
		}
		log.Println("Path: ", path)
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

	return nil, fmt.Errorf("tidak ditemukan jalur dari %s ke %s", wr.startURL, wr.endURL)
}

// fetchLinks mengambil semua tautan dari halaman yang diberikan
func (wr *Base) fetchLinks(pageURL string) ([]string, error) {
	resp, err := wr.getWithTimeout(pageURL, 30*time.Second) // Set timeout 30 detik
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gagal mengambil halaman %s dengan status: %d", pageURL, resp.StatusCode)
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
				// Pengecekan untuk Wikipedia bahasa Indonesia
				if strings.HasPrefix(pageURL, "https://id.wikipedia.org") {
					link := "https://id.wikipedia.org" + href
					linkCh <- link
				} else {
					// Jika bukan Wikipedia bahasa Indonesia, maka anggap sebagai Wikipedia bahasa Inggris
					link := "https://en.wikipedia.org" + href
					linkCh <- link
				}
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

// getWithTimeout mengambil halaman dengan batas waktu tertentu
func (wr *Base) getWithTimeout(url string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	return client.Get(url)
}

// buildPath membangun jalur dari URL akhir ke URL awal
func (wr *Base) buildPath() []string {
	var path []string
	currentPage := wr.endURL

	for currentPage != wr.startURL {
		path = append([]string{getTitle(currentPage)}, path...)
		currentPage = wr.pathToLink[currentPage]
	}

	path = append([]string{getTitle(wr.startURL)}, path...)

	return path
}

// getTitle mendapatkan judul halaman dari URL
func getTitle(url string) string {

	// Pengecekan awalan URL untuk Wikipedia bahasa Inggris
	if strings.HasPrefix(url, "https://en.wikipedia.org/wiki/") {
		title := strings.TrimPrefix(url, "https://en.wikipedia.org/wiki/")
		index := strings.Index(title, "/")
		if index != -1 {
			title = title[:index]
		}
		return title
	}

	// Pengecekan awalan URL untuk Wikipedia bahasa Indonesia
	if strings.HasPrefix(url, "https://id.wikipedia.org/wiki/") {
		title := strings.TrimPrefix(url, "https://id.wikipedia.org/wiki/")
		index := strings.Index(title, "/")
		if index != -1 {
			title = title[:index]
		}
		return title
	}

	return ""
}

// Visit mengembalikan jumlah URL yang telah dikunjungi
func (wr *Base) Visit() int {
	return len(wr.visitedURL)
}

// ArticleURL mengubah judul artikel menjadi URL Wikipedia
func ArticleURL(title string) string {
	formattedTitle := strings.ReplaceAll(title, " ", "_")
	return "https://en.wikipedia.org/wiki/" + formattedTitle
}
