package main

import (
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/ayushin78/go_mini_projects/html_link_parser"
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

	*rootlink = strings.TrimSpace(*rootlink)
	rootURL, _ := url.Parse(*rootlink)
	domain := rootURL.Hostname()

	for len(links) > 0 {
		currentlink := links[0]
		links = links[1:]

		_, visited := visitedLinks[currentlink]
		fmt.Println(currentlink.Href)

		if !visited && haveSameDomain(currentlink.Href, domain) {

			links = htmllinkparser.ExtractLinks(currentlink.Href, links)

			visitedLinks[currentlink] = true
		}
	}

	fmt.Println(len(visitedLinks))
}

func haveSameDomain(link string, domain string) bool {

	link = strings.TrimSpace(link)
	u, err := url.Parse(link)
	if err != nil {
		fmt.Println("error found")
		return false
	}

	if u.Hostname() == domain {
		fmt.Println(link + "---- " + domain)
		return true
	}
	return false
}
