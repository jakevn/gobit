package gobit

import (
	"math"
)

type Buf struct {
	Bytes []byte
	pos   uint32
	size  uint32
}

func NewBuf(byteSize uint32) *Buf {
	return &Buf{make([]byte, byteSize), 0, byteSize * 8}
}

func (b *Buf) BitSize() uint32 {
	return b.size
}

func (b *Buf) ByteSize() uint32 {
	return b.size / 8
}

func (b *Buf) Pos() uint32 {
	return b.pos
}

func (b *Buf) SetPos(pos uint32) {
	if pos > b.size {
		b.pos = b.size
	} else {
		b.pos = pos
	}
}

func (b *Buf) Reset() {
	b.pos = 0
	for i, _ := range b.Bytes {
		b.Bytes[i] = 0
	}
}

func (b *Buf) CanWrite(bits uint32) bool {
	return b.pos+bits <= b.size
}

func (b *Buf) CanRead(bits uint32) bool {
	return b.pos+bits <= b.size
}

func (b *Buf) WriteBool(value bool) {
	if value {
		b.writeByte(1, 1)
	} else {
		b.writeByte(0, 1)
	}
}

func (b *Buf) ReadBool() bool {
	return b.readByte(1) == 1
}

func (b *Buf) WriteByte(value byte) {
	b.WriteBytePart(value, 8)
}

func (b *Buf) WriteBytePart(value byte, bits uint32) {
	b.writeByte(value, bits)
}

func (b *Buf) ReadByte() byte {
	return b.ReadBytePart(8)
}

func (b *Buf) ReadBytePart(bits uint32) byte {
	return b.readByte(bits)
}

func (b *Buf) WriteUint16(value uint16) {
	b.WriteUint16Part(value, 16)
}

func (b *Buf) WriteUint16Part(value uint16, bits uint32) {
	w := byte(value >> 0)
	x := byte(value >> 8)

	switch (bits + 7) / 8 {
	case 1:
		b.writeByte(w, bits)
		break
	case 2:
		b.writeByte(w, 8)
		b.writeByte(x, bits-8)
		break
	}
}

func (b *Buf) ReadUint16() uint16 {
	return b.ReadUint16Part(16)
}

func (b *Buf) ReadUint16Part(bits uint32) uint16 {
	var w, x int32
	w, x = 0, 0

	switch (bits + 7) / 8 {
	case 1:
		w = int32(b.readByte(bits))
		break
	case 2:
		w = int32(b.readByte(8))
		x = int32(b.readByte(bits - 8))
		break
	}

	return uint16(w | (x << 8))
}

func (b *Buf) WriteInt16(value int16) {
	b.WriteInt16Part(value, 16)
}

func (b *Buf) WriteInt16Part(value int16, bits uint32) {
	b.WriteUint16Part(uint16(value), bits)
}

func (b *Buf) ReadInt16() int16 {
	return int16(b.ReadUint16Part(16))
}

func (b *Buf) ReadInt16Part(bits uint32) int16 {
	return int16(b.ReadUint16Part(bits))
}

func (b *Buf) WriteUint32(value uint32) {
	b.WriteUint32Part(value, 32)
}

func (b *Buf) ReadUint32() uint32 {
	return b.ReadUint32Part(32)
}

func (b *Buf) WriteInt32(value int32) {
	b.WriteInt32Part(value, 32)
}

func (b *Buf) WriteInt32Part(value int32, bits uint32) {
	b.WriteUint32Part(uint32(value), 32)
}

func (b *Buf) ReadInt32() int32 {
	return b.ReadInt32Part(32)
}

func (b *Buf) ReadInt32Part(bits uint32) int32 {
	return int32(b.ReadUint32Part(bits))
}

func (b *Buf) WriteUint32Part(value uint32, bits uint32) {
	w := byte(value >> 0)
	x := byte(value >> 8)
	y := byte(value >> 16)
	z := byte(value >> 24)

	switch (bits + 7) / 8 {
	case 1:
		b.writeByte(w, bits)
		break
	case 2:
		b.writeByte(w, 8)
		b.writeByte(x, bits-8)
		break
	case 3:
		b.writeByte(w, 8)
		b.writeByte(x, 8)
		b.writeByte(y, bits-16)
		break
	case 4:
		b.writeByte(w, 8)
		b.writeByte(x, 8)
		b.writeByte(y, 8)
		b.writeByte(z, bits-24)
		break
	}
}

