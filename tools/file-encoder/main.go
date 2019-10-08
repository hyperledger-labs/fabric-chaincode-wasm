package main

import (
	"encoding/hex"
	"fmt"
	"github.com/atotto/clipboard"
	"io/ioutil"
)

func main() {

	fmt.Println("Enter file path : ")
	var filepath string
	fmt.Scanf("%s", &filepath)

	fmt.Println("Trying to read file : " + filepath)

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	encodedFile := hex.EncodeToString(file)
	clipboard.WriteAll(encodedFile)
	fmt.Printf("Copied to clipboard. Encoded file: " + encodedFile)
}
