package savegame

import (
	"encoding/binary"
	"io"
	"unicode/utf16"
)

type SaveGame struct {
	Header    *Header
	FirstName string
	LastName  string
}

type Header struct {
	Version int32 // Savegame version
	Day     int8  // In-game day
	Month   int8  // In-game month
	Year    int16 // In-game year

	Actives        int32 // Total money (coints + property value)
	Difficulty     uint8 // Difficulty level
	Rank           uint8 // In game rank 0 - 15
	GameMode       uint8
	SpecialAbility uint8
}

func Read(r io.Reader) (*SaveGame, error) {
	savegame := new(SaveGame)
	h, err := readHeader(r)
	if err != nil {
		return nil, err
	}
	savegame.Header = h
	firstName, err := ReadUTF16ByteString(r)
	if err != nil {
		return nil, err
	}
	savegame.FirstName = firstName
	lastName, err := ReadUTF16ByteString(r)
	if err != nil {
		return nil, err
	}
	savegame.LastName = lastName
	return savegame, nil
}

func readHeader(r io.Reader) (*Header, error) {
	header := new(Header)
	err := binary.Read(r, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func ReadUTF16ByteString(r io.Reader) (string, error) {
	utf16surrogate := make([]byte, 2)
	var acc []uint16
	for {
		_, err := r.Read(utf16surrogate)
		if err != nil {
			return "", err
		}
		val := binary.LittleEndian.Uint16(utf16surrogate)
		if val == 0 {
			break
		}
		acc = append(acc, val)
	}
	return string(utf16.Decode(acc)), nil
}
