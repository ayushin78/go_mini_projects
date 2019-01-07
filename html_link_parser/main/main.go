package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	htmllinkparser "github.com/ayushin78/go_mini_projects/html_link_parser"
	"golang.org/x/net/html"
)

func main() {
	var rootLink = flag.String("path", "./example.html", "path of the HTML file")
	flag.Parse()
	fmt.Println(*rootLink)

	resp, err := http.Get(*rootLink)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error : %d %s", resp.StatusCode, resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
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
