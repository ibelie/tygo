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

func (t *Enum) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *Method) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *Object) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t UnknownType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t SimpleType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
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
		%s += %stygo.SizeVarint(uint64(%s))
	}`, name, size, tagsize, name), map[string]string{TYGO_PATH: ""}
	case SimpleType_BYTES:
		fallthrough
	case SimpleType_STRING:
		return fmt.Sprintf(`
	if len(%s) > 0 {
		l := len([]byte(%s))
		%s += %stygo.SizeVarint(l) + l
	}`, name, name, size, tagsize), map[string]string{TYGO_PATH: ""}
	case SimpleType_BOOL:
		return fmt.Sprintf(`
	if %s {
		%s += %s1
	}`, name, size, tagsize), nil
	case SimpleType_FLOAT32:
		return fmt.Sprintf(`
	if %s {
		%s += %s4
	}`, name, size, tagsize), nil
	case SimpleType_FLOAT64:
		return fmt.Sprintf(`
	if %s {
		%s += %s8
	}`, name, size, tagsize), nil
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value: %d", t)
		return "", nil
	}
}

func (t *EnumType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return fmt.Sprintf(`
	if %s != 0 {
		%s += %stygo.SizeVarint(uint64(%s))
	}`, name, size, tagsize, name), map[string]string{TYGO_PATH: ""}
}

func (t *InstanceType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *FixedPointType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return fmt.Sprintf(`
	if %s != %d {
		%s += %stygo.SizeVarint(uint64((%s - %d) * %d))
	}`, name, t.Floor, size, tagsize, name, t.Floor, pow10(t.Precision)), map[string]string{TYGO_PATH: ""}
}

func (t *ListType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	if t.E.IsPrimitive() {
		element_s, element_p := t.E.ByteSizeGo("s", "e", "")
		element_p[TYGO_PATH] = ""
		return fmt.Sprintf(`
	if len(%s) > 0 {
		var s uint64
		for _, e := range %s {%s
		}
		%s += %stygo.SizeVarint(s) + s
	}`, name, name, addIndent(element_s, 2), size, tagsize), element_p
	} else {
		element_s, element_p := t.E.ByteSizeGo("size", "e", tagsize)
		element_p[TYGO_PATH] = ""
		return fmt.Sprintf(`
	if len(%s) > 0 {
		for _, e := range %s {%s else {
				%s += %s1
			}
		}
	}`, name, name, addIndent(element_s, 2), size, tagsize), element_p
	}
}

func (t *DictType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	key_s, key_p := t.K.ByteSizeGo("s", "k", "1 + ")
	value_s, value_p := t.V.ByteSizeGo("s", "v", "1 + ")
	inner_p := map[string]string{TYGO_PATH: ""}
	inner_p = update(inner_p, key_p)
	inner_p = update(inner_p, value_p)
	return fmt.Sprintf(`
	if len(%s) > 0 {
		for k, v := range %s {
			var s uint64%s%s
			%s += %stygo.SizeVarint(s) + s
		}
	}`, name, name, addIndent(key_s, 2), addIndent(value_s, 2), size, tagsize), inner_p
}

func (t *VariantType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return "", nil
}
