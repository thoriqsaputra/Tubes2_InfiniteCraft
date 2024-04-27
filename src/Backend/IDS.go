package main

import (
	"log"
	"strings"
	"github.com/gocolly/colly"
    "container/list"
	"net/url"
	"sync"
)

// Mendefiniskan struktur Cache
type Cache struct {
    maxEntries int
    ll         *list.List
    cache      map[interface{}]*list.Element
}

// Mendefiniskan struktur entry
type entry struct {
    key   interface{}
    value interface{}
}

// Membuat instansi baru dari Cache
func New(maxEntries int) *Cache {
    return &Cache{
        maxEntries: maxEntries,
        ll:         list.New(),
        cache:      make(map[interface{}]*list.Element),
    }
}

// Menambahkan data ke dalam cache
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

//  Mendapatkan data dari cache
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

// Menghapus data terlama dari cache
func (c *Cache) RemoveOldest() {
    if c.cache == nil {
        return
    }
    ele := c.ll.Back()
    if ele != nil {
        c.removeElement(ele)
    }
}

// Menghapus elemen dari cache
func (c *Cache) removeElement(e *list.Element) {
    c.ll.Remove(e)
    kv := e.Value.(*entry)
    delete(c.cache, kv.key)
}


// Program Utama

// Fungsi extractPageNameAndLang digunakan untuk mengekstrak nama halaman dan bahasa dari URL Wikipedia
func extractPageNameAndLang(pageURL string) (string, string, error) {
    u, err := url.Parse(pageURL)
    if err != nil {
        return "", "", err
    }
    lang := strings.Split(u.Host, ".")[0]
    return strings.TrimPrefix(u.Path, "/wiki/"), lang, nil
}

// Fungsi fetchPageLinks digunakan untuk mengambil semua tautan halaman dari halaman Wikipedia tertentu
func fetchPageLinks(pageName string, lang string) (map[string]struct{}, error) {
    
    // Membuat instance baru dari Collector menggunakan library colly
    c := colly.NewCollector()

    // Membuat map untuk menyimpan tautan yang ditemukan
    links := make(map[string]struct{})

    // Mengatur fungsi callback yang akan dipanggil ketika elemen HTML dengan selector "a[href]" ditemukan
    c.OnHTML("div#mw-content-text p a[href]", func(e *colly.HTMLElement) {
        // Mengambil nilai dari atribut href
        link := e.Attr("href")
        // Jika tautan tidak dimulai dengan "/wiki/", kita abaikan
        if !strings.HasPrefix(link, "/wiki/") {
            return
        }
        // Jika tautan mengandung ":" atau "Category:", kita abaikan
        if strings.Contains(link, ":") || strings.Contains(link, "Category:") {
            return
        }
        // Menghapus prefix "/wiki/" dari tautan dan menyimpannya di map
        pageTitle := strings.TrimPrefix(link, "/wiki/")
        links[pageTitle] = struct{}{}
    })

    // Mengunjungi halaman Wikipedia dan memulai proses scraping
    err := c.Visit("https://" + lang + ".wikipedia.org/wiki/" + pageName)
    
    if err != nil {
        return nil, err
    }

    return links, nil
}

// Membuat cache untuk menyimpan tautan halaman
var linkCache = New(7000);


// Fungsi ini menggunakan cache untuk menghindari pengambilan data yang sama berulang kali
func fetchPageLinksCached(pageName string, lang string) (map[string]struct{}, error) {
    // Mencoba mendapatkan tautan dari cache menggunakan nama halaman sebagai kunci
    links, ok := linkCache.Get(pageName)
    // Jika tautan ditemukan di cache, kembalikan tautan tersebut dan nil sebagai error
    if ok {
        return links.(map[string]struct{}), nil
    }

    // Jika tautan tidak ditemukan di cache, ambil tautan dari halaman Wikipedia menggunakan fungsi fetchPageLinks
    links, err := fetchPageLinks(pageName, lang)
    
    if err != nil {
        return nil, err
    }

    // Menambahkan tautan yang baru diambil ke cache untuk penggunaan di masa mendatang
    linkCache.Add(pageName, links)
    
    return links.(map[string]struct{}), nil
}

// Node merepresentasikan simpul dalam graf pencarian
type Node struct {
    Name     string
    Path     []string
    Depth    int
    Children []string
}

// Mendefinisikan struktur job untuk menyimpan informasi pekerjaan
type job struct {
    pageName string
    lang     string
}

// Fungsi workers digunakan untuk memproses pekerjaan dari channel jobs dan mengirim hasilnya ke channel results
func workers(id int, jobs <-chan job, results chan<- *Node) {
    for j := range jobs {
		// Loop ini akan berjalan selama ada pekerjaan di channel jobs
        links, err := fetchPageLinksCached(j.pageName, j.lang)
        if err != nil {
            log.Printf("Failed to fetch links for page %s: %v", j.pageName, err)
            continue
        }

		// Membuat slice untuk menyimpan tautan yang ditemukan
        children := make([]string, 0, len(links))

		// Menambahkan tautan ke slice children
        for link := range links {
            children = append(children, link)
        }

		// Mengirim Node baru ke channel results
        results <- &Node{j.pageName, nil, 0, children}
    }
}

