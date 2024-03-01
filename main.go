package main

import (
	"log"

	"github.com/lukemassa/go-integer-range/ranges"
)

func main() {
	err := ranges.FixFile("foo/example.go", false)
	if err != nil {
		log.Fatal(err)
	}

}
