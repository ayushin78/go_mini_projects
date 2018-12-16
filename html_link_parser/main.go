package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type link struct {
	Href string
	Text string
}

func main() {
	content, _ := ioutil.ReadFile("example.html")
	doc, _ := html.Parse(strings.NewReader(string(content)))

	links := extractLinks(doc, []link{})

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(showError))
	mh := myHandler(links)
	mux.Handle("/links/", mh)
	http.ListenAndServe(":8080", mux)
}

func myHandler(links []link) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/links.html")
		err := t.Execute(w, links)
		if err != nil {
			panic(err)
		}
	}
}

func showError(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404 : Page not found"))
}

func extractLinks(n *html.Node, links []link) []link {
	if n.Type == html.ElementNode && n.Data == "a" {
		href := n.Attr[0].Val
		data := n.FirstChild.Data
		l := link{href, data}
		links = append(links, l)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = extractLinks(c, links)
	}
	return links
}

func viewLinks(links []link) {
	for _, link := range links {
		fmt.Printf("href : %v\ttext : %v\n", link.Href, link.Text)
	}
}
