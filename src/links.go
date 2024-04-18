package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageLinks map[string][]string

    func fetchLinks(pageNames []string) (PageLinks, error) {
        links := make(PageLinks)
        errs := make(chan error)
        linkCh := make(chan struct {
            pageName string
            links    []string
        })

        for _, pageName := range pageNames {
            go func(pageName string) {
                pageLinks, err := fetchPageLinks(pageName)
                if err != nil {
                    errs <- err
                    return
                }
                linkCh <- struct {
                    pageName string
                    links    []string
                }{pageName, pageLinks}
            }(pageName)
        }

        for i := 0; i < len(pageNames); i++ {
            select {
            case err := <-errs:
                return nil, err
            case result := <-linkCh:
                links[result.pageName] = result.links
            }
        }

        return links, nil
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
        doc.Find("#bodyContent a").Each(func(i int, s *goquery.Selection) {
            href, exists := s.Attr("href")
            if exists && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
                pageTitle := strings.TrimPrefix(href, "/wiki/")
                links = append(links, pageTitle)
            }
        })

        return links, nil
    }

