package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const in = "Hello, OTUS!"

func main() {
	out := stringutil.Reverse(in)
	fmt.Println(out)
}
