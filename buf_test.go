package gobit

import "testing"

func TestNewBufProperSize(t *testing.T) {
	buf := NewBuf(1000)
	if buf.BitSize() != 8000 {
		t.Error("Expected size of 8000, got ", buf.BitSize())
	}
}

func TestNewBufProperPos(t *testing.T) {
	buf := NewBuf(1000)
	if buf.Pos() != 0 {
		t.Error("Expected pos of 0, got ", buf.Pos())
	}
}

func TestWriteReadBool(t *testing.T) {
	buf := NewBuf(1)
	var b bool
	b = true
	buf.WriteBool(b)
	r := buf.ReadBool()
	if r != b {
		t.Error("Expected true, got false")
	}
}

func TestWriteReadByte(t *testing.T) {
	buf := NewBuf(1)
	var b byte
	b = 123
	buf.WriteByte(b)
	r := buf.ReadByte()
	if r != b {
		t.Error("Expected 123, got ", r)
	}
}

func TestWriteReadBytePartial(t *testing.T) {
	buf := NewBuf(1)
	var b byte
	b = 5
	buf.WriteBytePart(b, 4)
	r := buf.ReadBytePart(4)
	if r != b {
		t.Error("Expected 5, got ", r)
	}
}

func TestWriteReadInt16(t *testing.T) {
	buf := NewBuf(2)
	var i int16
	i = 20912
	buf.WriteInt16(i)
	r := buf.ReadInt16()
	if r != i {
		t.Error("Expected 20912, got ", r)
	}
}

func TestWriteReadInt16Partial(t *testing.T) {
	buf := NewBuf(2)
	var i int16
	i = -1024
	buf.WriteInt16Part(i, 11)
	r := buf.ReadInt16Part(11)
	if r != i {
		t.Error("Expected -1024, got ", r)
	}
}

func TestWriteReadUint16(t *testing.T) {
	buf := NewBuf(2)
	var i uint16
	i = 31234
	buf.WriteUint16(i)
	r := buf.ReadUint16()
	if r != i {
		t.Error("Expected 31234, got ", r)
	}
}

func TestWriteReadUint16Partial(t *testing.T) {
	buf := NewBuf(2)
	var i uint16
	i = 942
	buf.WriteUint16Part(i, 12)
	r := buf.ReadUint16Part(12)
	if r != i {
		t.Error("Expected 942, got ", r)
	}
}

func TestWriteReadUint32Part(t *testing.T) {
	buf := NewBuf(3)
	var i uint32
	i = 1234567
	buf.WriteUint32Part(i, 24)
	r := buf.ReadUint32Part(24)
	if r != i {
		t.Error("Expected 1234567, got ", r)
	}
}

func TestWriteReadUint32(t *testing.T) {
	buf := NewBuf(4)
	var i uint32
	i = 31392234
	buf.WriteUint32(i)
	r := buf.ReadUint32()
	if r != i {
		t.Error("Expected 31392234, got ", r)
	}
}

func TestWriteReadUint64Part(t *testing.T) {
	buf := NewBuf(7)
	var i uint64
	i = 1234567890
	buf.WriteUint64Part(i, 56)
	r := buf.ReadUint64Part(56)
	if r != i {
		t.Error("Expected 1234567890, got ", r)
	}
}

func TestWriteReadUint64(t *testing.T) {
	buf := NewBuf(8)
	var i uint64
	i = 12345678901234
	buf.WriteUint64(i)
	r := buf.ReadUint64()
	if r != i {
		t.Error("Expected 12345678901234, got ", r)
	}
}

func TestWriteReadInt64Part(t *testing.T) {
	buf := NewBuf(7)
	var i int64
	i = -1234567
	buf.WriteInt64Part(i, 56)
	r := buf.ReadInt64Part(56)
	if r != i {
		t.Error("Expected -1234567, got ", r)
	}
}

func TestWriteReadFloat32(t *testing.T) {
	buf := NewBuf(4)
	var f float32
	f = 1.234567
	buf.WriteFloat32(f)
	r := buf.ReadFloat32()
	if r != f {
		t.Error("Expected 1.234567, got ", r)
	}
}

func TestWriteReadFloat64(t *testing.T) {
	buf := NewBuf(8)
	var f float64
	f = 1.234567890
	buf.WriteFloat64(f)
	r := buf.ReadFloat64()
	if r != f {
		t.Error("Expcted 1.234567890, got ", f)
	}
}

func TestWriteReadString(t *testing.T) {
	buf := NewBuf(128)
	s := "This is a test string."
	buf.WriteString(s)
	r := buf.ReadString()
	if r != s {
		t.Error("Expected "+s+", got ", r)
	}
}
