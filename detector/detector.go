package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
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
	_ = byteReader

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
	var offset int64
	//count := 0
	//nothing := 0
	for chunkType != "IEND" { //this check does not appear to work. Runs forever
		//fmt.Printf("chunk number %d\n", nothing)
		//nothing++
		newBuf := string(buffer) //buffer holds the current chunk, right?
		bs, err := hex.DecodeString(newBuf)
		//^this decodes the string. Are there header bytes in need to ignore
		if err != nil {
			//panic(err)
			//fmt.Println("not a prob")
		} else {
			//fmt.Println(string(bs))
			valid := validate(string(bs))
			if valid == false {
				fmt.Println("this is a payload")
			} else {
				fmt.Println("This png is safe")
			}

		}
		//fmt.Println(string(bs))

	}
	fmt.Println(offset)

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
