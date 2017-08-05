package compression

import "bufio"

type decompressedState int

const (
	dsReadFromStream decompressedState = iota
	dsCopyFromBuffer
)

type Reader struct {
	bs                           *BitStream
	buffer                       CyclicBuffer // cyclic buffer which holds recent 520 bytes to be used while decompressing incrementally
	state                        decompressedState
	bytesToTakeFromBufferCounter uint32
	offsetInBuffer               uint32
}

func NewReader(wrapped *bufio.Reader) *Reader {
	return &Reader{
		bs:     NewBitStream(wrapped),
		buffer: NewCyclicBuffer(520),
		state:  dsReadFromStream,
		bytesToTakeFromBufferCounter: 0,
		offsetInBuffer:               0,
	}
}

//Read len(b) decompressed bytes from LZSS stream
func (d *Reader) Read(b []byte) (num int, err error) {
	c := len(b)
	for ; num < c; num++ {
		bt, err := d.DecompressByte()
		if err != nil {
			return 0, err
		}
		b[num] = bt
	}
	return
}

func (d *Reader) DecompressByte() (b byte, e error) {
	switch d.state {

	//read from stream
	case dsReadFromStream:

		// read LZSS codeword flag (1 bit)
		isInBuffer, err := d.bs.ReadBit()
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
			var decompressBlockSize uint32 = 2

			for {
				partial, err := d.bs.ReadBits(shift)
				if err != nil {
					return 0, err
				}
				decompressBlockSize += partial
				if partial != (1<<uint8(shift))-1 {
					break
				}
				shift += 1
			}

			d.bytesToTakeFromBufferCounter = decompressBlockSize
			d.offsetInBuffer = baseIndex + 1
			d.state = dsCopyFromBuffer
			return d.DecompressByte()
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

	//read from cyclic buffer
	case dsCopyFromBuffer:
		indexInBuffer := substractModulo(int(d.buffer.startIndex), int(d.offsetInBuffer), 520)

		b = d.buffer.buffer[indexInBuffer]
		d.buffer.WriteFront(b)
		d.bytesToTakeFromBufferCounter -= 1
		if d.bytesToTakeFromBufferCounter == 0 {
			d.state = dsReadFromStream
		}
	}
	return
}

func substractModulo(a int, b int, modulus int) (y int) {
	first := a % modulus
	second := b % modulus
	y = (first - second + modulus) % modulus
	return
}
