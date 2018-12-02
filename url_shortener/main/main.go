package urlshort

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ayushin78/go_mini_projects/url_shortener"
)

func main() {
	content, err := ioutil.ReadFile("./../example.yml")
	if err != nil {
		fmt.Println("Error : File could not be read")
		return
	}

	mux := http.NewServeMux()
	yamlhandler, err := urlshort.YAMLHandler(content, mux)

	log.Println("Listening...")
	http.ListenAndServe(":3000", yamlhandler)

}
