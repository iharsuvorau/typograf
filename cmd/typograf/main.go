package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/iharsuvorau/typograf"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if s := fmt.Sprintf("%s", b); len(s) > 0 {
		out, err := typograf.Typogrify(s)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out)
	} else {
		fmt.Println("no input")
	}
}
