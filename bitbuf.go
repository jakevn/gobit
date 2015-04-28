package bitbuf

import "math"

type bitbuf struct {
	buf  []byte
	pos  uint32
	size uint32
}

func NewBitbuf(size uint32) *bitbuf {
	return &bitbuf{make([]byte, size), 0, size * 8}
}

func (b *bitbuf) Size() uint32 {
	return b.size
}

func (b *bitbuf) ByteSize() uint32 {
	return b.size / 8
}

func (b *bitbuf) Pos() uint32 {
	return b.pos
}

func (b *bitbuf) SetPos(pos uint32) {
	if pos > b.size {
		b.pos = b.size
	} else {
		b.pos = pos
	}
}

func (b *bitbuf) Reset() {
	b.pos = 0
	for i, _ := range b.buf {
		b.buf[i] = 0
	}
}

func (b *bitbuf) CanWrite(bits uint32) bool {
	return b.pos+bits <= b.size
}

func (b *bitbuf) CanRead(bits uint32) bool {
	return b.pos+bits <= b.size
}

func (b *bitbuf) WriteBool(value bool) {
	if value {
		b.writeByte(1, 1)
	} else {
		b.writeByte(0, 1)
	}
}

func (b *bitbuf) ReadBool() bool {
	return b.readByte(1) == 1
}

func (b *bitbuf) WriteByte(value byte) {
	b.WriteBytePart(value, 8)
}

func (b *bitbuf) WriteBytePart(value byte, bits uint32) {
	b.writeByte(value, bits)
}

func (b *bitbuf) ReadByte() byte {
	return b.ReadBytePart(8)
}

func (b *bitbuf) ReadBytePart(bits uint32) byte {
	return b.readByte(bits)
}

func (b *bitbuf) WriteUint16(value uint16) {
	b.WriteUint16Part(value, 16)
}

func (b *bitbuf) WriteUint16Part(value uint16, bits uint32) {
	if bits <= 8 {
		b.writeByte(byte(value&0xFF), bits)
	} else {
		b.writeByte(byte(value&0xFF), 8)
		b.writeByte(byte(value>>8), bits-8)
	}
}

func (b *bitbuf) ReadUint16() uint16 {
	return b.ReadUint16Part(16)
}

func (b *bitbuf) ReadUint16Part(bits uint32) uint16 {
	if bits <= 8 {
		return uint16(b.readByte(bits))
	} else {
		return uint16(b.readByte(8) | (b.readByte(bits-8) << 8))
	}
}

func (b *bitbuf) WriteInt16(value int16) {
	b.WriteInt16Part(value, 16)
}

func (b *bitbuf) WriteInt16Part(value int16, bits uint32) {
	b.WriteUint16Part(uint16(value), bits)
}

func (b *bitbuf) ReadInt16() int16 {
	return int16(b.ReadUint16Part(16))
}

func (b *bitbuf) ReadInt16Part(bits uint32) int16 {
	return int16(b.ReadUint16Part(bits))
}

func (b *bitbuf) WriteUint32(value uint32) {
	b.WriteUint32Part(value, 32)
}

func (b *bitbuf) ReadUint32() uint32 {
	return b.ReadUint32Part(32)
}

func (b *bitbuf) WriteInt32Part(value int32, bits uint32) {
	b.WriteUint32Part(uint32(value), 32)
}

func (b *bitbuf) ReadInt32Part(bits uint32) int32 {
	return int32(b.ReadUint32Part(32))
}

func (b *bitbuf) WriteUint32Part(value uint32, bits uint32) {
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

func (b *bitbuf) ReadUint32Part(bits uint32) uint32 {
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

func (b *bitbuf) WriteFloat32(value float32) {
	b.WriteUint32(math.Float32bits(value))
}

func (b *bitbuf) ReadFloat32() float32 {
	return math.Float32frombits(b.ReadUint32())
}

func (b *bitbuf) writeByte(value byte, bits uint32) {
	if bits == 0 {
		return
	}

	p := b.pos >> 3
	used := b.pos & 0x7
	free := uint32(8 - used)
	remain := uint32(free - bits)

	if remain >= 0 {
		mask := byte(0xFF>>free) | (0xFF << (8 - remain))
		b.buf[p] = byte((b.buf[p] & mask) | (value << used))
	} else {
		b.buf[p] = byte(b.buf[p]&(0xFF>>free) | (value << used))
		b.buf[p+1] = byte((b.buf[p+1] & (0xFF << (bits - free))) | (value >> free))
	}

	b.pos += bits
}

func (b *bitbuf) readByte(bits uint32) byte {
	if bits == 0 {
		return 0
	}

	var value byte
	p := b.pos >> 3
	used := b.pos % 8

	if used == 0 && bits == 8 {
		value = b.buf[p]
	} else {
		first := b.buf[p] >> used
		remain := bits - (8 - used)

		if remain < 1 {
			value = byte(first & (0xFF >> (8 - bits)))
		} else {
			second := b.buf[p+1] & (0xFF >> (8 - remain))
			value = byte(first | (second << (bits - remain)))
		}
	}

	b.pos += bits
	return value
}
