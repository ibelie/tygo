// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"log"
	"strings"
)

func (t *Enum) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
	return fmt.Sprintf(`
	x, err := input.ReadVarint()
	*%s = %s(x)`, name, t.Name), nil
}

func (t *Method) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
	return "", nil
}

func (t *Object) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
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

	return fmt.Sprintf(`
	if %s != nil {%s
	}`, name, strings.Join(fields, "")), pkgs
}

func (t UnknownType) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
	return "", nil
}

func (t SimpleType) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
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
	if x, err := input.ReadVarint(); err == nil {
		%s = %s(x)
	} else {
		return err
	}`, t, name, t), nil
	case SimpleType_BYTES:
		fallthrough
	case SimpleType_STRING:
		return fmt.Sprintf(`
	// type: %s
	if len(%s) > 0 {
		l := len([]byte(%s))
		%s += %stygo.SizeVarint(uint64(l)) + l
	}`, t, name, name, size, tagsize_s), nil
	case SimpleType_BOOL:
		return fmt.Sprintf(`
	// type: %s
	if x, err := input.ReadByte(); err == nil {
		%s = x != 0
	} else {
		return err
	}`, t, name), nil
	case SimpleType_FLOAT32:
		return fmt.Sprintf(`
	// type: %s
	if x, err := input.ReadFixed32(); err == nil {
		%s = math.Float32frombits(x)
	} else {
		return err
	}`, t, name), MATH_PKG
	case SimpleType_FLOAT64:
		return fmt.Sprintf(`
	// type: %s
	if x, err := input.ReadFixed64(); err == nil {
		%s = math.Float64frombits(x)
	} else {
		return err
	}`, t, name), MATH_PKG
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value: %d", t)
		return "", nil
	}
}

func (t *FixedPointType) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
	return fmt.Sprintf(`
	// type: %s
	if x, err := input.ReadVarint(); err == nil {
		%s = float64(x) / %d + %d
	} else {
		return err
	}`, t, name, pow10(t.Precision), t.Floor), nil
}

func (t *EnumType) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
	return fmt.Sprintf(`
	// type: %s
	if x, err := input.ReadVarint(); err == nil {
		%s = %s(x)
	} else {
		return err
	}`, t, name, t.Name), nil
}

func (t *InstanceType) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
	return fmt.Sprintf(`
	// type: %s
	if err := %s.Deserialize(input); err != nil {
		return err
	}`, t, name), nil
}

func (t *ListType) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, false)
}

func (t *DictType) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, false)
}

func (t *VariantType) DeserializeGo(name string, preFieldNum string, fieldNum int) (string, map[string]string) {
	return t._ByteSizeGo(size, name, preFieldNum, fieldNum, ignore, false)
}