func (b *Buf) ReadUint32Part(bits uint32) uint32 {
	var w, x, y, z int32
	w, x, y, z = 0, 0, 0, 0

	switch (bits + 7) / 8 {
	case 1:
		w = int32(b.readByte(bits))
		break
	case 2:
		w = int32(b.readByte(8))
		x = int32(b.readByte(bits - 8))
		break
	case 3:
		w = int32(b.readByte(8))
		x = int32(b.readByte(8))
		y = int32(b.readByte(bits - 16))
		break
	case 4:
		w = int32(b.readByte(8))
		x = int32(b.readByte(8))
		y = int32(b.readByte(8))
		z = int32(b.readByte(bits - 24))
		break
	}

	return uint32(w | (x << 8) | (y << 16) | (z << 24))
}

func (b *Buf) WriteFloat32(value float32) {
	b.WriteUint32(math.Float32bits(value))
}

func (b *Buf) ReadFloat32() float32 {
	return math.Float32frombits(b.ReadUint32())
}

func (b *Buf) WriteUint64Part(value uint64, bits uint32) {
	if bits <= 32 {
		b.WriteUint32Part(uint32(value&0xFFFFFFFF), bits)
	} else {
		b.WriteUint32Part(uint32(value), 32)
		b.WriteUint32Part(uint32(value>>32), bits-32)
	}
}

func (b *Buf) ReadUint64Part(bits uint32) uint64 {
	if bits <= 32 {
		return uint64(b.ReadUint32Part(bits))
	}
	a := uint64(b.ReadUint32Part(32))
	x := uint64(b.ReadUint32Part(bits - 32))
	return a | (x << 32)
}

func (b *Buf) WriteUint64(value uint64) {
	b.WriteUint64Part(value, 64)
}

func (b *Buf) ReadUint64() uint64 {
	return b.ReadUint64Part(64)
}

func (b *Buf) WriteInt64Part(value int64, bits uint32) {
	b.WriteUint64Part(uint64(value), bits)
}

func (b *Buf) ReadInt64Part(bits uint32) int64 {
	return int64(b.ReadUint64Part(bits))
}

func (b *Buf) WriteInt64(value int64) {
	b.WriteInt64Part(value, 64)
}

func (b *Buf) ReadInt64() int64 {
	return b.ReadInt64Part(64)
}

func (b *Buf) WriteFloat64(value float64) {
	b.WriteUint64(math.Float64bits(value))
}

func (b *Buf) ReadFloat64() float64 {
	return math.Float64frombits(b.ReadUint64())
}

func (b *Buf) WriteString(value string) {
	b.WriteByteArray([]byte(value))
}

func (b *Buf) ReadString() string {
	return string(b.ReadByteArray())
}

func (b *Buf) WriteByteArray(value []byte) {
	b.WriteInt32(int32(len(value)))
	for _, v := range value {
		b.writeByte(v, 8)
	}
}

func (b *Buf) ReadByteArray() []byte {
	count := b.ReadInt32()
	res := make([]byte, count)
	for i, _ := range res {
		res[i] = b.readByte(8)
	}
	return res
}

func (b *Buf) writeByte(value byte, bits uint32) {
	if bits == 0 {
		return
	}

	p := b.pos >> 3
	used := b.pos & 0x7
	free := uint32(8 - used)
	remain := uint32(free - bits)

	if remain >= 0 {
		mask := byte(0xFF>>free) | (0xFF << (8 - remain))
		b.Bytes[p] = byte((b.Bytes[p] & mask) | (value << used))
	} else {
		b.Bytes[p] = byte(b.Bytes[p]&(0xFF>>free) | (value << used))
		b.Bytes[p+1] = byte((b.Bytes[p+1] & (0xFF << (bits - free))) | (value >> free))
	}

	b.pos += bits
}

func (b *Buf) readByte(bits uint32) byte {
	if bits == 0 {
		return 0
	}

	var value byte
	p := b.pos >> 3
	used := b.pos % 8

	if used == 0 && bits == 8 {
		value = b.Bytes[p]
	} else {
		first := b.Bytes[p] >> used
		remain := bits - (8 - used)
		if remain > 8 {
			remain = 0
		}

		if remain < 1 {
			value = byte(first & (0xFF >> (8 - bits)))
		} else {
			second := b.Bytes[p+1] & (0xFF >> (8 - remain))
			value = byte(first | (second << (bits - remain)))
		}
	}

	b.pos += bits
	return value
}
