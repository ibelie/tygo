// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"io"
	"log"
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

func TAG_SIZE(fieldNum uint32) int {
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

func (p *ProtoBuf) ReadTag(cutoff uint32) (uint32, error) {
	return 0, nil
}

func (p *ProtoBuf) SkipField(fieldNum uint32) error {
	return nil
}

//=============================================================================

func (t *Enum) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *Method) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *Object) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t UnknownType) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t SimpleType) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	switch t {
	case SimpleType_INT32:
		fallthrough
	case SimpleType_INT64:
		fallthrough
	case SimpleType_UINT32:
		fallthrough
	case SimpleType_UINT64:
		return fmt.Sprintf(`
	if %s != 0 {
		size += %s + tygo.SizeVarint(uint64(%s))
	}`, name, tagsize, name), map[string]string{TYGO_PATH: ""}
	case SimpleType_BYTES:
		fallthrough
	case SimpleType_STRING:
		return fmt.Sprintf(`
	if len(%s) > 0 {
		l := len([]byte(%s))
		size += %s + tygo.SizeVarint(l) + l
	}`, name, name, tagsize), map[string]string{TYGO_PATH: ""}
	case SimpleType_BOOL:
		return fmt.Sprintf(`
	if %s {
		size += %s + 1
	}`, name, tagsize), nil
	case SimpleType_FLOAT32:
		return fmt.Sprintf(`
	if %s {
		size += %s + 4
	}`, name, tagsize), nil
	case SimpleType_FLOAT64:
		return fmt.Sprintf(`
	if %s {
		size += %s + 8
	}`, name, tagsize), nil
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value: %d", t)
		return "", nil
	}
}

func (t *EnumType) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *InstanceType) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *FixedPointType) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *ListType) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *DictType) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *VariantType) ByteSizeGo(name string, tagsize string) (string, map[string]string) {
	return "", nil
}
