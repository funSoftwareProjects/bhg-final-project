package insertPNG

/*
	Authors: Tanner Selvig
	Desc: Inserts data into a PNG image
	Notes: Some code adapted from "Black Hat Go" by Tom Steele, Chris Patten, and Dan Kottmann
*/

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
)

var secretChunks = []byte("hello")

type ancChunk struct {
	Size uint32
	Type uint32
	Data []byte
	CRC  uint32
}

// Function that inserts data
func ParsePNG(filename string) int {
	// Open file and create a reader for it
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

	// Now we will iterate through each chunk until we know where the EOF is
	chunkType := ""
	var offset int64
	count := 0
	for chunkType != "IEND" && count < 40 {
		// Read the size of the chunk so we know what to skip 4 bytes
		var chunkSize uint32
		err := binary.Read(byteReader, binary.BigEndian, &chunkSize)
		// Read the type of the chunk 4 bytes
		var chunkType = make([]byte, 4)
		err = binary.Read(byteReader, binary.BigEndian, chunkType)
		_ = err
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
	w, err := os.Create("payload.png")
	if err != nil {
		panic(err)
	}
	// Write the secret chunk and ending chunk to the payload png
	byteReader.Seek(0, 0)
	var tmpBuf = make([]byte, offset)
	var tmpEndBuf = make([]byte, 12)
	byteReader.Read(tmpBuf)
	byteReader.Read(tmpEndBuf)
	w.Write(tmpBuf)
	specialChunk := makeChunk(GetFileBytes("insertPNG/main.exe")) // Here is where the file to be disected is specified
	w.Write(specialChunk)
	w.Write(tmpEndBuf)
	w.Close()
	return 0
}

// Creates an ancillary chunk and gives it back as a byte array
func makeChunk(data []byte) []byte {
	var chunk ancChunk
	chunk.Size = uint32(len(data))
	chunk.Data = data
	chunk.Type = binary.BigEndian.Uint32([]byte("rNDm"))
	chunk.CRC = makeCRC(chunk)
	return marshalChunk(chunk)
}

// Converts an ancillary chunk structure to a byte array
func marshalChunk(chunk ancChunk) []byte {
	bytesMSB := new(bytes.Buffer)
	err := binary.Write(bytesMSB, binary.BigEndian, chunk.Size)
	if err != nil {
		panic(err)
	}
	err = binary.Write(bytesMSB, binary.BigEndian, chunk.Type)
	if err != nil {
		panic(err)
	}
	err = binary.Write(bytesMSB, binary.BigEndian, chunk.Data)
	if err != nil {
		panic(err)
	}
	err = binary.Write(bytesMSB, binary.BigEndian, chunk.CRC)
	if err != nil {
		panic(err)
	}
	return bytesMSB.Bytes()
}

// Generates the CRC for the chunk; taken from Black Hat Go
func makeCRC(chunk ancChunk) uint32 {
	bytesMSB := new(bytes.Buffer)
	err := binary.Write(bytesMSB, binary.BigEndian, chunk.Type)
	if err != nil {
		panic(err)
	}
	err = binary.Write(bytesMSB, binary.BigEndian, chunk.Data)
	if err != nil {
		panic(err)
	}
	return crc32.ChecksumIEEE(bytesMSB.Bytes())
}
