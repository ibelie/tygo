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

func MAKE_TAG(fieldNum uint32, wireType WireType) uint32 {
	return (fieldNum << WireTypeBits) | uint32(wireType)
}

func MAX_TAG(fieldNum uint32) uint32 {
	return MAKE_TAG(fieldNum, WireTypeMask)
}

func TAG_SIZE(fieldNum uint32) int {
	return SizeVarint(uint64(fieldNum << WireTypeBits))
}

func TAG_FIELD(tag uint32) uint32 {
	return tag >> WireTypeBits
}

func TAG_WIRE(tag uint32) uint32 {
	return tag & WireTypeMask
}

type ProtoBuf struct {
	offset int
	Buffer []byte
}

func (p *ProtoBuf) WriteUvarint(x uint64) {
	for x >= 0x80 {
		p.Buffer[p.offset] = byte(x) | 0x80
		x >>= 7
		p.offset++
	}
	p.Buffer[p.offset] = byte(x)
}

func (p *ProtoBuf) ReadUvarint() (uint64, error) {
	var x uint64
	var s uint
	for i, b := range p.Buffer[p.offset:] {
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, fmt.Errorf("[Tygo][ProtoBuf] ReadUvarint: %v", p.Buffer[:i+1])
			}
			p.offset += i + 1
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, io.EOF
}

func ReadTag(reader io.Reader, cutoff uint32) (uint32, error) {
	return 0, nil
}

func SkipField(reader io.Reader, fieldNum uint32) error {
	return nil
}

func (t *Enum) BsGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *Method) BsGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *Object) BsGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t SimpleType) BsGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *ObjectType) BsGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *FixedPointType) BsGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *ListType) BsGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *DictType) BsGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *VariantType) BsGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}
