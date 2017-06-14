// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"io"
	"log"
	"strings"
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
	return fmt.Sprintf(`
	if %s != 0 {
		%s += %stygo.SizeVarint(uint64(%s))
	}`, name, size, tagsize, name), map[string]string{TYGO_PATH: ""}
}

func (t *Method) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return "", nil
}

func (t *Object) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	pkgs := map[string]string{TYGO_PATH: ""}
	var fields []string
	if t.HasParent() {
		fields = append(fields, fmt.Sprintf(`
		%s += %s.%s.ByteSize()`, size, name, t.Parent.Name))
	}
	p_num := 0
	var p_name string
	if t.HasParent() {
		if t.Parent.PkgPath == "" {
			p_name, p_num = t.Parent.Object.MaxFieldNum()
		} else {
			p_name = t.Parent.Name
		}
	}
	if p_name != "" {
		fields = append(fields, fmt.Sprintf(`
		preFieldNum := %s.%s.MaxFieldNum()`, name, p_name))
	}
	for i, field := range t.Fields {
		var ts string
		if p_name == "" {
			ts = fmt.Sprintf("%d + ", TAG_SIZE(p_num+i+1))
		} else {
			ts = fmt.Sprintf("tygo.TAG_SIZE(preFieldNum + %d) + ", p_num+i+1)
		}
		field_s, field_p := field.ByteSizeGo(size, fmt.Sprintf("%s.%s", name, field.Name), ts)
		pkgs = update(pkgs, field_p)
		fields = append(fields, fmt.Sprintf(`
		// %s%s
`, field, addIndent(field_s, 1)))
	}
	return fmt.Sprintf(`
	if %s != nil {%s
	}`, name, strings.Join(fields, "")), pkgs
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
	// %s
	if %s != 0 {
		%s += %stygo.SizeVarint(uint64(%s))
	}`, t, name, size, tagsize, name), map[string]string{TYGO_PATH: ""}
	case SimpleType_BYTES:
		fallthrough
	case SimpleType_STRING:
		return fmt.Sprintf(`
	// %s
	if len(%s) > 0 {
		l := len([]byte(%s))
		%s += %stygo.SizeVarint(uint64(l)) + l
	}`, t, name, name, size, tagsize), map[string]string{TYGO_PATH: ""}
	case SimpleType_BOOL:
		return fmt.Sprintf(`
	// %s
	if %s {
		%s += %s1
	}`, t, name, size, tagsize), nil
	case SimpleType_FLOAT32:
		return fmt.Sprintf(`
	// %s
	if %s {
		%s += %s4
	}`, t, name, size, tagsize), nil
	case SimpleType_FLOAT64:
		return fmt.Sprintf(`
	// %s
	if %s {
		%s += %s8
	}`, t, name, size, tagsize), nil
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value: %d", t)
		return "", nil
	}
}

func (t *EnumType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return fmt.Sprintf(`
	// %s
	if %s != 0 {
		%s += %stygo.SizeVarint(uint64(%s))
	}`, t, name, size, tagsize, name), map[string]string{TYGO_PATH: ""}
}

func (t *InstanceType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	var zero string
	if t.IsPtr {
		zero = "nil"
	} else {
		zero = "0"
	}
	return fmt.Sprintf(`
	// %s
	if %s != %s {
		s := %s.ByteSize()
		%s += %stygo.SizeVarint(uint64(s)) + s
	}`, t, name, zero, name, size, tagsize), map[string]string{TYGO_PATH: ""}
}

func (t *FixedPointType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	return fmt.Sprintf(`
	// %s
	if %s != %d {
		%s += %stygo.SizeVarint(uint64((%s - %d) * %d))
	}`, t, name, t.Floor, size, tagsize, name, t.Floor, pow10(t.Precision)), map[string]string{TYGO_PATH: ""}
}

func (t *ListType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	if t.E.IsPrimitive() {
		element_s, element_p := t.E.ByteSizeGo("s", "e", "")
		element_p[TYGO_PATH] = ""
		return fmt.Sprintf(`
	// %s
	if len(%s) > 0 {
		s := 0
		for _, e := range %s {%s else {
				s++
			}
		}
		%s += %stygo.SizeVarint(uint64(s)) + s
	}`, t, name, name, addIndent(element_s, 2), size, tagsize), element_p
	} else {
		element_s, element_p := t.E.ByteSizeGo(size, "e", tagsize)
		element_p[TYGO_PATH] = ""
		return fmt.Sprintf(`
	// %s
	if len(%s) > 0 {
		for _, e := range %s {%s else {
				%s += %s1
			}
		}
	}`, t, name, name, addIndent(element_s, 2), size, tagsize), element_p
	}
}

func (t *DictType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	key_s, key_p := t.K.ByteSizeGo("s", "k", "1 + ")
	value_s, value_p := t.V.ByteSizeGo("s", "v", "1 + ")
	pkgs := map[string]string{TYGO_PATH: ""}
	pkgs = update(pkgs, key_p)
	pkgs = update(pkgs, value_p)
	return fmt.Sprintf(`
	// %s
	if len(%s) > 0 {
		for k, v := range %s {
			s := 0%s%s
			%s += %stygo.SizeVarint(uint64(s)) + s
		}
	}`, t, name, name, addIndent(key_s, 2), addIndent(value_s, 2), size, tagsize), pkgs
}

func (t *VariantType) ByteSizeGo(size string, name string, tagsize string) (string, map[string]string) {
	pkgs := map[string]string{TYGO_PATH: ""}
	var cases []string
	for i, st := range t.Ts {
		type_s, type_p := st.Go()
		if type_s == "nil" {
			continue
		}
		variant_s, variant_p := st.ByteSizeGo("s", "v", fmt.Sprintf("%d + ", TAG_SIZE(i+1)))
		cases = append(cases, fmt.Sprintf(`
		// %s
		case %s:%s else {
				s += %d + 1
			}
			`, st, type_s, addIndent(variant_s, 2), TAG_SIZE(i+1)))
		pkgs = update(pkgs, type_p)
		pkgs = update(pkgs, variant_p)
	}
	return fmt.Sprintf(`
	// %s
	if %s != nil {
		s := 0
		switch v := %s.(type) {%s
		default:
			panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for %s: %%v", v))
		}
		%s += %stygo.SizeVarint(uint64(s)) + s
	}`, t, name, name, strings.Join(cases, ""), t), pkgs
}
