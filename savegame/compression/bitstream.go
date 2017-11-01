package compression

import (
	"encoding/binary"
	"io"
)

type BitStream struct {
	r                   io.Reader
	buffer              uint32
	nextBuffer          uint32
	numBitsInBuffer     uint32
	NumProcessedBuffers int
}

func NewBitStream(r io.Reader) *BitStream {
	return &BitStream{
		r: r,
	}
}

func (bs *BitStream) ReadBits(numberOfBits uint32) (uint32, error) {
	if numberOfBits > bs.numBitsInBuffer {
		numMissingBits := numberOfBits - bs.numBitsInBuffer
		err := binary.Read(bs.r, binary.LittleEndian, &bs.nextBuffer)
		if err != nil {
			return 0, err
		}
		bs.NumProcessedBuffers++
		res := bs.buffer | (bs.nextBuffer << bs.numBitsInBuffer)
		res &= (1 << numberOfBits) - 1
		bs.buffer = bs.nextBuffer >> numMissingBits
		bs.numBitsInBuffer = 32 - numMissingBits
		return res, nil
	}
	res := bs.buffer & ((1 << numberOfBits) - 1)
	bs.buffer >>= numberOfBits
	bs.numBitsInBuffer -= numberOfBits
	return res, nil
}

func (bs *BitStream) ReadOneBit() (bool, error) {
	v, err := bs.ReadBits(1)
	if err != nil {
		return false, err
	}
	return v != 0, nil
}
