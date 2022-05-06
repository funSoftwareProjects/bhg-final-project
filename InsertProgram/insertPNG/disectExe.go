package main

import  (
	"bytes"
	"os"
	"bufio"
)

func main() {
	// Opening a file and creating a reader for it
	fd, err := os.Open("../main.exe")
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
	w, err := os.Create("remastered.exe")
	if err != nil { panic(err) }
	var tmpBuf = make([]byte, 1)
	n := 1
	for n > 0 {
		n, _ = byteReader.Read(tmpBuf)
		w.Write(tmpBuf)
	}
	w.Close();
}