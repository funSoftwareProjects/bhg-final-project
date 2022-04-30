package insertPNG

/*
Author: Tanner Selvig
Desc: Inserts data into a PNG image
Usage: 
Notes: Some code adapted from "Black Hat Go" by Tom Steele, Chris Patten, and 
Dan Kottmann
*/

import  (
	"bytes"
	"encoding/binary"
	"os"
	"bufio"
	"fmt"
)

var secretChunks = []byte("hello");

// Function that inserts data
func ParsePNG(filename string) int {
	// Open file and create a reader for it
	fd, err := os.Open(filename)
	if err != nil { panic(err) }
	stats, err := fd.Stat()
	if err != nil { panic(err) }
	var size = stats.Size()
	buffer := make([]byte, size)
	reader := bufio.NewReader(fd)
	reader.Read(buffer)
	byteReader := bytes.NewReader(buffer)
	_ = byteReader;

	// Need to validate that the file is a PNG
	header := make([]byte, 8)
	err = binary.Read(byteReader, binary.BigEndian, header);
	if err != nil { panic(err) }
	if string(header[1:4]) != "PNG" {
		fmt.Println("Error: Not a PNG file.")
		os.Exit(1);
	}

	// Now we will iterate through each chunk until we know where the EOF is
	chunkType := ""
	var offset int64;
	count := 0
	for chunkType != "IEND" && count < 4 {
		// Read the size of the chunk so we know what to skip 4 bytes
		var chunkSize uint32
		err := binary.Read(byteReader, binary.BigEndian, &chunkSize)
		fmt.Println(chunkSize)
		// Read the type of the chunk 4 bytes
		var chunkType = make([]byte, 4)
		err = binary.Read(byteReader, binary.BigEndian, chunkType)
		_ = err
		fmt.Println(string(chunkType[0:4]))
		// If it is IEND then seek back 8 bytes and take note of where it's at
		if string(chunkType[0:4]) == "IEND" {
			offset, _ = byteReader.Seek(-8, 1)
			break
		}
		// Seek the size of the data
		byteReader.Seek(int64(chunkSize), 1)
		// Seek the CRC
		byteReader.Seek(4, 1)
		count++
	}
	fmt.Println(offset)
	w, err := os.Create("payload.png")
	if err != nil { panic(err) }
	byteReader.Seek(0, 0)
	var tmpBuf = make([]byte, offset)
	var tmpEndBuf = make([]byte, 12)
	byteReader.Read(tmpBuf)
	byteReader.Read(tmpEndBuf)
	w.Write(tmpBuf)
	w.Write(secretChunks)
	w.Write(tmpEndBuf)
	w.Close();
	return 0;
} 