package main

import (
	"log"
	"net/http"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
    // Initialize Gin router
    r := gin.Default()

    // Apply CORS middleware
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your frontend URL
    r.Use(cors.New(config))

    // Define a route
    r.POST("pathfinder", prosessPathFinder)
    // Run the server
    r.Run(":8080")
}

type RequestData struct {
	StartArticle string `json:"start_article"`
	TargetArticle string `json:"target_article"`
	SolutionType string `json:"solution_type"`
	Method string `json:"method"`
	Language string `json:"language"`
}

type Article struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Image string `json:"image"`
}

type ResponseData struct {
	Path [][]string `json:"path"`
	Links int `json:"links"`
	Duration int64 `json:"duration"`
	Degrees int `json:"degrees"`
	Language string `json:"language"`
}

func prosessPathFinder(c *gin.Context) {
	var requestData RequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Request data: %v", requestData)
	// Process data using your algorithm
	result := processAlgorithm(requestData)
	log.Printf("Response data: %v", result)
	// Send result back to frontend
	c.JSON(http.StatusOK, result)
}

func processAlgorithm(data RequestData) ResponseData {
	var startURL string = data.StartArticle
	var endURL string = data.TargetArticle
	var solution_type string = data.SolutionType
	var path [][]string
	var duration time.Duration
	var err error 
	var links int
	var degrees int
	var paths []string

	if data.Method == "BFS" {
		bfsInstance := NewBase(startURL, endURL)
		startTime := time.Now()
		paths, err = bfsInstance.Bfs()
		links = bfsInstance.Visit()
		duration = time.Since(startTime)
		degrees = len(paths) - 1
		path = make([][]string, 1)
		path[0] = paths
	} else if data.Method == "IDS" {
		start := time.Now()
		if (solution_type == "one") {
			paths := IDS(startURL, endURL,5)
			path = make([][]string, 1)
			path[0] = paths
			degrees = len(paths) - 1
			} else if (solution_type == "all") {
		path = IDSMany(startURL, endURL, 5, 20)
		degrees = len(path) - 1
		}
		duration = time.Since(start)
		links = articlesChecked
	}

	if err != nil {
		log.Fatalf("Error finding path: %v", err)
	}
	return ResponseData{
		Path: path,
		Links: links,
		Duration: int64(duration.Milliseconds()),
		Degrees: degrees,
		Language: data.Language,
	}
}
