package main

import (
	"encoding/binary"
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

	sizeBytes := make([]byte, 4)
	f.Read(sizeBytes)
	_ = int64(binary.LittleEndian.Uint32(sizeBytes))

	stream := compression.NewReader(f)
	_, err = io.Copy(os.Stdout, stream)
	if err != nil {
		log.Println(err)
	}
}
