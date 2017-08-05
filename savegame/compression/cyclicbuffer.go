package compression

import "errors"

var CyclicBufferOverflowError = errors.New("Cyclic Buffer Overflow")

type CyclicBuffer struct {
	buffer     []byte
	startIndex int32
	endIndex   int32
	size       int32
}

func NewCyclicBuffer(size int32) CyclicBuffer {
	return CyclicBuffer{
		buffer: make([]byte, size),
		size:   size,
	}
}

func (b *CyclicBuffer) WriteFront(val byte) {
	b.buffer[b.startIndex] = val
	b.startIndex += 1
	b.startIndex %= b.size
}

func (b *CyclicBuffer) ReadBack() (r byte, err error) {
	if b.startIndex == b.endIndex {
		return 0, CyclicBufferOverflowError
	}
	r = b.buffer[b.endIndex]
	b.endIndex += 1
	b.endIndex %= b.size
	return
}
