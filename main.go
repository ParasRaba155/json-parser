package main

import (
	"flag"
	"fmt"
	"os"
)

var filepath = flag.String("filepath", "", "the path of the JSON file")

func main() {
	flag.Parse()
	if filepath == nil || *filepath == "" {
		fmt.Println("must provide the filepath")
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Printf("the filepath: %v\n", *filepath)
}
