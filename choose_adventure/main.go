package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {

	adventure, err := parseJSON()
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/cyoa/", http.HandlerFunc(showIntro))
	ah := adventureHandler(adventure, mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", ah)

}

func parseJSON() (map[string]interface{}, error) {
	var adventure map[string]interface{}
	var filepath = flag.String("path", "./story.json", "path of the JSON file")
	flag.Parse()
	content, err := ioutil.ReadFile(*filepath)
	if err != nil {
		fmt.Println("Error : File could not be read")
		return nil, err
	}

	err = json.Unmarshal(content, &adventure)
	if err != nil {
		return nil, err
	}

	return adventure, nil
}

func adventureHandler(adventure map[string]interface{}, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/story.html")
		path := strings.TrimLeft(r.URL.Path, "/cyoa/")
		arc, ok := adventure[path]
		if ok {
			err := t.Execute(w, arc)
			if err != nil {
				panic(err)
			}
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func showIntro(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/cyoa/intro", http.StatusFound)
}
