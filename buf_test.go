package gobit

import "testing"

func NewBufProperSize(t *testing.T) {
	buf := NewBuf(1000)
	if buf.Size() != 8000 {
		t.Error("Expected size of 8000, got ", buf.Size)
	}
}

func NewBufProperPos(t *testing.T) {
	buf := NewBuf(1000)
	if buf.Pos() != 0 {
		t.Error("Expected pos of 0, got ", buf.Size)
	}
}

func WriteReadBool(t *testing.T) {
	buf := NewBuf(1)
	var b bool
	b = true
	buf.WriteBool(b)
	r := buf.ReadBool()
	if r != b {
		t.Error("Expected true, got false")
	}
}

func WriteReadByte(t *testing.T) {
	buf := NewBuf(1)
	var b byte
	b = 123
	buf.WriteByte(b)
	r := buf.ReadByte()
	if r != b {
		t.Error("Expected 123, got ", r)
	}
}

func WriteReadBytePartial(t *testing.T) {
	buf := NewBuf(1)
	var b byte
	b = 5
	buf.WriteBytePart(b, 4)
	r := buf.ReadBytePart(4)
	if r != b {
		t.Error("Expected 5, got ", r)
	}
}

func WriteReadInt16(t *testing.T) {
	buf := NewBuf(16)
	var i int16
	i = 20912
	buf.WriteInt16(i)
	r := buf.ReadInt16()
	if r != i {
		t.Error("Expected 20912, got ", r)
	}
}

func WriteReadInt16Partial(t *testing.T) {
	buf := NewBuf(11)
	var i int16
	i = -1024
	buf.WriteInt16Part(i, 11)
	r := buf.ReadInt16Part(11)
	if r != i {
		t.Error("Expected -1024, got ", r)
	}
}

func WriteReadUint16(t *testing.T) {
	buf := NewBuf(16)
	var i uint16
	i = 31234
	buf.WriteUint16(i)
	r := buf.ReadUint16()
	if r != i {
		t.Error("Expected 31234, got ", r)
	}
}

func WriteReadUint16Partial(t *testing.T) {
	buf := NewBuf(12)
	var i uint16
	i = 942
	buf.WriteUint16Part(i, 12)
	r := buf.ReadUint16Part(12)
	if r != i {
		t.Error("Expected 942, got ", r)
	}
}

func WriteReadUint32Part(t *testing.T) {
	buf := NewBuf(24)
	var i uint32
	i = 1234567
	buf.WriteUint32Part(i, 24)
	r := buf.ReadUint32Part(24)
	if r != i {
		t.Error("Expected 1234567, got ", r)
	}
}
