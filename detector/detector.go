/*
Authors: Tanner Selvig, Will Brant
Desc: Takes a file and returns it as a byte array
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

func main() {

	var filename string
	fmt.Println("Heyo. To test a png for code payloads, provide the name of a png, such as test.png")
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
	//count := 0
	//nothing := 0
	for chunkType != "IEND" { //this check does not appear to work. Runs forever
		var chunkSize uint32
		err := binary.Read(byteReader, binary.BigEndian, &chunkSize)
		if err != nil {panic(err)}
		var chunkType = make([]byte, 4)
		err = binary.Read(byteReader, binary.BigEndian, chunkType)
		if err != nil {panic(err)}
		if string(chunkType[0:4]) == "IEND" {
			break
		}
		firstDataBuf := make([]byte, 2)
		err = binary.Read(byteReader, binary.BigEndian, firstDataBuf)
		byteReader.Seek(int64(chunkSize) - 2, 1)
		byteReader.Seek(4, 1)
		
		//^this decodes the string. Are there header bytes in need to ignore
		if err != nil {
			panic(err)
			//fmt.Println("not a prob")
		} else {
			//fmt.Println(string(bs))
			if(!validateBytes(firstDataBuf)) {
				malPNG = true
				break;
			}
			
		}
		//fmt.Println(string(bs))
		
	}
	if malPNG == true {
		fmt.Println("this is a payload")
	} else {
		fmt.Println("This png is safe")
	}
}

func validateBytes(firstBytes []byte) bool {
	// fmt.Printf("First byte: %c, Second: %c\n", firstBytes[0], firstBytes[1])
	// There is a chance that we will get false positives when we use this method but there
	// is really no other good way to identify an executable file
	if firstBytes[0] == 'M' && firstBytes[1] == 'Z' { // Executable files start with MZ in their header
		// fmt.Printf("Made it inside\n")
		return false
	}
	return true
}

func validate(sample string) bool {
	//decode chunk to string
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

// uInt32ToInt converts a 4 byte big-endian buffer to int.
func uInt32ToInt(buf []byte) (int, error) {
	if len(buf) == 0 || len(buf) > 4 {
		return 0, errors.New("invalid buffer")
	}
	return int(binary.BigEndian.Uint32(buf)), nil
}
