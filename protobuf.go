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
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF /*-*/, 0xFF, 0xFF, 0xFF,
	0x35 /*0*/, 0x36 /*1*/, 0x37 /*2*/, 0x38 /*3*/, 0x39 /*4*/, 0x3A /*5*/, 0x3B /*6*/, 0x3C, /*7*/
	0x3D /*8*/, 0x3E /*9*/, 0xFF /***/, 0xFF /***/, 0xFF /***/, 0xFF /***/, 0xFF /***/, 0xFF, /***/
	0xFF /***/, 0x01 /*A*/, 0x02 /*B*/, 0x03 /*C*/, 0x04 /*D*/, 0x05 /*E*/, 0x06 /*F*/, 0x07, /*G*/
	0x08 /*H*/, 0x09 /*I*/, 0x0A /*J*/, 0x0B /*K*/, 0x0C /*L*/, 0x0D /*M*/, 0x0E /*N*/, 0x0F, /*O*/
	0x10 /*P*/, 0x11 /*Q*/, 0x12 /*R*/, 0x13 /*S*/, 0x14 /*T*/, 0x15 /*U*/, 0x16 /*V*/, 0x17, /*W*/
	0x18 /*X*/, 0x19 /*Y*/, 0x1A /*Z*/, 0xFF /***/, 0xFF /***/, 0xFF /***/, 0xFF /***/, 0x3F, /*_*/
	0xFF /***/, 0x1B /*a*/, 0x1C /*b*/, 0x1D /*c*/, 0x1E /*d*/, 0x1F /*e*/, 0x20 /*f*/, 0x21, /*g*/
	0x22 /*h*/, 0x23 /*i*/, 0x24 /*j*/, 0x25 /*k*/, 0x26 /*l*/, 0x27 /*m*/, 0x28 /*n*/, 0x29, /*o*/
	0x2A /*p*/, 0x2B /*q*/, 0x2C /*r*/, 0x2D /*s*/, 0x2E /*t*/, 0x2F /*u*/, 0x30 /*v*/, 0x31, /*w*/
	0x32 /*x*/, 0x33 /*y*/, 0x34 /*z*/, 0xFF /***/, 0xFF /***/, 0xFF /***/, 0xFF /***/, 0xFF, /***/
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
	src := []byte(s)
	n := len(src) / 4 * 4
	for si := 0; si < n; si += 4 {
		// Convert 4x 6bit source bytes into 3 bytes
		val := uint(SymbolEncodeMap[src[si]])<<18 |
			uint(SymbolEncodeMap[src[si+1]])<<12 |
			uint(SymbolEncodeMap[src[si+2]])<<6 |
			uint(SymbolEncodeMap[src[si+3]])

		p.Buffer[p.offset+0] = byte(val >> 16)
		p.Buffer[p.offset+1] = byte(val >> 8)
		p.Buffer[p.offset+2] = byte(val >> 0)
		p.offset += 3
	}

	var val uint
	remain := len(src) - n
	for j := 0; j < remain; j++ {
		val |= uint(SymbolEncodeMap[src[n+j]]) << (18 - uint(j)*6)
	}
	for j := 0; j < remain; j++ {
		p.Buffer[p.offset] = byte(val >> (16 - uint(j)*8))
		p.offset++
	}
}

func SymbolEncodedLen(data string) int {
	return (len(data)*6 + 7) / 8
}

func DecodeSymbol(src []byte) string {
	di := 0
	dst := make([]byte, len(src)*4/3)
	n := len(src) / 3 * 3
	for si := 0; si < n; si += 3 {
		// Convert 3x 8bit source bytes into 4 bytes
		val := uint(src[si+0])<<16 | uint(src[si+1])<<8 | uint(src[si+2])

		dst[di+0] = SymbolDecodeMap[val>>18&0x3F]
		dst[di+1] = SymbolDecodeMap[val>>12&0x3F]
		dst[di+2] = SymbolDecodeMap[val>>6&0x3F]
		dst[di+3] = SymbolDecodeMap[val&0x3F]
		di += 4
	}

	switch len(src) - n {
	case 1:
		dst[di+0] = SymbolDecodeMap[uint(src[n])>>2&0x3F]
	case 2:
		val := uint(src[n])<<8 | uint(src[n+1])
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
