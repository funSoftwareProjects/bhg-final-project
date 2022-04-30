package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"bufio"
)

//usage: go run getFiles.go /
//Windows usage: .\fileSearch.exe C:\Users\

var regexes = []*regexp.Regexp{
	regexp.MustCompile(`\.go$`), //how it decides to pick files should be improved, which regexp option to use?
}

func walkFn(path string, f os.FileInfo, err error) error {
	f, err := os.Create("/tmp/output")
	check(err)
	defer f.Close()

	for _, r := range regexes {
		if r.MatchString(path) {
			//fmt.Printf("%s\n", path)
			n, err := f.WriteString(path)
			check(err)
			fmt.Printf("wrote %s", n)
		}
	}
	f.Sync()
	return nil
}

func main() {
	root := os.Args[1]
	if err := filepath.Walk(root, walkFn); err != nil {
		log.Panicln(err)
	}
}

