package main

import (
	"flag"
	"fmt"
	"html/template"
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

	foundLinks := crawl(*rootlink)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(showError))
	mh := myHandler(foundLinks)
	mux.Handle("/links/", mh)
	fmt.Println("Server started at port 8080")
	http.ListenAndServe(":8080", mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Whoa, Go is neat!</h1>")
}

func crawl(rootlink string) map[string]string {
	visitedLinks := make(map[string]string)

	rootlink = strings.TrimSpace(rootlink)
	rootURL, err := url.Parse(rootlink)
	if err != nil {
		log.Fatal("Invalid url")
		return visitedLinks
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
			currentlink.Href = currentURL.String()
			_, visited := visitedLinks[currentlink.Href]
			if err != nil {
				continue // invalid url
			}

			if !visited && haveSameHost(currentURL, rootURL) {
				fmt.Printf("\n start -> %v \n", currentURL.String())

				go crawlLink(currentlink.Href, links)
				visitedLinks[currentlink.Href] = currentlink.Text
			}

		case <-time.After(2 * time.Second):
			fmt.Println("channel closed")
			return visitedLinks
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

func myHandler(foundLinks map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/links.html")
		err := t.Execute(w, foundLinks)
		if err != nil {
			panic(err)
		}
	}
}

func showError(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404 : Page not found"))
}
