package main

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