package insertPNG

/*
Authors: Tanner Selvig
Desc: Takes a file and returns it as a byte array
Dan Kottmann
*/

import  (
	"bytes"
	"os"
	"bufio"
)

func GetFileBytes(filename string) []byte {
	// Opening a file and creating a reader for it
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

	// Read and write the file out at the same time byte by byte
	var bigAssBuf []byte
	var tmpBuf = make([]byte, 1)
	n := 1
	for n > 0 {
		n, _ = byteReader.Read(tmpBuf)
		bigAssBuf = append(bigAssBuf, tmpBuf[0])
	}
	return bigAssBuf
}