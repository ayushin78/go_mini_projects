package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ayushin78/go_mini_projects/html_link_parser"
	"golang.org/x/net/html"
)

func main() {
	var rootlink = flag.String("path", "./../example.html", "path of the HTML file")
	flag.Parse()

	visitedLinks := make(map[htmllinkparser.Link]bool)

	links := []htmllinkparser.Link{
		htmllinkparser.Link{
			Href: *rootlink,
			Text: "root",
		},
	}

	for len(links) > 0 {
		currentlink := links[0]
		links = links[1:]

		_, visited := visitedLinks[currentlink]
		fmt.Println(currentlink.Href)

		if !visited && haveSameDomain(currentlink.Href, "www.taylorswift.com") {

			resp, err := http.Get(currentlink.Href)
			if err != nil {
				log.Fatal(err)
			}

			doc, _ := html.Parse(resp.Body)
			links = htmllinkparser.ExtractLinks(doc, links)

			visitedLinks[currentlink] = true
		}
	}

	fmt.Println(len(visitedLinks))
}

func haveSameDomain(link string, domain string) bool {

	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.Hostname())
	if u.Hostname() == domain {
		return true
	}
	return false
}
