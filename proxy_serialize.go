// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func writeTag(preFieldNum string, fieldNum int, wireType WireType, indent int) string {
	if preFieldNum == "" {
		if fieldNum <= 0 {
			return ""
		} else {
			tagbuf := &ProtoBuf{Buffer: make([]byte, TAG_SIZE(fieldNum))}
			tagbuf.WriteTag(fieldNum, wireType)
			var tagbytes []string
			for _, i := range tagbuf.Buffer {
				tagbytes = append(tagbytes, strconv.Itoa(int(i)))
			}
			return fmt.Sprintf(`
	%soutput.WriteBytes(%s) // tag: %d MAKE_TAG(%d, %s=%d)`, strings.Repeat("\t", indent),
				strings.Join(tagbytes, ", "), _MAKE_TAG(fieldNum, wireType), fieldNum, wireType, wireType)
		}
	} else {
		return fmt.Sprintf(`
	%soutput.WriteTag(%s + %d, %d)`, strings.Repeat("\t", indent), preFieldNum, fieldNum, wireType)
	}
}

func (t *Enum) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return fmt.Sprintf(`
	if %s != 0 {
		output.WriteVarint(uint64(%s))
	}`, name, name), nil
}

func (t *Method) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return "", nil
}

func (t *Object) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	var pkgs map[string]string
	var fields []string
	if t.HasParent() {
		fields = append(fields, fmt.Sprintf(`
		%s.%s.Serialize(output)`, name, t.Parent.Name))
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
		p_name = "preFieldNum"
	}

	for i, field := range t.Fields {
		field_s, field_p := field.SerializeGo(size, fmt.Sprintf("%s.%s", name, field.Name), p_name, p_num+i+1, true)
		pkgs = update(pkgs, field_p)
		fields = append(fields, fmt.Sprintf(`
		// property: %s.%s%s
`, name, field.Name, addIndent(field_s, 1)))
	}

	return fmt.Sprintf(`
	if %s != nil {%s
	}`, name, strings.Join(fields, "")), pkgs
}

func (t UnknownType) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return "", nil
}

