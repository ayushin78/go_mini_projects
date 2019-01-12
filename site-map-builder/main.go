package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/ayushin78/go_mini_projects/html_link_parser"
)

func main() {
	var rootlink = flag.String("path", "https://example.com", "path of the url")
	flag.Parse()

	visitedLinks := make(map[htmllinkparser.Link]bool)

	crawl(*rootlink, visitedLinks)

	fmt.Println(len(visitedLinks))
}

func crawl(rootlink string, visitedLinks map[htmllinkparser.Link]bool) {
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

		_, visited := visitedLinks[currentlink]
		currentURL, err := getAbsoluteURL(currentlink.Href, rootURL)
		if err != nil {
			continue // invalid url
		}

		if !visited && haveSameHost(currentURL, rootURL) {
			links = htmllinkparser.ExtractLinks(currentURL.String(), links)
			visitedLinks[currentlink] = true
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
		u = rootURL.ResolveReference(rootURL)
	}
	return u, nil
}

func haveSameHost(currentURL *url.URL, rootURL *url.URL) bool {
	if currentURL.Hostname() == rootURL.Hostname() {
		return true
	}

	return false
}
