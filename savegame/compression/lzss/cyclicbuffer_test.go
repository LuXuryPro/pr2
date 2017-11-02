package lzss

import "testing"

func TestBuffer(t *testing.T) {
	b := NewCyclicBuffer(520)
	if b.size != 520 {
		t.FailNow()
	}
	if b.numElements == 1 {
		t.FailNow()
	}
	b.WriteFront(byte(10))
	v, err := b.GetFromOffset(1)
	if err != nil {
		t.Fatal("Error not expected")
	}
	if v != byte(10) {
		t.Fatal("Excepted: %d found: %d", byte(10), v)
	}
}
