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
	fmt.Printf(".")
	ioutil.ReadAll(os.Stdin)
}

func writeTestFile() {
	f, err := os.OpenFile(dummyFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("What a beautiful string") // 23 bytes
}
