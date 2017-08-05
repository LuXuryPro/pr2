package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gospino/pr2/savegame"
	"github.com/gospino/pr2/savegame/compression"
	"log"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cp", "", "write cpu profile `file`")
var filename = flag.String("fn", "", "compressed file name")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	f, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", os.Args[1], err)
		os.Exit(-1)
	}
	defer f.Close()

	bufferedReader := bufio.NewReader(f)
	sizeBytes := make([]byte, 4)
	bufferedReader.Read(sizeBytes)
	_ = int64(binary.LittleEndian.Uint32(sizeBytes))

	stream := compression.NewReader(bufferedReader)
	h := savegame.ReadHeader(stream)
	j, _ := json.MarshalIndent(h, "", "    ")
	fmt.Println(string(j))

}
