package compression

import (
	"io"

	"github.com/pkg/errors"
)

const dictionarySize = 520 // Retrived in process of reverse engineering the game

type decompressedState int

const (
	dsReadFromStream decompressedState = iota
	dsCopyFromBuffer
)

type Reader struct {
	bs                       *BitStream
	buffer                   CyclicBuffer
	state                    decompressedState
	numBytesToTakeFromBuffer uint
	offsetInBuffer           uint
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		bs:     NewBitStream(r),
		buffer: NewCyclicBuffer(dictionarySize),
		state:  dsReadFromStream,
	}
}

func (d *Reader) Read(b []byte) (num int, err error) {
	c := len(b)
	for ; num < c; num++ {
		bt, err := d.decompressByte()
		if err != nil {
			if err == io.EOF {
				return num, err
			} else {
				return 0, err
			}
		}
		b[num] = bt
	}
	return
}

func (d *Reader) decompressByte() (b byte, e error) {
	switch d.state {
	case dsReadFromStream:

		// read LZSS codeword flag (1 bit)
		isInBuffer, err := d.bs.ReadOneBit()
		if err != nil {
			return 0, err
		}

		if isInBuffer {
			// now we will be reading bytes from cyclic buffer that holds
			// last bytes

			// first of all lets decode position in buffer relative to end and
			// number of bytes to copy
			shift, err := d.bs.ReadBits(3)
			if err != nil {
				return 0, err
			}

			shift++

			baseIndex, err := d.bs.ReadBits(shift)
			if err != nil {
				return 0, err
			}
			baseIndex += ((1 << uint8(shift)) - 2)

			shift = 2
			var numBytes uint32 = 2

			for {
				partial, err := d.bs.ReadBits(shift)
				if err != nil {
					return 0, err
				}
				numBytes += partial
				if partial != (1<<uint8(shift))-1 {
					break
				}
				shift += 1
			}

			d.numBytesToTakeFromBuffer = uint(numBytes)
			d.offsetInBuffer = uint(baseIndex) + 1
			d.state = dsCopyFromBuffer
			return d.decompressByte()
		} else {
			//read byte in normal way
			b, err := d.bs.ReadBits(8)
			if err != nil {
				return 0, err
			}
			// add it to cyclic buffer
			d.buffer.WriteFront(byte(b))
			return byte(b), nil
		}
	case dsCopyFromBuffer:
		b, err := d.buffer.GetFromOffset(d.offsetInBuffer)
		if err != nil {
			return 0, errors.Wrap(err, "Error while reading data from LZSS window buffer")
		}
		//fmt.Printf("buffer:\n%s, offset %d, res: %x\n", d.buffer.String(), d.offsetInBuffer, b)
		d.buffer.WriteFront(b)
		d.numBytesToTakeFromBuffer--
		if d.numBytesToTakeFromBuffer == 0 {
			d.state = dsReadFromStream
		}
		return byte(b), nil
	}
	return
}
