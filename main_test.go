package main

import (
	"io/fs"
	"os"
	"strings"
	"testing"
)

func TestParsing(t *testing.T) {
	testDataDir := os.DirFS("testdata")
	fs.WalkDir(testDataDir, ".", func(path string, d fs.DirEntry, err error) error {
		// skip the directory
		if d.IsDir() {
			return nil
		}
		t.Run(path, func(t *testing.T) {
			fileContent, err := os.ReadFile("testdata/" + path)
			if err != nil {
				t.Errorf("expected to read the file, got: %v", err)
			}
			parser := NewParser(fileContent)
			_, err = parser.Parse()
			if strings.Contains(path, "invalid") && err == nil {
				t.Errorf("expected invalid json, got valid")
			}
			if strings.Contains(path, "valid") && err != nil {
				t.Errorf("expected valid json, got: %v", err)
			}
		})

		return nil
	})
}
