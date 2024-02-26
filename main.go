package main

import (
	"log"

	"github.com/lukemassa/for-loop-to-integer-range-converter/pkg/ranges"
)

func main() {
	err := ranges.FixFile("foo/example.go", false)
	if err != nil {
		log.Fatal(err)
	}

}
