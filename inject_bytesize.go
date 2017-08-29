// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"log"
	"strings"
)

func tagSize(preFieldNum string, fieldNum int) (string, map[string]string) {
	if preFieldNum == "" {
		if fieldNum <= 0 {
			return "", nil
		} else {
			return fmt.Sprintf("%d + ", TAG_SIZE(fieldNum)), nil
		}
	} else {
		return fmt.Sprintf("tygo.TAG_SIZE(%s + %d) + ", preFieldNum, fieldNum), updateTygo(nil)
	}
}

func (t *Enum) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return fmt.Sprintf(`
	if %s != 0 {
		%s = tygo.SizeVarint(uint64(%s))
	}`, name, size, name), updateTygo(nil)
}

func (t *Object) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	var pkgs map[string]string
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
		p_name = "preFieldNum"
	}

	for i, field := range t.VisibleFields() {
		field_s, field_p := field.ByteSizeGo(size, fmt.Sprintf("%s.%s", name, field.Name), p_name, p_num+i+1, true)
		pkgs = update(pkgs, field_p)
		fields = append(fields, fmt.Sprintf(`
		// property: %s.%s%s
`, name, field.Name, addIndent(field_s, 1)))
	}

	return fmt.Sprintf(`
	if %s != nil {%s
		%s.SetCachedSize(%s)
	}`, name, strings.Join(fields, ""), name, size), pkgs
}

func (t UnknownType) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return "", nil
}

func (t SimpleType) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	tagsize_s, tagsize_p := tagSize(preFieldNum, fieldNum)
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
	if %s != 0 {
		%s += %stygo.SizeVarint(uint64(%s))
	}`, t, name, size, tagsize_s, name), updateTygo(tagsize_p)
		} else {
			return fmt.Sprintf(`
	// type: %s
	%s += %stygo.SizeVarint(uint64(%s))`, t, size, tagsize_s, name), updateTygo(tagsize_p)
		}
	case SimpleType_BYTES:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		%s += %stygo.SizeBuffer(%s)
	}`, t, name, size, tagsize_s, name), updateTygo(tagsize_p)
		} else {
			return fmt.Sprintf(`
	// type: %s
	%s += %stygo.SizeBuffer(%s)`, t, size, tagsize_s, name), updateTygo(tagsize_p)
		}
	case SimpleType_STRING:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		%s += %stygo.SizeBuffer([]byte(%s))
	}`, t, name, size, tagsize_s, name), updateTygo(tagsize_p)
		} else {
			return fmt.Sprintf(`
	// type: %s
	%s += %stygo.SizeBuffer([]byte(%s))`, t, size, tagsize_s, name), updateTygo(tagsize_p)
		}
	case SimpleType_SYMBOL:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		%s += %stygo.SizeSymbol(%s)
	}`, t, name, size, tagsize_s, name), updateTygo(tagsize_p)
		} else {
			return fmt.Sprintf(`
	// type: %s
	%s += %stygo.SizeSymbol(%s)`, t, size, tagsize_s, name), updateTygo(tagsize_p)
		}
	case SimpleType_BOOL:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if %s {
		%s += %s1
	}`, t, name, size, tagsize_s), tagsize_p
		} else {
			return fmt.Sprintf(`
	// type: %s
	%s += %s1`, t, size, tagsize_s), tagsize_p
		}
	case SimpleType_FLOAT32:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if %s != 0 {
		%s += %s4
	}`, t, name, size, tagsize_s), tagsize_p
		} else {
			return fmt.Sprintf(`
	// type: %s
	%s += %s4`, t, size, tagsize_s), tagsize_p
		}
	case SimpleType_FLOAT64:
		if ignore {
			return fmt.Sprintf(`
	// type: %s
	if %s != 0 {
		%s += %s8
	}`, t, name, size, tagsize_s), tagsize_p
		} else {
			return fmt.Sprintf(`
	// type: %s
	%s += %s8`, t, size, tagsize_s), tagsize_p
		}
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for ByteSize: %d", t)
		return "", nil
	}
}

func (t *FixedPointType) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	tagsize_s, tagsize_p := tagSize(preFieldNum, fieldNum)
	if ignore {
		return fmt.Sprintf(`
	// type: %s
	if %s != %d {
		%s += %stygo.SizeVarint(uint64(%s))
	}`, t, name, t.Floor, size, tagsize_s, t.ToVarint(name)), updateTygo(tagsize_p)
	} else {
		return fmt.Sprintf(`
	// type: %s
	%s += %stygo.SizeVarint(uint64(%s))`, t, size, tagsize_s, t.ToVarint(name)), updateTygo(tagsize_p)
	}
}

