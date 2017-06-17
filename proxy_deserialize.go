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

func desVar() string {
	desVarCount++
	return fmt.Sprintf("tmp_%d", desVarCount)
}

func tagInt(preFieldNum string, fieldNum int, wireType WireType) string {
	if preFieldNum == "" {
		return strconv.Itoa(_MAKE_TAG(fieldNum, wireType))
	} else {
		return _MAKE_TAG_STR(fmt.Sprintf(`(%s + %d)`, preFieldNum, fieldNum), wireType)
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
			return fmt.Sprintf(`ExpectBytes(%s) { // tag: %d MAKE_TAG(%d, %s=%d)`,
				strings.Join(tagbytes, ", "), _MAKE_TAG(fieldNum, wireType),
				fieldNum, wireType, wireType)
		}
	} else {
		return fmt.Sprintf(`ExpectTag(%s + %d, %d) { // %s`,
			preFieldNum, fieldNum, wireType, wireType)
	}
}

func (t *Enum) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	return fmt.Sprintf(`
	x, err := %s.ReadVarint()
	*%s = %s(x)`, input, name, t.Name), WireVarint, nil
}

func (t *Method) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	return "", WireVarint, nil
}

func (t *Object) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	var pkgs map[string]string
	var parents []string
	if t.HasParent() {
		parents = append(parents, fmt.Sprintf(`
		if err := %s.%s.Deserialize(%s); err == nil {
			%s.Reset()
		} else {
			return err
		}`, name, t.Parent.Name, input, input))
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
		parents = append(parents, fmt.Sprintf(`
		preFieldNum := %s.%s.MaxFieldNum()`, name, p_name))
		p_name = "preFieldNum"
	}

	d1 := desVar()
	var d string
	var field_s string
	var field_w WireType
	var field_p map[string]string
	var fields []string
	for i, field := range t.Fields {
		if i == 0 {
			field_s, field_w, field_p = field.DeserializeGo("tag", input, fmt.Sprintf("%s.%s", name, field.Name), p_name, p_num+i+1)
			pkgs = update(pkgs, field_p)
		}

		var label string
		if p_name == "" {
			label = fmt.Sprintf(" // MAKE_TAG(%d, %s=%d)", p_num+i+1, field_w, field_w)
		}
		if d != "" {
			label = fmt.Sprintf(`%s
			object_%s:`, label, d)
			d = ""
		}

		var next string
		if i < len(t.Fields)-1 {
			next_s, next_w, next_p := t.Fields[i+1].DeserializeGo("tag", input, fmt.Sprintf("%s.%s", name, t.Fields[i+1].Name), p_name, p_num+i+2)
			pkgs = update(pkgs, next_p)
			d = desVar()
			next = fmt.Sprintf(`
				if %s.%s
					goto object_%s // goto case %d
				}`, input, expectTag(p_name, p_num+i+2, next_w), d, i+2)
			field_s, field_w, field_p = next_s, next_w, next_p
		} else {
			next = fmt.Sprintf(`
				if %s.ExpectEnd() {
					break object_%s // end for %s
				}`, input, d1, t.Name)
		}

		var listTag string
		if l, ok := field.Type.(*ListType); ok && l.E.IsPrimitive() {
			listTag = fmt.Sprintf(" || tag == %s", tagInt(p_name, p_num+i+1, WireBytes))
		}

		fields = append(fields, fmt.Sprintf(`
		// property: %s.%s
		case %d:
			if tag == %s%s {%s%s
				continue object_%s // next tag for %s%s
			}`, name, field.Name, i+1, tagInt(p_name, p_num+i+1, field_w), listTag, label,
			addIndent(field_s, 3), d1, t.Name, next))
	}

	var cutoff string
	if p_name == "" {
		cutoff = strconv.Itoa(_MAKE_CUTOFF(p_num + len(t.Fields)))
	} else {
		cutoff = _MAKE_CUTOFF_STR(fmt.Sprintf("(%s + %d)", p_name, p_num+len(t.Fields)))
	}

	var switchFlag string
	if p_name == "" {
		if p_num == 0 {
			switchFlag = _TAG_FIELD_STR("tag")
		} else {
			switchFlag = fmt.Sprintf("(%s) - %d", _TAG_FIELD_STR("tag"), p_num)
		}
	} else {
		if p_num == 0 {
			switchFlag = fmt.Sprintf("(%s) - %s", _TAG_FIELD_STR("tag"), p_name)
		} else {
			switchFlag = fmt.Sprintf("(%s) - %s - %d", _TAG_FIELD_STR("tag"), p_name, p_num)
		}
	}

	return fmt.Sprintf(`%s
	object_%s: for !%s.ExpectEnd() {
		var tag int
		if tag, err = %s.ReadTag(%s); err != nil {
			return
		}
		switch %s {%s
		}
		if err = %s.SkipField(tag); err != nil {
			return
		}
	}`, strings.Join(parents, ""), d1, input, input, cutoff, switchFlag, strings.Join(fields, ""), input), WireBytes, pkgs
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
	if x, err := %s.ReadBuf(); err == nil {
		%s = make([]byte, len(x))
		copy(%s, x)
	} else {
		return err
	}`, t, input, name, name), WireBytes, nil
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
	} else if err := %s.Deserialize(&tygo.ProtoBuf{Buffer: x}); err != nil {
		return err
	}`, t, input, name), WireBytes, updateTygo(nil)
}

