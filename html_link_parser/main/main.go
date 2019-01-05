package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ayushin78/go_mini_projects/html_link_parser"
	"golang.org/x/net/html"
)

func main() {
	var filepath = flag.String("path", "./../example.html", "path of the HTML file")
	flag.Parse()
	content, err := ioutil.ReadFile(*filepath)
	if err != nil {
		fmt.Println("Error : File could not be read")
		return
	}

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
