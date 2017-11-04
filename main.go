package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rzaluska/pr2/savegame"
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

	stream, _ := compression.NewReader(f)
	s, _ := savegame.Read(stream)
	j, _ := json.MarshalIndent(s, "", "\t")
	fmt.Printf("%s\n", j)
}