func (t *ListType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	tempInput := TEMP_PREFIX + "i"
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
		list_s := fmt.Sprintf(`
	if x, err := %s.ReadBuf(); err == nil {
		%s := &tygo.ProtoBuf{Buffer: x}
		var %s %s
		for !%s.ExpectEnd() {%s
		}
		%s = append(%s, %s)
	} else {
		return err
	}`, input, tempInput, v, type_s, tempInput, addIndent(element_s, 2), name, name, v)
		if tag_s == "" {
			return fmt.Sprintf(`
	// type: %s%s`, t, list_s), WireBytes, pkgs
		} else {
			return fmt.Sprintf(`
	// type: %s
	loop_%s: for {%s
		if !%s.%s
			break loop_%s // end for %s
		}
	}`, t, v, addIndent(list_s, 1), input, tag_s, v, t), WireBytes, pkgs
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
	loop_%s: for {
		var %s %s%s
		%s = append(%s, %s)
		if !%s.%s
			break loop_%s // end for %s
		}
	}`, t, v, v, type_s, addIndent(element_s, 1), name, name, v, input, tag_s, v, t), WireBytes, pkgs
		}
	} else {
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
	}`, t, input, tempInput, v, type_s, tempInput, addIndent(element_s, 2), name, name, v), element_w, pkgs
		} else {
			return fmt.Sprintf(`
	// type: %s
	if %s == %s {
		loop_%s: for {
			var %s %s%s
			%s = append(%s, %s)
			if !%s.%s
				break loop_%s // end for %s
			}
		}
	} else if x, err := %s.ReadBuf(); err == nil {
		%s := &tygo.ProtoBuf{Buffer: x}
		var %s %s
		for !%s.ExpectEnd() {%s
		}
		%s = append(%s, %s)
	} else {
		return err
	}`, t, tag, tagInt(preFieldNum, fieldNum, element_w), v, v, type_s,
				addIndent(element_s, 2), name, name, v, input, tag_s, v, t, input, tempInput,
				v, type_s, tempInput, addIndent(element_s, 2), name, name, v), element_w, pkgs
		}
	}
}

