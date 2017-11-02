package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rzaluska/pr2/savegame/compression"
)

func main() {
	filename := flag.String("fn", "", "compressed file name")
	flag.Parse()
	f, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", os.Args[1], err)
		os.Exit(-1)
	}
	defer f.Close()

	stream, err := compression.NewReader(f)
	if err != nil {
		log.Fatalf("Error while reading save file: %s", err)
	}
	_, err = io.Copy(os.Stdout, stream)
	if err != nil {
		log.Printf("Error while reading save file: %s\n", err)
	}
}
