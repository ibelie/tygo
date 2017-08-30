// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

const (
	WireTypeBits = 3
	WireTypeMask = (1 << WireTypeBits) - 1
)

type WireType int

const (
	WireVarint WireType = iota
	WireFixed64
	WireBytes
	WireStart
	WireEnd
	WireFixed32
)

func (i WireType) String() string {
	switch i {
	case WireVarint:
		return "WireVarint"
	case WireFixed64:
		return "WireFixed64"
	case WireBytes:
		return "WireBytes"
	case WireStart:
		return "WireStart"
	case WireEnd:
		return "WireEnd"
	case WireFixed32:
		return "WireFixed32"
	default:
		log.Panicf("[Tygo][WireType] Unexpect enum value: %d", i)
		return "Unknown"
	}
}

func _MAKE_TAG(fieldNum int, wireType WireType) int {
	return (fieldNum << WireTypeBits) | int(wireType)
}

func _MAKE_TAG_STR(fieldNum string, wireType WireType) string {
	return fmt.Sprintf("((%s << %d) | %d)", fieldNum, WireTypeBits, wireType)
}

func _MAKE_CUTOFF(fieldNum int) int {
	max_tag := _MAKE_TAG(fieldNum, WireTypeMask)
	if max_tag <= 0x7F {
		return 0x7F
	} else if max_tag <= 0x3FFF {
		return 0x3FF
	} else {
		return max_tag
	}
}

func _MAKE_CUTOFF_STR(fieldNum string) string {
	return _MAKE_TAG_STR(fieldNum, WireTypeMask)
}

func _TAG_FIELD(tag int) int {
	return tag >> WireTypeBits
}

func _TAG_FIELD_STR(tag string) string {
	return fmt.Sprintf("%s >> %d", tag, WireTypeBits)
}

func _TAG_WIRE(tag int) WireType {
	return WireType(tag & WireTypeMask)
}