func (t *DictType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	tempInput := TEMP_PREFIX + "i"
	if strings.HasPrefix(input, TEMP_PREFIX) {
		tempInput = input + "i"
	}
	tempTag := TEMP_PREFIX + "g"
	if strings.HasPrefix(input, TEMP_PREFIX) {
		tempTag = input + "g"
	}

	k := desVar()
	key_t_s, key_t_p := t.K.Go()
	key_d_s, key_d_w, key_d_p := t.K.DeserializeGo(tempTag, tempInput, k, "", 1)
	v := desVar()
	value_t_s, value_t_p := t.V.Go()
	value_d_s, value_d_w, value_d_p := t.V.DeserializeGo(tempTag, tempInput, v, "", 2)
	value_e := expectTag("", 2, value_d_w)
	var pkgs map[string]string
	pkgs = update(pkgs, key_t_p)
	pkgs = update(pkgs, value_t_p)
	pkgs = update(pkgs, key_d_p)
	pkgs = update(pkgs, value_d_p)
	tag_s := expectTag(preFieldNum, fieldNum, WireBytes)

	dict_s := fmt.Sprintf(`
	if x, err := %s.ReadBuf(); err == nil {
		%s := &tygo.ProtoBuf{Buffer: x}
		var %s %s
		var %s %s
		dict_%s: for !%s.ExpectEnd() {
			%s, err := %s.ReadTag(%d)
			if err != nil {
				return err
			}
			switch %s {
			// dict key
			case 1:
				if %s == %d { // MAKE_TAG(1, %s=%d)%s
					if %s.%s
						goto dict_%s // goto case 2
					}
					continue dict_%s // next tag for %s
				}
			case 2:
				if %s == %d { // MAKE_TAG(2, %s=%d)
				dict_%s:%s
					if %s.ExpectEnd() {
						break dict_%s // end for %s
					}
					continue dict_%s // next tag for %s
				}
			}
			if err := %s.SkipField(%s); err != nil {
				return err
			}
		}
		%s[%s] = %s
	} else {
		return err
	}`, input, tempInput, k, key_t_s, v, value_t_s, k, tempInput, tempTag, tempInput,
		_MAKE_CUTOFF(2), _TAG_FIELD_STR(tempTag), tempTag, _MAKE_TAG(1, key_d_w),
		key_d_w, key_d_w, addIndent(key_d_s, 4), tempInput, value_e, v, k, t, tempTag,
		_MAKE_TAG(2, value_d_w), value_d_w, value_d_w, v, addIndent(value_d_s, 4),
		tempInput, k, t, k, t, tempInput, tempTag, name, k, v)

	if tag_s == "" {
		return fmt.Sprintf(`
	// type: %s%s`, t, dict_s), WireBytes, pkgs
	} else {
		return fmt.Sprintf(`
	// type: %s
	loop_%s: for {%s
		if !%s.%s
			break loop_%s // end for %s
		}
	}`, t, k, addIndent(dict_s, 1), input, tag_s, k, t), WireBytes, pkgs
	}
}

func (t *VariantType) DeserializeGo(tag string, input string, name string, preFieldNum string, fieldNum int) (string, WireType, map[string]string) {
	v := desVar()
	tempInput := TEMP_PREFIX + "i"
	if strings.HasPrefix(input, TEMP_PREFIX) {
		tempInput = input + "i"
	}
	tempTag := TEMP_PREFIX + "g"
	if strings.HasPrefix(input, TEMP_PREFIX) {
		tempTag = input + "g"
	}
	var pkgs map[string]string

	var cases []string
	for i, ts := range t.Ts {
		if s, ok := ts.(SimpleType); ok && s == SimpleType_NIL {
			continue
		}
		variant_s, variant_w, variant_p := ts.DeserializeGo(tempTag, tempInput, name, "", i+1)
		pkgs = update(pkgs, variant_p)
		cases = append(cases, fmt.Sprintf(`
		case %d:
			if %s == %d { // MAKE_TAG(%d, %s=%d)%s
				continue variant_%s // next tag for %s
			}`, i+1, tempTag, _MAKE_TAG(i+1, variant_w), i+1, variant_w, variant_w, addIndent(variant_s, 3), v, t))
	}

	return fmt.Sprintf(`
	// type: %s
	if x, err := %s.ReadBuf(); err == nil {
		%s := &tygo.ProtoBuf{Buffer: x}
		variant_%s: for !%s.ExpectEnd() {
			%s, err := %s.ReadTag(%d)
			if err != nil {
				return err
			}
			switch %s {%s
			}
			if err := %s.SkipField(%s); err != nil {
				return err
			}
		}
	} else {
		return err
	}`, t, input, tempInput, v, tempInput, tempTag, tempInput, _MAKE_CUTOFF(len(t.Ts)),
		_TAG_FIELD_STR(tempTag), strings.Join(cases, ""), tempInput, tempTag), WireBytes, pkgs
}
