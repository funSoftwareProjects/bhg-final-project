package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	//"github.com/spf13/pflag"
	//"github.com/blackhat-go/bhg/ch-13/imgInject/Models"
)

/*
const (
	endChunkType = "IEND"
)

//Header holds the first byte (aka magic byte)
type Header struct {
	Header uint64 //  0:8
}

//Chunk represents a data byte chunk
type Chunk struct {
	Size uint32
	Type uint32
	Data []byte
	CRC  uint32
}

//MetaChunk inherits a Chunk struct
type MetaChunk struct {
	Chk    Chunk
	Offset int64
}
*/
var (
	//flags = pflag.FlagSet{SortFlags: false}
	opts models.CmdLineOpts
	png  pnglib.MetaChunk
)

func main() {
	pull, x := os.Open("Silmarillion_Sticker (2).png")
	//pull, x := os.Open("test.txt")
	if x != nil {
		log.Fatal(x)
	}
	valid, x := preProcess(pull)
	if x != nil {
		log.Fatal(x)
	}
	png.ProcessImage(valid, &opts)
	//if mc.validate(valid) == "correct" {
	//	fmt.Println("This is a PNG. Commence steganography")
	//}
	//mc.validate(valid)
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

	f, errx := os.Create("test.txt")
	if errx != nil {
		fmt.Println(errx)
		return nil, errx
	}
	_, errz := f.WriteString(fmt.Sprintf("%x", bReader))
	if errz != nil {
		fmt.Println(errz)
		f.Close()
		return nil, errz
	}
	//fmt.Printf("%x", bReader)
	return bReader, err

}

func (mc *MetaChunk) validate(b *bytes.Reader) {
	var header Header

	if err := binary.Read(b, binary.BigEndian, &header.Header); err != nil {
		log.Fatal(err)
	}

	bArr := make([]byte, 8)
	binary.BigEndian.PutUint64(bArr, header.Header)

	if string(bArr[1:4]) != "PNG" {
		log.Fatal("Provided file is not a valid PNG format")
	}
	//return "correct"
}
