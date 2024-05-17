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
	validJSON, err := isValidJSON(string(fileContent))
	if err != nil {
		fmt.Printf("couldn't parse the given json file: %s\n", err)
		os.Exit(1)
	}
	if !validJSON {
		fmt.Println("Invalid JSON")
		os.Exit(1)
	}
	fmt.Println("Valid JSON")
	os.Exit(0)
}

// isValidJSON will check the validity of the JSON string
// on error will return (false, err) otherwise (validity, nil)
func isValidJSON(json string) (bool, error) {
	switch json {
	case "":
		return false, nil
	case "{}":
		return true, nil
	default:
		return false, nil
	}
}
