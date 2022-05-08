package main

/*
Authors: Tanner Selvig
Desc: Inserts a file into a PNG image
Usage: In the folder "InsertProgram" run go build, then run ./main.
*/

import (
	"bufio"
	"fmt"
	"main/insertPNG"
	"os"
	"strings"
)

func main() {
	// Read in the filename of the PNG from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter PNG name:")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	// Insert the payload into the PNG
	insertPNG.ParsePNG(text)
}
