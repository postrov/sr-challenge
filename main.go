package main

import (
	"fmt"

	"pasza.org/sr-challenge/parser"
)

func main() {
	defer fmt.Println("dupa")
	s := parser.Parser()
	fmt.Println(s)
}