func (t SimpleType) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	switch t {
	case SimpleType_INT32:
		fallthrough
	case SimpleType_INT64:
		fallthrough
	case SimpleType_UINT32:
		fallthrough
	case SimpleType_UINT64:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if %s != 0 {%s
		output.WriteVarint(uint64(%s))
	}`, t, name, writeTag(preFieldNum, fieldNum, WireVarint, 1), name), nil
		} else {
			return fmt.Sprintf(`
	// type: %s%s
	output.WriteVarint(uint64(%s))`, t, writeTag(preFieldNum, fieldNum, WireVarint, 0), name), nil
		}
	case SimpleType_BYTES:
		fallthrough
	case SimpleType_STRING:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {%s
		output.WriteBuf([]byte(%s))
	}`, t, name, writeTag(preFieldNum, fieldNum, WireBytes, 1), name), nil
		} else {
			return fmt.Sprintf(`
	// type: %s
	{%s
		output.WriteBuf([]byte(%s))
	}`, t, writeTag(preFieldNum, fieldNum, WireBytes, 1), name), nil
		}
	case SimpleType_BOOL:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if %s {%s
		output.WriteBytes(1)
	}`, t, name, writeTag(preFieldNum, fieldNum, WireVarint, 1)), nil
		} else {
			return fmt.Sprintf(`
	// type: %s%s
	if %s {
		output.WriteBytes(1)
	} else {
		output.WriteBytes(0)
	}`, t, writeTag(preFieldNum, fieldNum, WireVarint, 1), name), nil
		}
	case SimpleType_FLOAT32:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if %s != 0 {%s
		output.WriteFixed32(math.Float32bits(%s))
	}`, t, name, writeTag(preFieldNum, fieldNum, WireFixed32, 1), name), MATH_PKG
		} else {
			return fmt.Sprintf(`
	// type: %s%s
	output.WriteFixed32(math.Float32bits(%s))`, t, writeTag(preFieldNum, fieldNum, WireFixed32, 0), name), MATH_PKG
		}
	case SimpleType_FLOAT64:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if %s != 0 {%s
		output.WriteFixed64(math.Float64bits(%s))
	}`, t, name, writeTag(preFieldNum, fieldNum, WireFixed64, 1), name), MATH_PKG
		} else {
			return fmt.Sprintf(`
	// type: %s%s
	output.WriteFixed64(math.Float64bits(%s))`, t, writeTag(preFieldNum, fieldNum, WireFixed64, 0), name), MATH_PKG
		}
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value: %d", t)
		return "", nil
	}
}

func (t *FixedPointType) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	if ignore {
		return fmt.Sprintf(`
	// type: %s
	if %s != %d {%s
		output.WriteVarint(uint64(%s))
	}`, t, name, t.Floor, writeTag(preFieldNum, fieldNum, WireVarint, 1), t.ToVarint(name)), nil
	} else {
		return fmt.Sprintf(`
	// type: %s%s
	output.WriteVarint(uint64(%s))`, t, writeTag(preFieldNum, fieldNum, WireVarint, 0), t.ToVarint(name)), nil
	}
}

func (t *EnumType) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	if ignore {
		return fmt.Sprintf(`
	// type: %s
	if %s != 0 {%s
		output.WriteVarint(uint64(%s))
	}`, t, name, writeTag(preFieldNum, fieldNum, WireVarint, 1), name), nil
	} else {
		return fmt.Sprintf(`
	// type: %s%s
	output.WriteVarint(uint64(%s))`, t, writeTag(preFieldNum, fieldNum, WireVarint, 0), name), nil
	}
}

func (t *InstanceType) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	if ignore {
		var zero string
		if t.IsPtr {
			zero = "nil"
		} else {
			zero = "0"
		}
		return fmt.Sprintf(`
	// type: %s
	if %s != %s {%s
		output.WriteVarint(uint64(%s.CachedSize()))
		%s.Serialize(output)
	}`, t, name, zero, writeTag(preFieldNum, fieldNum, WireBytes, 1), name, name), nil
	} else {
		return fmt.Sprintf(`
	// type: %s
	{%s
		output.WriteVarint(uint64(%s.CachedSize()))
		%s.Serialize(output)
	}`, t, writeTag(preFieldNum, fieldNum, WireBytes, 1), name, name), nil
	}
}

func (t *ListType) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	tempSize := TEMP_PREFIX
	if strings.HasPrefix(size, TEMP_PREFIX) {
		tempSize = size + "p"
	}
	var pkgs map[string]string

	if l, ok := t.E.(*ListType); ok {
		if l.E.IsPrimitive() {
			element_s, element_p := t.E.SerializeGo(size, "e", preFieldNum, fieldNum, true)
			pkgs = update(pkgs, element_p)
			return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		for _, e := range %s {
			// list element%s else {%s
				output.WriteBytes(0)
			}
		}
	}`, t, name, name, addIndent(element_s, 2), writeTag(preFieldNum, fieldNum, WireBytes, 3)), pkgs
		} else {
			bytesize_s, bytesize_p := t.E.CachedSizeGo(tempSize, "e", "", 0, true)
			serialize_s, serialize_p := t.E.SerializeGo(tempSize, "e", "", 0, true)
			pkgs = update(pkgs, bytesize_p)
			pkgs = update(pkgs, serialize_p)

			return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		for _, e := range %s {
			%s := 0
			// list element size%s%s
			output.WriteVarint(uint64(%s))
			// list element serialize%s
		}
	}`, t, name, name, tempSize, addIndent(bytesize_s, 2), writeTag(preFieldNum, fieldNum, WireBytes, 2), tempSize, addIndent(serialize_s, 2)), pkgs
		}
	} else if !t.E.IsPrimitive() {
		element_s, element_p := t.E.SerializeGo(size, "e", "", 0, true)
		pkgs = update(pkgs, element_p)
		var checkNil string
		if _, ok := t.E.(*VariantType); !ok {
			checkNil = `
				log.Printf("[Tygo][Serialize] Nil in a list is treated as an empty object contents default properties!")`
			pkgs = update(pkgs, LOG_PKG)
		}

		return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		for _, e := range %s {
			// list element%s%s else {%s
				output.WriteBytes(0)
			}
		}
	}`, t, name, name, writeTag(preFieldNum, fieldNum, WireBytes, 2), addIndent(element_s, 2), checkNil), pkgs

	} else {
		var bytesize_s string
		var bytesize_p map[string]string
		if st, ok := t.E.(SimpleType); ok && st == SimpleType_BOOL {
			bytesize_s = fmt.Sprintf(`
		%s := len(%s)`, tempSize, name)
		} else if ok && st == SimpleType_FLOAT32 {
			bytesize_s = fmt.Sprintf(`
		%s := len(%s) * 4`, tempSize, name)
		} else if ok && st == SimpleType_FLOAT64 {
			bytesize_s = fmt.Sprintf(`
		%s := len(%s) * 8`, tempSize, name)
		} else {
			bytesize_s, bytesize_p = t.E.CachedSizeGo(tempSize, "e", "", 0, false)
			bytesize_s = fmt.Sprintf(`
		%s := 0
		for _, e := range %s {
			// list element size%s
		}`, tempSize, name, addIndent(bytesize_s, 2))
			pkgs = update(pkgs, bytesize_p)
		}
		serialize_s, serialize_p := t.E.SerializeGo(tempSize, "e", "", 0, false)
		pkgs = update(pkgs, serialize_p)

		return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {%s%s
		output.WriteVarint(uint64(%s))
		for _, e := range %s {
			// list element serialize%s
		}
	}`, t, name, bytesize_s, writeTag(preFieldNum, fieldNum, WireBytes, 1), tempSize, name, addIndent(serialize_s, 2)), pkgs
	}
}

func (t *DictType) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	tempSize := TEMP_PREFIX
	if strings.HasPrefix(size, TEMP_PREFIX) {
		tempSize = size + "p"
	}
	key_b_s, key_b_p := t.K.CachedSizeGo(tempSize, "k", "", 1, true)
	key_s_s, key_s_p := t.K.SerializeGo(tempSize, "k", "", 1, true)
	value_b_s, value_b_p := t.V.CachedSizeGo(tempSize, "v", "", 2, true)
	value_s_s, value_s_p := t.V.SerializeGo(tempSize, "v", "", 2, true)
	var pkgs map[string]string
	pkgs = update(pkgs, key_b_p)
	pkgs = update(pkgs, key_s_p)
	pkgs = update(pkgs, value_b_p)
	pkgs = update(pkgs, value_s_p)

	return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		for k, v := range %s {
			%s := 0
			// dict key size%s
			// dict value size%s%s
			output.WriteVarint(uint64(%s))
			// dict key serialize%s
			// dict value serialize%s
		}
	}`, t, name, name, tempSize, addIndent(key_b_s, 2), addIndent(value_b_s, 2), writeTag(preFieldNum, fieldNum, WireBytes, 2), tempSize, addIndent(key_s_s, 2), addIndent(value_s_s, 2)), pkgs
}

