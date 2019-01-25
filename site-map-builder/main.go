package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ayushin78/go_mini_projects/html_link_parser"
)

func main() {
	var rootlink = flag.String("url", "https://example.com", "path of the url")
	flag.Parse()

	visitedLinks := make(map[string]bool)

	crawl(*rootlink, visitedLinks)

	fmt.Println(len(visitedLinks))
}

func crawl(rootlink string, visitedLinks map[string]bool) {
	rootlink = strings.TrimSpace(rootlink)
	rootURL, err := url.Parse(rootlink)
	if err != nil {
		log.Fatal("Invalid url")
		return
	}

	links := []htmllinkparser.Link{
		htmllinkparser.Link{
			Href: rootURL.String(),
			Text: "root",
		},
	}

	for len(links) > 0 {
		currentlink := links[0]
		links = links[1:]
		currentURL, err := getAbsoluteURL(currentlink.Href, rootURL)
		_, visited := visitedLinks[currentURL.String()]
		if err != nil {
			continue // invalid url
		}

		if !visited && haveSameHost(currentURL, rootURL) {
			fmt.Printf("\n start : %v \n", currentURL.String())
			resp, err := http.Get(currentURL.String())
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				log.Fatalf("status code error : %d %s", resp.StatusCode, resp.Status)
			}

			links = htmllinkparser.ExtractLinks(resp.Body, links)
			visitedLinks[currentURL.String()] = true
		}
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
