package compression

import "bufio"
import "encoding/binary"

type BitStream struct {
	r               *bufio.Reader
	buffer          uint32
	nextBuffer      uint32
	numBitsInBuffer uint32
}

func NewBitStream(r *bufio.Reader) *BitStream {
	return &BitStream{
		r:               r,
		buffer:          0,
		numBitsInBuffer: 0,
	}
}

func (bs *BitStream) ReadBits(bitNum uint32) (uint32, error) {
	if bitNum > bs.numBitsInBuffer {
		numMissingBits := bitNum - bs.numBitsInBuffer
		err := binary.Read(bs.r, binary.LittleEndian, &bs.nextBuffer)
		if err != nil {
			return 0, err
		}
		res := bs.buffer | (bs.nextBuffer << bs.numBitsInBuffer)
		res &= (1 << bitNum) - 1
		bs.buffer = bs.nextBuffer >> numMissingBits
		bs.numBitsInBuffer = 32 - numMissingBits
		return res, nil
	}
	res := bs.buffer & ((1 << bitNum) - 1)
	bs.buffer >>= bitNum
	bs.numBitsInBuffer -= bitNum
	return res, nil
}

func (ils *BitStream) ReadBit() (bool, error) {
	v, err := ils.ReadBits(1)
	if err != nil {
		return false, err
	}
	return v != 0, nil
}
