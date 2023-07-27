package main

import (
	"fmt"

	"pasza.org/sr-challenge/parser"
)

func main() {
	s := parser.Parser()
	fmt.Println(s)
}
