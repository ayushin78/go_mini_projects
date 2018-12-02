package urlshort

import (
	"fmt"
	"io/ioutil"
)

func main() {
	content, err := ioutil.ReadFile("./../example.yml")
	if err != nil {
		fmt.Println("Error : File could not be read")
		return
	}

}
