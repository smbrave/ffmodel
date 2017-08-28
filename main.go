package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		return
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	result, err := DDL2Model(string(content))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	fmt.Printf("%s\n", result)
}
