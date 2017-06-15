// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"io"
)

const (
	WireTypeBits = 3
	WireTypeMask = (1 << WireTypeBits) - 1
)

type WireType uint8

const (
	WireVarint WireType = iota
	WireFixed64
	WireBytes
	WireStartGroup
	WireEndGroup
	WireFixed32
)

func MAKE_TAG(fieldNum uint32, wireType WireType) uint32 {
	return (fieldNum << WireTypeBits) | uint32(wireType)
}

func MAX_TAG(fieldNum uint32) uint32 {
	return MAKE_TAG(fieldNum, WireTypeMask)
}

func TAG_SIZE(fieldNum int) int {
	return SizeVarint(uint64(fieldNum << WireTypeBits))
}

func TAG_FIELD(tag uint32) uint32 {
	return tag >> WireTypeBits
}

func TAG_WIRE(tag uint32) uint32 {
	return tag & WireTypeMask
}

func SizeVarint(x uint64) int {
	n := 0
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}

type ProtoBuf struct {
	offset int
	Buffer []byte
}

func (p *ProtoBuf) Write(b []byte) (n int, err error) {
	n = copy(p.Buffer[p.offset:], b)
	p.offset += n
	if len(b) > n {
		err = fmt.Errorf("[Tygo][ProtoBuf] Write out of range: %d", len(b)-n)
	}
	return
}

func (p *ProtoBuf) Read(b []byte) (n int, err error) {
	n = copy(b, p.Buffer[p.offset:])
	p.offset += n
	if n == 0 && len(b) != 0 {
		err = io.EOF
	}
	return
}

func (p *ProtoBuf) WriteVarint(x uint64) {
	for x >= 0x80 {
		p.Buffer[p.offset] = byte(x) | 0x80
		x >>= 7
		p.offset++
	}
	p.Buffer[p.offset] = byte(x)
	p.offset++
}

func (p *ProtoBuf) ReadVarint() (uint64, error) {
	var x uint64
	var s uint
	for i, b := range p.Buffer[p.offset:] {
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, fmt.Errorf("[Tygo][ProtoBuf] ReadVarint overflow: %v", p.Buffer[:i+1])
			}
			p.offset += i + 1
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, io.EOF
}

func (p *ProtoBuf) WriteBytes(x ...byte) {
	for _, b := range x {
		p.Buffer[p.offset] = b
		p.offset++
	}
}

func (p *ProtoBuf) ReadByte() (byte, error) {
	if p.offset >= len(p.Buffer) {
		return 0, io.EOF
	} else {
		b := p.Buffer[p.offset]
		p.offset++
		return b, nil
	}
}

func (p *ProtoBuf) WriteFixed32(x uint32) {
	p.Buffer[p.offset+0] = byte(x)
	p.Buffer[p.offset+1] = byte(x >> 8)
	p.Buffer[p.offset+2] = byte(x >> 16)
	p.Buffer[p.offset+3] = byte(x >> 24)
	p.offset += 4
}

func (p *ProtoBuf) ReadFixed32() (uint32, error) {
	if p.offset+4 > len(p.Buffer) {
		return 0, io.EOF
	} else {
		x := uint32(p.Buffer[p.offset+0]) |
			uint32(p.Buffer[p.offset+1])<<8 |
			uint32(p.Buffer[p.offset+2])<<16 |
			uint32(p.Buffer[p.offset+3])<<24
		p.offset += 4
		return x, nil
	}
}

func (p *ProtoBuf) WriteFixed64(x uint64) {
	p.Buffer[p.offset+0] = byte(x)
	p.Buffer[p.offset+1] = byte(x >> 8)
	p.Buffer[p.offset+2] = byte(x >> 16)
	p.Buffer[p.offset+3] = byte(x >> 24)
	p.Buffer[p.offset+4] = byte(x >> 32)
	p.Buffer[p.offset+5] = byte(x >> 40)
	p.Buffer[p.offset+6] = byte(x >> 48)
	p.Buffer[p.offset+7] = byte(x >> 56)
	p.offset += 8
}

func (p *ProtoBuf) ReadFixed64() (uint64, error) {
	if p.offset+8 > len(p.Buffer) {
		return 0, io.EOF
	} else {
		x := uint32(p.Buffer[p.offset+0]) |
			uint32(p.Buffer[p.offset+1])<<8 |
			uint32(p.Buffer[p.offset+2])<<16 |
			uint32(p.Buffer[p.offset+3])<<24 |
			uint32(p.Buffer[p.offset+4])<<32 |
			uint32(p.Buffer[p.offset+5])<<40 |
			uint32(p.Buffer[p.offset+6])<<48 |
			uint32(p.Buffer[p.offset+7])<<56
		p.offset += 8
		return x, nil
	}
}

func (p *ProtoBuf) WriteTag(fieldNum uint32, wireType WireType) {
	p.WriteVarint(uint64(MAKE_TAG(fieldNum, wireType)))
}

func (p *ProtoBuf) ReadTag(cutoff uint32) (uint32, error) {
	return 0, nil
}

func (p *ProtoBuf) SkipField(fieldNum uint32) error {
	return nil
}
