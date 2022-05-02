package main

import (
	//"Users/Will/Desktop/School/COSC/Cyber/stegan/bhg-final-project/models"
	//"Users/Will/Desktop/School/COSC/Cyber/stegan/bhg-final-project/pnglib"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	//"github.com/spf13/pflag"
)

var (
	//flags = pflag.FlagSet{SortFlags: false}
	opts CmdLineOpts
	png  MetaChunk
)

func main() {
	pull, x := os.Open("images/Silmarillion_Sticker (2).png")
	//pull, x := os.Open("test.txt")
	if x != nil {
		log.Fatal(x)
	}
	valid, x := preProcess(pull)
	if x != nil {
		log.Fatal(x)
	}
	png.ProcessImage(valid, &opts)
	if validate(valid) == "correct" {
		fmt.Println("This is a PNG. Commence steganography")
	}
	//mc.validate(valid)
	//validate(valid)
}

//func (mc *MetaChunk) validate(b *bytes.Reader) {
func validate(b *bytes.Reader) string {
	var header Header

	if err := binary.Read(b, binary.BigEndian, &header.Header); err != nil {
		log.Fatal(err)
	}

	bArr := make([]byte, 8)
	binary.BigEndian.PutUint64(bArr, header.Header)

	if string(bArr[1:4]) != "PNG" {
		log.Fatal("Provided file is not a valid PNG format")
		return "lame"
	} else {
		return "correct"
	}
}

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
type CmdLineOpts struct {
	Input   string
	Output  string
	Meta    bool
	Supress bool
	Offset  string
	Inject  bool
	Payload string
	Type    string
	Encode  bool
	Decode  bool
	Key     string
}

func (mc *MetaChunk) ProcessImage(b *bytes.Reader, c *CmdLineOpts) {
	mc.validate(b)
	if (c.Offset != "") && (c.Encode == false && c.Decode == false) {
		var m MetaChunk
		m.Chk.Data = []byte(c.Payload)
		m.Chk.Type = m.strToInt(c.Type)
		m.Chk.Size = m.createChunkSize()
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload: % X\n", m.Chk.Data)
		WriteData(b, c, bmb)
	}
	if (c.Offset != "") && c.Encode {
		var m MetaChunk
		m.Chk.Data = XorEncode([]byte(c.Payload), c.Key)
		m.Chk.Type = m.strToInt(c.Type)
		m.Chk.Size = m.createChunkSize()
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload Encode: % X\n", m.Chk.Data)
		WriteData(b, c, bmb)
	}
	if (c.Offset != "") && c.Decode {
		var m MetaChunk
		offset, _ := strconv.ParseInt(c.Offset, 10, 64)
		b.Seek(offset, 0)
		m.readChunk(b)
		origData := m.Chk.Data
		m.Chk.Data = XorDecode(m.Chk.Data, c.Key)
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", origData)
		fmt.Printf("Payload Decode: % X\n", m.Chk.Data)
		WriteData(b, c, bmb)
	}
	if c.Meta {
		count := 1 //Start at 1 because 0 is reserved for magic byte
		var chunkType string
		for chunkType != endChunkType {
			mc.getOffset(b)
			mc.readChunk(b)
			fmt.Println("---- Chunk # " + strconv.Itoa(count) + " ----")
			fmt.Printf("Chunk Offset: %#02x\n", mc.Offset)
			fmt.Printf("Chunk Length: %s bytes\n", strconv.Itoa(int(mc.Chk.Size)))
			fmt.Printf("Chunk Type: %s\n", mc.chunkTypeToString())
			fmt.Printf("Chunk Importance: %s\n", mc.checkCritType())
			//if c.Suppress == false {
			//	fmt.Printf("Chunk Data: %#x\n", mc.Chk.Data)
			//	} //else if c.Suppress {
			//		fmt.Printf("Chunk Data: %s\n", "Suppressed")
			//	}
			fmt.Printf("Chunk CRC: %x\n", mc.Chk.CRC)
			chunkType = mc.chunkTypeToString()
			count++
		}
	}
}

func (mc *MetaChunk) marshalData() *bytes.Buffer {
	bytesMSB := new(bytes.Buffer)
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Size); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.CRC); err != nil {
		log.Fatal(err)
	}

	return bytesMSB
}

func (mc *MetaChunk) readChunk(b *bytes.Reader) {
	mc.readChunkSize(b)
	mc.readChunkType(b)
	mc.readChunkBytes(b, mc.Chk.Size)
	mc.readChunkCRC(b)
}

func (mc *MetaChunk) readChunkSize(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Size); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkType(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkBytes(b *bytes.Reader, cLen uint32) {
	mc.Chk.Data = make([]byte, cLen)
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkCRC(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.CRC); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) getOffset(b *bytes.Reader) {
	offset, _ := b.Seek(0, 1)
	mc.Offset = offset
}

func (mc *MetaChunk) chunkTypeToString() string {
	h := fmt.Sprintf("%x", mc.Chk.Type)
	decoded, _ := hex.DecodeString(h)
	result := fmt.Sprintf("%s", decoded)
	return result
}

func (mc *MetaChunk) checkCritType() string {
	fChar := string([]rune(mc.chunkTypeToString())[0])
	if fChar == strings.ToUpper(fChar) {
		return "Critical"
	}
	return "Ancillary"
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
	} else {
		fmt.Println("Valid PNG so let us continue!")
	}
}

func (mc *MetaChunk) createChunkSize() uint32 {
	return uint32(len(mc.Chk.Data))
}

func (mc *MetaChunk) createChunkCRC() uint32 {
	bytesMSB := new(bytes.Buffer)
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
	return crc32.ChecksumIEEE(bytesMSB.Bytes())
}

func (mc *MetaChunk) strToInt(s string) uint32 {
	t := []byte(s)
	return binary.BigEndian.Uint32(t)
}

//---------------------------------------------------------------------
//encoders.go
//---------------------------------------------------------------------
func encodeDecode(input []byte, key string) []byte {
	var bArr = make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		bArr[i] += input[i] ^ key[i%len(key)]
	}
	return bArr
}

//XorEncode returns encoded byte array
func XorEncode(decode []byte, key string) []byte {
	return encodeDecode(decode, key)
}

//XorDecode returns decoded byte array
func XorDecode(encode []byte, key string) []byte {
	return encodeDecode(encode, key)
}

//-----------------------------------------------------------------
//reader.go
//-----------------------------------------------------------------
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

//-----------------------------------------------------------------
//writer.go
//-----------------------------------------------------------------

//WriteData writes new data to offset
func WriteData(r *bytes.Reader, c *CmdLineOpts, b []byte) {
	offset, err := strconv.ParseInt(c.Offset, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	w, err := os.OpenFile(c.Output, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("Fatal: Problem writing to the output file!")
	}
	r.Seek(0, 0)

	var buff = make([]byte, offset)
	r.Read(buff)
	w.Write(buff)
	w.Write(b)
	if c.Decode {
		r.Seek(int64(len(b)), 1) // right bitshift to overwrite encode chunk
	}
	_, err = io.Copy(w, r)
	if err == nil {
		fmt.Printf("Success: %s created\n", c.Output)
	}
}
