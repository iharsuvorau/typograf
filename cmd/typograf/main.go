package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/iharsuvorau/typograf"
)

func main() {
	in := flag.String("i", "", "input text, must be in double quotes")
	flag.Parse()

	if len(*in) > 0 {
		out, err := typograf.Typografy(*in)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(out)
	} else {
		fmt.Println("input text required")
	}
}
