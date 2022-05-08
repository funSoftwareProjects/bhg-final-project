/*
Authors: Tanner Selvig, Will Brant
Desc: Validates if the file contains a payload
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func main() {

	var filename string
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("This program validates whether a file contains a payload")
	fmt.Println("Please enter in a png file for testing (ex: test.png or payload.png)\n")
	fmt.Scanln(&filename)
	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	stats, err := fd.Stat()
	if err != nil {
		panic(err)
	}
	var size = stats.Size()
	buffer := make([]byte, size)
	reader := bufio.NewReader(fd)
	reader.Read(buffer)
	byteReader := bytes.NewReader(buffer)

	// Need to validate that the file is a PNG
	header := make([]byte, 8)
	err = binary.Read(byteReader, binary.BigEndian, header)
	if err != nil {
		panic(err)
	}
	if string(header[1:4]) != "PNG" {
		fmt.Println("Error: Not a PNG file.")
		os.Exit(1)
	}

	chunkType := ""
	malPNG := false

	for chunkType != "IEND" {
		var chunkSize uint32
		err := binary.Read(byteReader, binary.BigEndian, &chunkSize)
		if err != nil {
			panic(err)
		}
		var chunkType = make([]byte, 4)
		err = binary.Read(byteReader, binary.BigEndian, chunkType)
		if err != nil {
			panic(err)
		}
		if string(chunkType[0:4]) == "IEND" {
			break
		}
		firstDataBuf := make([]byte, 2)
		err = binary.Read(byteReader, binary.BigEndian, firstDataBuf)
		byteReader.Seek(int64(chunkSize)-2, 1)
		byteReader.Seek(4, 1)

		if err != nil {
			panic(err)

		} else {
			if !validateBytes(firstDataBuf) {
				malPNG = true
				break
			}

		}
	}
	if malPNG == true {
		fmt.Println("This file may contain a payload. Careful!")
		fmt.Println("----------------------------------------------")
	} else {
		fmt.Println("This png is safe!")
		fmt.Println("----------------------------------------------")
	}
}

func validateBytes(firstBytes []byte) bool {
	// There is a chance that we will get false positives when we use this method but there
	// is really no other good way to identify an executable file
	if firstBytes[0] == 'M' && firstBytes[1] == 'Z' { // Executable files start with MZ in their header
		return false
	}
	return true
}

//for future iterations, pattern matching functions such as this can be used to further validate
func validate(sample string) bool {
	check := true
	for _, j := range sample {
		switch j {
		case '{':
			check = false
		case '}':
			check = false
		case ';':
			check = false
		}
	}
	return check
}
