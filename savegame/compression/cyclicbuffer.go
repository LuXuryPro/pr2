package compression

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

var CyclicBufferOverflowError = errors.New("Cyclic Buffer Overflow")

type CyclicBuffer struct {
	buffer      []byte
	startIndex  uint
	endIndex    uint
	size        uint
	numElements uint
}

func NewCyclicBuffer(size uint) CyclicBuffer {
	return CyclicBuffer{
		buffer: make([]byte, size),
		size:   size,
	}
}

func (b *CyclicBuffer) WriteFront(val byte) {
	b.buffer[b.startIndex] = val
	b.startIndex += 1
	b.startIndex %= b.size

	if b.startIndex == b.endIndex {
		b.endIndex++
		b.endIndex %= b.size
	}
	b.numElements++
}

func (b *CyclicBuffer) ReadBack() (r byte, err error) {
	if b.startIndex == b.endIndex {
		return 0, CyclicBufferOverflowError
	}
	r = b.buffer[b.endIndex]
	b.endIndex += 1
	b.endIndex %= b.size
	b.numElements--
	return
}

func (b *CyclicBuffer) GetFromOffset(offset uint) (r byte, err error) {
	if offset > b.numElements {
		return 0, errors.Errorf("Attempt to read element outside buffer range: startIndex: %d endIndex: %d numElements: %d in:offset: %d", b.startIndex, b.endIndex, b.numElements, offset)
	}
	indexInBuffer := subtractModulo(b.startIndex, offset, b.size)
	return b.buffer[indexInBuffer], nil
}

func (b *CyclicBuffer) String() string {
	var m []byte
	i := subtractModulo(b.startIndex, 1, b.size)
	for i != subtractModulo(b.endIndex, 1, b.size) {
		es := b.buffer[i]
		m = append(m, es)
		i = subtractModulo(i, 1, b.size)
	}
	return hex.Dump(m)
}

func subtractModulo(a, b, modulus uint) uint {
	first := a % modulus
	second := b % modulus
	return (first - second + modulus) % modulus
}
