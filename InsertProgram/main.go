package main

/*
Authors: Tanner Selvig
Desc: Inserts a file into a PNG image
Usage: In the folder "InsertProgram" run go build, then run ./main. 
Dan Kottmann
*/

import (
	"bufio"
	"main/insertPNG"
	"fmt"
	"os"
	"strings"
)
func main() {
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter PNG name:")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	insertPNG.ParsePNG(text)
}