func TAG_SIZE(fieldNum int) int {
	return SizeVarint(uint64(fieldNum << WireTypeBits))
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

func SizeBuffer(b []byte) int {
	x := len(b)
	n := x
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}

func SizeSymbol(s string) int {
	x := (len(s)*6 + 7) / 8
	n := x
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

func (p *ProtoBuf) Reset() {
	p.offset = 0
}

func (p *ProtoBuf) Bytes() []byte {
	return p.Buffer[p.offset:]
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

func (p *ProtoBuf) WriteBuf(b []byte) {
	p.WriteVarint(uint64(len(b)))
	p.offset += copy(p.Buffer[p.offset:], b)
}

func (p *ProtoBuf) ReadBuf() ([]byte, error) {
	if l, err := p.ReadVarint(); err != nil {
		return nil, err
	} else if l == 0 {
		return nil, nil
	} else if p.offset+int(l) > len(p.Buffer) {
		return nil, io.EOF
	} else {
		p.offset += int(l)
		return p.Buffer[p.offset-int(l) : p.offset], nil
	}
}

const SymbolDecodeMap = "-ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

var SymbolEncodeMap = [256]byte{
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF /*-*/, 0xFF, 0xFF,
	53 /*0*/, 54 /*1*/, 55 /*2*/, 56 /*3*/, 57 /*4*/, 58 /*5*/, 59 /*6*/, 60, /*7*/
	61 /*8*/, 62 /*9*/, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 1 /*A*/, 2 /*B*/, 3 /*C*/, 4 /*D*/, 5 /*E*/, 6 /*F*/, 7, /*G*/
	8 /*H*/, 9 /*I*/, 10 /*J*/, 11 /*K*/, 12 /*L*/, 13 /*M*/, 14 /*N*/, 15, /*O*/
	16 /*P*/, 17 /*Q*/, 18 /*R*/, 19 /*S*/, 20 /*T*/, 21 /*U*/, 22 /*V*/, 23, /*W*/
	24 /*X*/, 25 /*Y*/, 26 /*Z*/, 0xFF, 0xFF, 0xFF, 0xFF, 63, /*_*/
	0xFF, 27 /*a*/, 28 /*b*/, 29 /*c*/, 30 /*d*/, 31 /*e*/, 32 /*f*/, 33, /*g*/
	34 /*h*/, 35 /*i*/, 36 /*j*/, 37 /*k*/, 38 /*l*/, 39 /*m*/, 40 /*n*/, 41, /*o*/
	42 /*p*/, 43 /*q*/, 44 /*r*/, 45 /*s*/, 46 /*t*/, 47 /*u*/, 48 /*v*/, 49, /*w*/
	50 /*x*/, 51 /*y*/, 52 /*z*/, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
}

func (p *ProtoBuf) EncodeSymbol(s string) {
	si := 0
	src := []byte(s)
	n := (len(src) / 4) * 4
	for si < n {
		// Convert 4x 6bit source bytes into 3 bytes
		val := uint(SymbolEncodeMap[src[si]])<<18 |
			uint(SymbolEncodeMap[src[si+1]])<<12 |
			uint(SymbolEncodeMap[src[si+2]])<<6 |
			uint(SymbolEncodeMap[src[si+3]])

		p.Buffer[p.offset+0] = byte(val >> 16)
		p.Buffer[p.offset+1] = byte(val >> 8)
		p.Buffer[p.offset+2] = byte(val >> 0)
		p.offset += 3
		si += 4
	}

	if len(src) > si {
		var j, val uint
		remain := uint(len(src) - si)
		for j < remain {
			val |= uint(SymbolEncodeMap[src[si+int(j)]]) << (18 - j*6)
			j++
		}
		j = 0
		for j < remain {
			p.Buffer[p.offset] = byte(val >> (16 - j*8))
			p.offset++
			j++
		}
	}
}

func SymbolEncodedLen(data string) int {
	return (len(data)*6 + 7) / 8
}

func DecodeSymbol(src []byte) string {
	dst := make([]byte, len(src)*4/3)
	di, si := 0, 0
	n := (len(src) / 3) * 3
	for si < n {
		// Convert 3x 8bit source bytes into 4 bytes
		val := uint(src[si+0])<<16 | uint(src[si+1])<<8 | uint(src[si+2])

		dst[di+0] = SymbolDecodeMap[val>>18&0x3F]
		dst[di+1] = SymbolDecodeMap[val>>12&0x3F]
		dst[di+2] = SymbolDecodeMap[val>>6&0x3F]
		dst[di+3] = SymbolDecodeMap[val&0x3F]

		si += 3
		di += 4
	}

	switch len(src) - si {
	case 1:
		dst[di+0] = SymbolDecodeMap[uint(src[si])>>2&0x3F]
	case 2:
		val := uint(src[si])<<8 | uint(src[si+1])
		dst[di+0] = SymbolDecodeMap[val>>10&0x3F]
		dst[di+1] = SymbolDecodeMap[val>>4&0x3F]
	}

	if dst[len(dst)-1] == SymbolDecodeMap[0] {
		return string(dst[:len(dst)-1])
	} else {
		return string(dst)
	}
}

func (p *ProtoBuf) WriteSymbol(s string) {
	p.WriteVarint(uint64((len(s)*6 + 7) / 8))
	p.EncodeSymbol(s)
}

func (p *ProtoBuf) ReadSymbol() (string, error) {
	if l, err := p.ReadVarint(); err != nil {
		return "", err
	} else if l == 0 {
		return "", nil
	} else if p.offset+int(l) > len(p.Buffer) {
		return "", io.EOF
	} else {
		p.offset += int(l)
		return DecodeSymbol(p.Buffer[p.offset-int(l) : p.offset]), nil
	}
}

func (p *ProtoBuf) WriteVarint(x uint64) {
	for x >= 0x80 {
		p.Buffer[p.offset] = byte(x) | 0x80
		p.offset++
		x >>= 7
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
				return 0, fmt.Errorf("[Tygo][ProtoBuf] ReadVarint overflow: %v", p.Buffer[p.offset:p.offset+i+1])
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
	p.offset += copy(p.Buffer[p.offset:], x)
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
		x := uint64(p.Buffer[p.offset+0]) |
			uint64(p.Buffer[p.offset+1])<<8 |
			uint64(p.Buffer[p.offset+2])<<16 |
			uint64(p.Buffer[p.offset+3])<<24 |
			uint64(p.Buffer[p.offset+4])<<32 |
			uint64(p.Buffer[p.offset+5])<<40 |
			uint64(p.Buffer[p.offset+6])<<48 |
			uint64(p.Buffer[p.offset+7])<<56
		p.offset += 8
		return x, nil
	}
}

func (p *ProtoBuf) WriteTag(fieldNum int, wireType WireType) {
	p.WriteVarint(uint64(_MAKE_TAG(fieldNum, wireType)))
}

func (p *ProtoBuf) ReadTag(cutoff int) (int, bool, error) {
	if p.offset >= len(p.Buffer) {
		return 0, false, io.EOF
	}
	b1 := int(p.Buffer[p.offset])
	if b1 < 0x80 {
		p.offset++
		return b1, cutoff >= 0x7F || b1 <= cutoff, nil
	}
	if p.offset+1 >= len(p.Buffer) {
		return 0, false, io.EOF
	}
	b2 := int(p.Buffer[p.offset+1])
	if cutoff >= 0x80 && b2 < 0x80 {
		p.offset += 2
		b1 = (b2 << 7) + (b1 - 0x80)
		return b1, cutoff >= 0x3FFF || b1 <= cutoff, nil
	}
	x, err := p.ReadVarint()
	return int(x), int(x) <= cutoff, err
}

func (p *ProtoBuf) ExpectTag(fieldNum int, wireType WireType) bool {
	if p.offset >= len(p.Buffer) {
		return false
	}
	offset := p.offset
	tag := _MAKE_TAG(fieldNum, wireType)
	for tag >= 0x80 {
		if p.Buffer[offset] != byte(tag)|0x80 {
			return false
		}
		tag >>= 7
		offset++
		if offset >= len(p.Buffer) {
			return false
		}
	}
	if p.Buffer[offset] != byte(tag) {
		return false
	}
	p.offset = offset + 1
	return true
}

func (p *ProtoBuf) ExpectBytes(x ...byte) bool {
	if p.offset+len(x) > len(p.Buffer) {
		return false
	} else if bytes.Compare(x, p.Buffer[p.offset:p.offset+len(x)]) != 0 {
		return false
	}
	p.offset += len(x)
	return true
}

func (p *ProtoBuf) ExpectEnd() bool {
	return p.offset >= len(p.Buffer)
}

func (p *ProtoBuf) SkipField(tag int) (err error) {
	switch _TAG_WIRE(tag) {
	case WireVarint:
		_, err = p.ReadVarint()
	case WireFixed64:
		_, err = p.ReadFixed64()
	case WireBytes:
		_, err = p.ReadBuf()
	case WireFixed32:
		_, err = p.ReadFixed32()
	default:
		err = fmt.Errorf("[Tygo][WireType] Unexpect field type to skip: %d", tag)
	}
	return
}

func (t *Enum) WireType() WireType {
	return WireVarint
}

func (t *Object) WireType() WireType {
	return WireBytes
}

func (t UnknownType) WireType() WireType {
	return WireVarint
}

func (t SimpleType) WireType() WireType {
	switch t {
	case SimpleType_INT32:
		fallthrough
	case SimpleType_INT64:
		fallthrough
	case SimpleType_UINT32:
		fallthrough
	case SimpleType_UINT64:
		fallthrough
	case SimpleType_BOOL:
		return WireVarint
	case SimpleType_BYTES:
		return WireBytes
	case SimpleType_STRING:
		return WireBytes
	case SimpleType_SYMBOL:
		return WireBytes
	case SimpleType_FLOAT32:
		return WireFixed32
	case SimpleType_FLOAT64:
		return WireFixed64
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for WireType: %d", t)
		return WireVarint
	}
}

func (t *EnumType) WireType() WireType {
	return WireVarint
}

func (t *InstanceType) WireType() WireType {
	return WireBytes
}

func (t *FixedPointType) WireType() WireType {
	return WireVarint
}

func (t *ListType) WireType() WireType {
	return WireBytes
}

func (t *DictType) WireType() WireType {
	return WireBytes
}

func (t *VariantType) WireType() WireType {
	return WireBytes
}
