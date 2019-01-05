package htmllinkparser

import (
	"fmt"

	"golang.org/x/net/html"
)

// Link is used to store the actual Link in Href and text in Text
type Link struct {
	Href string
	Text string
}

// ExtractLinks extracts the anchor tags Link and text from the html node passed
func ExtractLinks(n *html.Node, links []Link) []Link {

	if n.Type == html.ElementNode && n.Data == "a" {
		href := n.Attr[0].Val
		data := n.FirstChild.Data
		l := Link{href, data}
		links = append(links, l)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = ExtractLinks(c, links)
	}
	return links
}

// ViewLinks views all the extracted links
func ViewLinks(links []Link) {
	for _, Link := range links {
		fmt.Printf("href : %v\ttext : %v\n", Link.Href, Link.Text)
	}
}
