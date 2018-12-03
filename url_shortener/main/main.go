package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ayushin78/go_mini_projects/url_shortener"
)

func main() {
	// filepath is used to accept a YAML file as a flag
	var filepath = flag.String("path", "./../example.yaml", "path of the YAML file")
	flag.Parse()
	content, err := ioutil.ReadFile(*filepath)
	if err != nil {
		fmt.Println("Error : File could not be read")
		return
	}

	/*
		This mux will be used as a fallback handler.
		showError function is registered to the newly created ServeMux.
		By doing this, it can be ensured that if the path is not registered
		then, instead of throwing errors, A 404 error message will be displayed.
	*/
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(showError))

	yamlhandler, err := urlshort.YAMLHandler(content, mux)

	log.Println("Listening...")
	http.ListenAndServe(":3000", yamlhandler)
}

func showError(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ERROR 404: Page not found")
}