func (t *EnumType) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	tagsize_s, tagsize_p := tagSize(preFieldNum, fieldNum)
	if ignore {
		return fmt.Sprintf(`
	// type: %s
	if %s != 0 {
		%s += %stygo.SizeVarint(uint64(%s))
	}`, t, name, size, tagsize_s, name), updateTygo(tagsize_p)
	} else {
		return fmt.Sprintf(`
	// type: %s
	%s += %stygo.SizeVarint(uint64(%s))`, t, size, tagsize_s, name), updateTygo(tagsize_p)
	}
}

func (t *InstanceType) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, false)
}

func (t *ListType) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, false)
}

func (t *DictType) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, false)
}

func (t *VariantType) ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, false)
}

//=============================================================================

func (t *InstanceType) _ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool, isCached bool) (string, map[string]string) {
	tagsize_s, tagsize_p := tagSize(preFieldNum, fieldNum)
	tempSize := TEMP_PREFIX
	if strings.HasPrefix(size, TEMP_PREFIX) {
		tempSize = size + "p"
	}
	bytesizeMethod := "ByteSize"
	if isCached {
		bytesizeMethod = "CachedSize"
	}
	if ignore {
		var zero string
		if t.IsPtr {
			zero = "nil"
		} else {
			zero = "0"
		}
		return fmt.Sprintf(`
	// type: %s
	if %s != %s {
		%s := %s.%s()
		%s += %stygo.SizeVarint(uint64(%s)) + %s
	}`, t, name, zero, tempSize, name, bytesizeMethod, size, tagsize_s, tempSize, tempSize), updateTygo(tagsize_p)
	} else {
		return fmt.Sprintf(`
	// type: %s
	{
		%s := %s.%s()
		%s += %stygo.SizeVarint(uint64(%s)) + %s
	}`, t, tempSize, name, bytesizeMethod, size, tagsize_s, tempSize, tempSize), updateTygo(tagsize_p)
	}
}

func (t *ListType) _ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool, isCached bool) (string, map[string]string) {
	tempSize := TEMP_PREFIX
	if strings.HasPrefix(size, TEMP_PREFIX) {
		tempSize = size + "p"
	}
	bytesizeMethod := t.E.ByteSizeGo
	if isCached {
		bytesizeMethod = t.E.CachedSizeGo
	}
	var pkgs map[string]string
	tagsize_s, tagsize_p := tagSize(preFieldNum, fieldNum)
	pkgs = update(pkgs, tagsize_p)

	if st, ok := t.E.(SimpleType); ok {
		fixedSize := 1
		switch st {
		case SimpleType_FLOAT64:
			fixedSize *= 2
			fallthrough
		case SimpleType_FLOAT32:
			fixedSize *= 4
			fallthrough
		case SimpleType_BOOL:
			return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		%s := len(%s) * %d
		%s += %stygo.SizeVarint(uint64(%s)) + %s
	}`, t, name, tempSize, name, fixedSize, size, tagsize_s, tempSize, tempSize), updateTygo(pkgs)
		}
	}

	if t.E.IsIterative() {
		element_s, element_p := bytesizeMethod(tempSize, "e", "", 0, true)
		pkgs = update(pkgs, element_p)
		return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		for _, e := range %s {
			%s := 0
			// list element%s
			%s += %stygo.SizeVarint(uint64(%s)) + %s
		}
	}`, t, name, name, tempSize, addIndent(element_s, 2), size, tagsize_s, tempSize, tempSize), updateTygo(pkgs)

	} else if !t.E.IsPrimitive() {
		element_s, element_p := bytesizeMethod(size, "e", preFieldNum, fieldNum, true)
		pkgs = update(pkgs, element_p)
		var checkNil string
		switch t.E.(type) {
		case *VariantType:
		case *ListType:
		default:
			checkNil = `
				log.Printf("[Tygo][ByteSize] Nil in a list is treated as an empty object contents default properties!")`
			pkgs = update(pkgs, LOG_PKG)
		}

		return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		for _, e := range %s {
			// list element%s else {%s
				%s += %s1
			}
		}
	}`, t, name, name, addIndent(element_s, 2), checkNil, size, tagsize_s), pkgs

	} else {
		element_s, element_p := bytesizeMethod(tempSize, "e", "", 0, false)
		pkgs = update(pkgs, element_p)

		return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		%s := 0
		for _, e := range %s {
			// list element%s
		}
		%s += %stygo.SizeVarint(uint64(%s)) + %s
	}`, t, name, tempSize, name, addIndent(element_s, 2), size, tagsize_s, tempSize, tempSize), updateTygo(pkgs)
	}
}

