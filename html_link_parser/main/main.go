package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ayushin78/go_mini_projects/html_link_parser"
	"golang.org/x/net/html"
)

func main() {
	content, _ := ioutil.ReadFile("../example.html")
	doc, _ := html.Parse(strings.NewReader(string(content)))

	links := htmllinkparser.ExtractLinks(doc, make([]htmllinkparser.Link, 0))
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(showError))
	mh := myHandler(links)
	mux.Handle("/links/", mh)
	fmt.Println("Server started at port 8080")
	http.ListenAndServe(":8080", mux)
}

func myHandler(links []htmllinkparser.Link) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("../templates/links.html")
		err := t.Execute(w, links)
		if err != nil {
			panic(err)
		}
	}
}

func showError(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404 : Page not found"))
}
