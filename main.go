package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

//here is a comment
func main() {
	pull, x := os.Open("getFiles.go")
	if x != nil {
		log.Fatal(x)
	}
	preProcess(pull)
}

func preProcess(dat *os.File) (*bytes.Reader, error) {
	stats, err := dat.Stat()
	if err != nil {
		return nil, err
	}
	var size = stats.Size()
	b := make([]byte, size)
	bufR := bufio.NewReader(dat)
	_, err = bufR.Read(b)
	bReader := bytes.NewReader(b)

	return bReader, err
}