func (t *VariantType) SerializeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	tempSize := TEMP_PREFIX
	if strings.HasPrefix(size, TEMP_PREFIX) {
		tempSize = size + "p"
	}
	var bytesize_cases []string
	var serialize_cases []string
	var pkgs map[string]string
	tagInteger := 0
	tagFloat32 := 0
	tagFloat64 := 0

	for i, st := range t.Ts {
		type_s, type_p := st.Go()
		if type_s == "nil" {
			continue
		} else if t, ok := st.(SimpleType); ok {
			switch t {
			case SimpleType_INT32:
				fallthrough
			case SimpleType_INT64:
				fallthrough
			case SimpleType_UINT32:
				fallthrough
			case SimpleType_UINT64:
				tagInteger = i + 1
			case SimpleType_FLOAT32:
				tagFloat32 = i + 1
			case SimpleType_FLOAT64:
				tagFloat64 = i + 1
			}
		}

		bytesize_s, bytesize_p := st.CachedSizeGo(tempSize, "v", "", i+1, false)
		serialize_s, serialize_p := st.SerializeGo(tempSize, "v", "", i+1, false)
		pkgs = update(pkgs, type_p)
		pkgs = update(pkgs, bytesize_p)
		pkgs = update(pkgs, serialize_p)
		bytesize_cases = append(bytesize_cases, fmt.Sprintf(`
		// variant type size: %s
		case %s:%s`, st, type_s, addIndent(bytesize_s, 2)))
		serialize_cases = append(serialize_cases, fmt.Sprintf(`
		// variant type serialize: %s
		case %s:%s`, st, type_s, addIndent(serialize_s, 2)))
	}

	if tagInteger != 0 {
		bytesize_cases = append(bytesize_cases, fmt.Sprintf(`
		// addition type size: int
		case int:
			%s += %d + tygo.SizeVarint(uint64(v))`, tempSize, TAG_SIZE(tagInteger)))
		serialize_cases = append(serialize_cases, fmt.Sprintf(`
		// addition type serialize: int
		case int:%s
			output.WriteVarint(uint64(v))`, writeTag("", tagInteger, WireVarint, 2)))
	} else if tagFloat32 != 0 {
		bytesize_cases = append(bytesize_cases, fmt.Sprintf(`
		// addition type size: int -> float32
		case int:
			%s += %d`, tempSize, TAG_SIZE(tagFloat32)+4))
		serialize_cases = append(serialize_cases, fmt.Sprintf(`
		// addition type serialize: int -> float32
		case int:%s
			output.WriteFixed32(math.Float32bits(float32(v)))`, writeTag("", tagFloat32, WireFixed32, 2)))
		pkgs = update(pkgs, MATH_PKG)
	} else if tagFloat64 != 0 {
		bytesize_cases = append(bytesize_cases, fmt.Sprintf(`
		// addition type size: int -> float64
		case int:
			%s += %d`, tempSize, TAG_SIZE(tagFloat64)+8))
		serialize_cases = append(serialize_cases, fmt.Sprintf(`
		// addition type serialize: int -> float64
		case int:%s
			output.WriteFixed64(math.Float64bits(float64(v)))`, writeTag("", tagFloat64, WireFixed64, 2)))
		pkgs = update(pkgs, MATH_PKG)
	}

	if tagFloat32 != 0 && tagFloat64 == 0 {
		bytesize_cases = append(bytesize_cases, fmt.Sprintf(`
		// addition type size: float64 -> float32
		case float64:
			%s += %d`, tempSize, TAG_SIZE(tagFloat32)+4))
		serialize_cases = append(serialize_cases, fmt.Sprintf(`
		// addition type serialize: float64 -> float32
		case float64:%s
			output.WriteFixed32(math.Float32bits(float32(v)))`, writeTag("", tagFloat32, WireFixed32, 2)))
		pkgs = update(pkgs, MATH_PKG)
	}

	var compareZero string
	if ignore {
		compareZero = fmt.Sprintf("if %s != nil ", name)
	}

	return fmt.Sprintf(`
	// type: %s
	%s{
		%s := 0
		switch v := %s.(type) {%s
		default:
			panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for %s: %%v", v))
		}%s
		output.WriteVarint(uint64(%s))
		switch v := %s.(type) {%s
		default:
			panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for %s: %%v", v))
		}
	}`, t, compareZero, tempSize, name, strings.Join(bytesize_cases, ""), t,
		writeTag(preFieldNum, fieldNum, WireBytes, 1), tempSize, name,
		strings.Join(serialize_cases, ""), t), pkgs
}

//=============================================================================
