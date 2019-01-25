package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ayushin78/go_mini_projects/html_link_parser"
)

func main() {
	var rootlink = flag.String("url", "https://example.com", "path of the url")
	flag.Parse()

	visitedLinks := make(map[string]bool)

	crawl2(*rootlink, visitedLinks)

	fmt.Println(len(visitedLinks))
}

func crawl2(rootlink string, visitedLinks map[string]bool) {
	rootlink = strings.TrimSpace(rootlink)
	rootURL, err := url.Parse(rootlink)
	if err != nil {
		log.Fatal("Invalid url")
		return
	}

	links := make(chan htmllinkparser.Link, 10)
	links <- htmllinkparser.Link{
		Href: rootURL.String(),
		Text: "root",
	}

	for {
		select {
		case currentlink := <-links:
			currentURL, err := getAbsoluteURL(currentlink.Href, rootURL)
			_, visited := visitedLinks[currentURL.String()]
			if err != nil {
				continue // invalid url
			}

			if !visited && haveSameHost(currentURL, rootURL) {
				fmt.Printf("\n start -> %v \n", currentURL.String())
				go crawlLink(currentURL.String(), links)
				visitedLinks[currentURL.String()] = true
			}

		case <-time.After(2 * time.Second):
			fmt.Println("channel closed")
			return
		}
	}
}

func crawlLink(currentURL string, links chan<- htmllinkparser.Link) {
	resp, err := http.Get(currentURL)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("status code error : %d %s", resp.StatusCode, resp.Status)
		return
	}

	extractedLinks := []htmllinkparser.Link{}

	extractedLinks = htmllinkparser.ExtractLinks(resp.Body, extractedLinks)

	for _, link := range extractedLinks {
		links <- link
	}
}

func getAbsoluteURL(rawurl string, rootURL *url.URL) (*url.URL, error) {
	rawurl = strings.TrimSpace(rawurl) // remove all leading spaces
	u, err := url.Parse(rawurl)        // parses rawuurl
	if err != nil {
		return u, err
	}

	// get absolute path from relative path
	if !u.IsAbs() {
		u = rootURL.ResolveReference(u)
	}
	return u, nil
}

func haveSameHost(currentURL *url.URL, rootURL *url.URL) bool {
	if currentURL.Hostname() == rootURL.Hostname() {
		return true
	}

	return false
}