var articlesChecked int

// Fungsi DFS (Depth-First Search) digunakan untuk mencari jalur dari node start ke node goal
// dengan kedalaman maksimum maxDepth. Fungsi ini menggunakan map visited untuk melacak node yang sudah dikunjungi
func DFS(start, goal string, maxDepth int, visited map[string]bool, startLang, goalLang string) *Node {
    // Membuat stack dan memasukkan node start ke dalam stack
    stack := []*Node{{start, []string{start}, 0, nil}}

    // Membuat channel untuk pekerjaan dan hasil
    jobs := make(chan job, 100)
    results := make(chan *Node, 100)

    // Membuat WaitGroup
    var wg sync.WaitGroup

    // Memulai sejumlah goroutine pekerja
    for w := 1; w <= 10; w++ {
        wg.Add(1)
        go func(w int) {
            defer wg.Done()
            workers(w, jobs, results)
        }(w)
    }

    // Membuat goroutine terpisah untuk menutup channel hasil setelah semua pekerjaan selesai
    go func() {
        wg.Wait()
        close(results)
    }()

    // Melakukan pencarian DFS
    for len(stack) > 0 {
        // Mengambil node terakhir dari stack
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        // Jika node adalah goal, kembalikan node tersebut
        if node.Name == goal {
            return node
        }

        // Jika kedalaman node kurang dari maxDepth
        if node.Depth < maxDepth {
            // Jika node tidak memiliki anak, kirim nama halaman ke channel pekerjaan
            if node.Children == nil {
                jobs <- job{node.Name, startLang}
                // Terima hasil dari channel hasil
                newNode := <-results
                node.Children = newNode.Children
            }

            // Untuk setiap tautan di anak node
            for _, link := range node.Children {
                // Jika tautan belum dikunjungi
                if !visited[link] {
                    visited[link] = true
                    articlesChecked++
                    // Membuat slice baru untuk menyimpan path
                    newPath := make([]string, len(node.Path))
                    // Menyalin path node saat ini ke slice baru
                    copy(newPath, node.Path)
                    // Menambahkan tautan ke path baru
                    newPath = append(newPath, link)
                    // Membuat node baru dan menambahkannya ke stack
                    newNode := &Node{link, newPath, node.Depth + 1, nil}
                    stack = append(stack, newNode)
                }
            }
            // Mengosongkan anak node
            node.Children = node.Children[:0] 
        }
    }

    // Menutup channel pekerjaan setelah semua pekerjaan telah dikirim
    close(jobs)

    // Jika tidak ada jalur yang ditemukan, kembalikan nil
    return nil
}


// Fungsi IDS (Iterative Deepening Search) digunakan untuk mencari jalur dari URL start ke URL goal
// dengan kedalaman maksimum maxDepth. Fungsi ini mengembalikan jalur dalam bentuk slice string
func IDS(startURL, goalURL string, maxDepth int) []string {
    // Mengekstrak nama halaman dan bahasa dari URL start
    start, startLang, err := extractPageNameAndLang(startURL)
    // Jika terjadi error, cetak pesan error dan kembalikan nil
    if err != nil {
        log.Printf("Failed to extract page name from URL %s: %v", startURL, err)
        return nil
    }

    // Mengekstrak nama halaman dan bahasa dari URL goal
    goal, goalLang, err := extractPageNameAndLang(goalURL)
    // Jika terjadi error, cetak pesan error dan kembalikan nil
    if err != nil {
        log.Printf("Failed to extract page name from URL %s: %v", goalURL, err)
        return nil
    }

    // Melakukan pencarian IDS
    for depth := 0; depth <= maxDepth; depth++ {
        // Membuat map untuk melacak halaman yang sudah dikunjungi
        visited := make(map[string]bool)
        visited[start] = true
        // Melakukan pencarian DFS dengan kedalaman tertentu
        node := DFS(start, goal, depth, visited, startLang, goalLang)
        // Jika jalur ditemukan, kembalikan jalur tersebut
        if node != nil {
            return node.Path
        }
    }
    // Jika tidak ada jalur yang ditemukan, kembalikan nil
    return nil
}

// Bonus
func IDSMany(startURL, goalURL string, maxDepth, maxPaths int) [][]string {
    allPaths := [][]string{}
    pathCache := make(map[string]bool)

    for i := 0; i < maxPaths; i++ {
        path := IDS(startURL, goalURL, maxDepth)
        if path != nil {
            // Convert the path to a string so it can be used as a map key
            pathStr := strings.Join(path, "->")
            if _, exists := pathCache[pathStr]; !exists {
                allPaths = append(allPaths, path)
                pathCache[pathStr] = true
            }
        } else {
            break
        }
    }

    return allPaths
}

// func main() {
//     start := time.Now()
//     path := IDS("https://en.wikipedia.org/wiki/Mike_Tyson", "https://en.wikipedia.org/wiki/Joko_Widodo", 3)

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