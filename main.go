package main

import (
	"flag"
	"fmt"
	"io"
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
	file, err := os.Open(*filepath)
	if err != nil {
		fmt.Printf("couldn't open the given file: %s\n", err)
		os.Exit(1)
	}
	fileContent, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("couldn't read the given file: %s\n", err)
		os.Exit(1)
	}
	parser := NewParser(fileContent)
	json, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(json)
	os.Exit(0)
}
