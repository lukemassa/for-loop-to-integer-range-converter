package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/lukemassa/go-integer-range/pkg/ranges"
)

func main() {

	var opts struct {
		Dryrun bool `short:"d" long:"dryrun"`
		// Example of positional arguments
		Args struct {
			Path string
		} `positional-args:"yes" required:"yes"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	filePath := opts.Args.Path

	// Check if the provided path is a file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	if !fileInfo.IsDir() {
		err = ranges.FixFile(filePath, opts.Dryrun)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		err = ranges.FixFile(path, opts.Dryrun)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}

	// Walk through the directory
	err = filepath.Walk(filePath, walkFunc)
	if err != nil {
		log.Fatal(err)
	}

}