func (t *DictType) _ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool, isCached bool) (string, map[string]string) {
	tempSize := TEMP_PREFIX
	if strings.HasPrefix(size, TEMP_PREFIX) {
		tempSize = size + "p"
	}
	bytesizeMethod := t.K.ByteSizeGo
	if isCached {
		bytesizeMethod = t.K.CachedSizeGo
	}
	key_s, key_p := bytesizeMethod(tempSize, "k", "", 1, true)
	bytesizeMethod = t.V.ByteSizeGo
	if isCached {
		bytesizeMethod = t.V.CachedSizeGo
	}
	value_s, value_p := bytesizeMethod(tempSize, "v", "", 2, true)
	tagsize_s, tagsize_p := tagSize(preFieldNum, fieldNum)
	pkgs := updateTygo(nil)
	pkgs = update(pkgs, key_p)
	pkgs = update(pkgs, value_p)
	pkgs = update(pkgs, tagsize_p)

	return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		for k, v := range %s {
			%s := 0
			// dict key%s
			// dict value%s
			%s += %stygo.SizeVarint(uint64(%s)) + %s
		}
	}`, t, name, name, tempSize, addIndent(key_s, 2), addIndent(value_s, 2), size, tagsize_s, tempSize, tempSize), pkgs
}

func (t *VariantType) _ByteSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool, isCached bool) (string, map[string]string) {
	tempSize := TEMP_PREFIX
	if strings.HasPrefix(size, TEMP_PREFIX) {
		tempSize = size + "p"
	}
	var cases []string
	tagInteger := 0
	tagFloat32 := 0
	tagFloat64 := 0
	tagsize_s, tagsize_p := tagSize(preFieldNum, fieldNum)
	pkgs := updateTygo(nil)
	pkgs = update(pkgs, LOG_PKG)
	pkgs = update(pkgs, tagsize_p)

	variantNum := 0
	for _, st := range t.Ts {
		type_s, type_p := st.Go()
		if type_s == "nil" {
			continue
		}
		variantNum++
		if t, ok := st.(SimpleType); ok {
			switch t {
			case SimpleType_INT32:
				fallthrough
			case SimpleType_INT64:
				fallthrough
			case SimpleType_UINT32:
				fallthrough
			case SimpleType_UINT64:
				tagInteger = variantNum
			case SimpleType_FLOAT32:
				tagFloat32 = variantNum
			case SimpleType_FLOAT64:
				tagFloat64 = variantNum
			}
		}

		bytesizeMethod := st.ByteSizeGo
		if isCached {
			bytesizeMethod = st.CachedSizeGo
		}
		variant_s, variant_p := bytesizeMethod(tempSize, "v", "", variantNum, false)
		cases = append(cases, fmt.Sprintf(`
		// variant type: %s
		case %s:%s`, st, type_s, addIndent(variant_s, 2)))
		pkgs = update(pkgs, type_p)
		pkgs = update(pkgs, variant_p)
	}

	if tagInteger != 0 {
		cases = append(cases, fmt.Sprintf(`
		// addition type: int
		case int:
			%s += %d + tygo.SizeVarint(uint64(v))`, tempSize, TAG_SIZE(tagInteger)))
	} else if tagFloat32 != 0 {
		cases = append(cases, fmt.Sprintf(`
		// addition type: int -> float32
		case int:
			%s += %d`, tempSize, TAG_SIZE(tagFloat32)+4))
	} else if tagFloat64 != 0 {
		cases = append(cases, fmt.Sprintf(`
		// addition type: int -> float64
		case int:
			%s += %d`, tempSize, TAG_SIZE(tagFloat64)+8))
	}

	if tagFloat32 != 0 && tagFloat64 == 0 {
		cases = append(cases, fmt.Sprintf(`
		// addition type: float64 -> float32
		case float64:
			%s += %d`, tempSize, TAG_SIZE(tagFloat32)+4))
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
			log.Panicf("[Tygo][Variant] Unexpect type for %s: %%v", v)
		}
		%s += %stygo.SizeVarint(uint64(%s)) + %s
	}`, t, compareZero, tempSize, name, strings.Join(cases, ""), t, size, tagsize_s, tempSize, tempSize), pkgs
}

//=============================================================================

func (t *Enum) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t.ByteSizeGo(size, name, preFieldNum, fieldNum, ignore)
}

func (t *Object) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t.ByteSizeGo(size, name, preFieldNum, fieldNum, ignore)
}

func (t UnknownType) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t.ByteSizeGo(size, name, preFieldNum, fieldNum, ignore)
}

func (t SimpleType) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t.ByteSizeGo(size, name, preFieldNum, fieldNum, ignore)
}

func (t *FixedPointType) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t.ByteSizeGo(size, name, preFieldNum, fieldNum, ignore)
}

func (t *EnumType) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t.ByteSizeGo(size, name, preFieldNum, fieldNum, ignore)
}

func (t *InstanceType) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, true)
}

func (t *ListType) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, true)
}

func (t *DictType) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, true)
}

func (t *VariantType) CachedSizeGo(size string, name string, preFieldNum string, fieldNum int, ignore bool) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, true)
}
