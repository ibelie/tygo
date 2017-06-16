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

var desVarCount int

func desVar() {
	desVarCount++
	return fmt.Sprintf("tmp_%d", desVarCount)
}

func tagInt(preFieldNum string, fieldNum int, wireType WireType) string {
	if preFieldNum == "" {
		if fieldNum <= 0 {
			return ""
		} else {
			return strconv.Itoa(MAKE_TAG(fieldNum, wireType))
		}
	} else {
		return fmt.Sprintf(`(((%s + %d) << %d) | %d)`, preFieldNum, fieldNum, WireTypeBits, wireType)
	}
}

func expectTag(preFieldNum string, fieldNum int, wireType WireType) string {
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
			return fmt.Sprintf(`ExpectBytes(%s) // tag: %d MAKE_TAG(%d, %s=%d)`,
				strings.Join(tagbytes, ", "), MAKE_TAG(fieldNum, wireType), fieldNum, wireType, wireType)
		}
	} else {
		return fmt.Sprintf(`ExpectTag(%s + %d, %d)`, preFieldNum, fieldNum, wireType)
	}
}

func (t *Enum) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	if desVarCount != 0 {
		log.Fatalf("[Tygo][Enum] desVarCount(%d)", desVarCount)
	}
	desVarCount = 0
	return fmt.Sprintf(`
	x, err := %s.ReadVarint()
	*%s = %s(x)`, input, name, t.Name), WireVarint, nil
}

func (t *Method) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	return "", WireVarint, nil
}

func (t *Object) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	if desVarCount != 0 {
		log.Fatalf("[Tygo][Enum] desVarCount(%d)", desVarCount)
	}
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

	for i, field := range t.Fields {
		field_s, field_p := field.ByteSizeGo(size, fmt.Sprintf("%s.%s", name, field.Name), p_name, p_num+i+1, true)
		pkgs = update(pkgs, field_p)
		fields = append(fields, fmt.Sprintf(`
		// property: %s.%s%s
`, name, field.Name, addIndent(field_s, 1)))
	}

	desVarCount = 0
	return fmt.Sprintf(`
	if %s != nil {%s
	}`, name, strings.Join(fields, "")), WireBytes, pkgs
}

func (t UnknownType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	return "", WireVarint, nil
}

func (t SimpleType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	switch t {
	case SimpleType_INT32:
		fallthrough
	case SimpleType_INT64:
		fallthrough
	case SimpleType_UINT32:
		fallthrough
	case SimpleType_UINT64:
		return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadVarint(); err == nil {
		%s = %s(x)
	} else {
		return err
	}`, t, input, name, t), WireVarint, nil
	case SimpleType_BYTES:
		return fmt.Sprintf(`
	// type: %s
	if %s, err := %s.ReadBuf(); err != nil {
		return err
	}`, t, name, input), WireBytes, nil
	case SimpleType_STRING:
		return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadBuf(); err == nil {
		%s = %s(x)
	} else {
		return err
	}`, t, input, name, t), WireBytes, nil
	case SimpleType_BOOL:
		return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadByte(); err == nil {
		%s = x != 0
	} else {
		return err
	}`, t, input, name), WireVarint, nil
	case SimpleType_FLOAT32:
		return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadFixed32(); err == nil {
		%s = math.Float32frombits(x)
	} else {
		return err
	}`, t, input, name), WireFixed32, MATH_PKG
	case SimpleType_FLOAT64:
		return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadFixed64(); err == nil {
		%s = math.Float64frombits(x)
	} else {
		return err
	}`, t, input, name), WireFixed64, MATH_PKG
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value: %d", t)
		return "", WireVarint, nil
	}
}

func (t *FixedPointType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadVarint(); err == nil {
		%s = float64(x) / %d + %d
	} else {
		return err
	}`, t, input, name, pow10(t.Precision), t.Floor), WireVarint, nil
}

func (t *EnumType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadVarint(); err == nil {
		%s = %s(x)
	} else {
		return err
	}`, t, input, name, t.Name), WireVarint, nil
}

func (t *InstanceType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadBuf(); err != nil {
		return err
	} else if err := %s.Deserialize(&tygo.ProbuBuf{Buffer: x}); err != nil {
		return err
	}`, t, input, name), WireBytes, updateTygo(nil)
}

func (t *ListType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	tempInput := TEMP_PREFIX
	if strings.HasPrefix(input, TEMP_PREFIX) {
		tempInput = input + "i"
	}
	v := desVar()
	type_s, type_p := t.E.Go()
	var pkgs map[string]string
	pkgs = update(pkgs, type_p)

	if l, ok := t.E.(*ListType); ok && !l.E.IsPrimitive() {
		element_s, element_w, element_p := t.E.DeserializeGo(tag, tempInput, v, "", 0)
		pkgs = update(pkgs, element_p)
		tag_s := expectTag(preFieldNum, fieldNum, element_w)
		if tag_s == "" {
			return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadBuf(); err == nil {
		%s := &tygo.ProtoBuf{Buffer: x}
		var %s %s
		for !%s.ExpectEnd() {%s
		}
		%s = append(%s, %s)
	} else {
		return err
	}`, t, input, tempInput, v, type_s, tempInput, addIndent(element_s, 2), name, name, v), WireBytes, pkgs
		} else {
			return fmt.Sprintf(`
	// type: %s
	parse_%s_loop: for {
		if x, err := %s.ReadBuf(); err == nil {
			%s := &tygo.ProtoBuf{Buffer: x}
			var %s %s
			for !%s.ExpectEnd() {%s
			}
			%s = append(%s, %s)
			if !%s.%s {
				break parse_%s_loop
			}
		} else {
			return err
		}
	}`, t, v, input, tempInput, v, type_s, tempInput, addIndent(element_s, 3), name, name, v, input, tag_s, v), WireBytes, pkgs
		}
	} else if !t.E.IsPrimitive() {
		element_s, element_w, element_p := t.E.DeserializeGo(tag, input, v, "", 0)
		pkgs = update(pkgs, element_p)
		tag_s := expectTag(preFieldNum, fieldNum, element_w)
		if tag_s == "" {
			return fmt.Sprintf(`
	// type: %s
	var %s %s%s
	%s = append(%s, %s)`, t, v, type_s, element_s, name, name, v), WireBytes, pkgs
		} else {
			return fmt.Sprintf(`
	// type: %s
	parse_%s_loop: for {
		var %s %s%s
		%s = append(%s, %s)
		if !%s.%s {
			break parse_%s_loop
		}
	}`, t, v, v, type_s, addIndent(element_s, 1), name, name, v, input, tag_s, v), WireBytes, pkgs
		}
	} else {
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

func (t *DictType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	tempTag := TEMP_PREFIX
	if strings.HasPrefix(input, TEMP_PREFIX) {
		tempTag = input + "g"
	}
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, false)
}

func (t *VariantType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, false)
}
