package savegame

import (
	"encoding/binary"
	"io"
	"log"
)

func ReadHeader(r io.Reader) *Header {
	header := new(Header)
	err := binary.Read(r, binary.LittleEndian, header)
	if err != nil {
		log.Println(err)
	}
	return header
}

type Header struct {
	Version int32 // Savegame version
	Day     int8  // In game day
	Month   int8  // In game month
	Year    int16 // In game year

	Actives        int32 // Total money (coints + property value)
	Difficulty     uint8 // Difficulty level
	Range          uint8 // In game range 0 - 15
	GameMode       uint8
	SpecialAbility uint8
}
