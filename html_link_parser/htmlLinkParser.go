package htmllinkparser

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// Link is used to store the actual Link in Href and text in Text
type Link struct {
	Href string
	Text string
}

// ExtractLinks extracts the anchor tags Link and text from the url passed
func ExtractLinks(r io.Reader, links []Link) []Link {

	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return links

		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						href := attr.Val
						data := ""

						if z.Next() == html.TextToken {
							data = z.Token().Data
						}
						link := Link{href, data}
						links = append(links, link)
						break
					}
				}
			}

		}

	}
}

// ViewLinks views all the extracted links
func ViewLinks(links []Link) {
	for _, Link := range links {
		fmt.Printf("href : %v\ttext : %v\n", Link.Href, Link.Text)
	}
}
