package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const dummyFilePath = "dummyfile"

func main() {
	defer os.Remove(dummyFilePath)
	writeTestFile()
	fmt.Printf("\x2A")
	ioutil.ReadAll(os.Stdin)
}

func writeTestFile() {
	f, err := os.OpenFile(dummyFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i := 0; i < 16; i++ { // Write 1024 bytes
		f.WriteString("What a beautiful string which is exactly 64 bytes in length!!!!!")
	}
}
