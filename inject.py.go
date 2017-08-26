// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
)

var (
	PY_WRITER  io.Writer
	PY_TYPES   map[string]bool
	PY_OBJECTS map[string]*Object
)

func Python(types []Type) []byte {
	var buffer bytes.Buffer
	buffer.Write([]byte(`
import typy
`))

	PY_WRITER = &buffer
	PY_TYPES = make(map[string]bool)
	var codes []string
	for _, t := range types {
		codes = append(codes, t.Python())
	}
	PY_TYPES = nil
	PY_WRITER = nil

	buffer.Write([]byte(strings.Join(codes, "")))
	return buffer.Bytes()
}

func (t *Enum) Python() string {
	if ok, exist := PY_TYPES[t.Name]; exist && ok {
		return ""
	}

	PY_TYPES[t.Name] = true
	var enums []string
	for _, name := range t.Sorted() {
		enums = append(enums, fmt.Sprintf(`
	%s = %d, "%s"`, name, t.Values[name], name))
	}

	return fmt.Sprintf(`
class %s(typy.Enum):%s
`, t.Name, strings.Join(enums, ""))
}

func (t *Object) Python() string {
	if ok, exist := PY_TYPES[t.Name]; exist && ok {
		return ""
	}

	PY_TYPES[t.Name] = true
	var fields []string
	var sequences []string
	for _, f := range t.VisibleFields() {
		sequences = append(sequences, fmt.Sprintf("'%s'", f.Name))
		fields = append(fields, fmt.Sprintf(`
	%s = %s`, f.Name, strings.Replace(f.Python(), "typy.", "typy.pb.", 1)))
	}

	var sequence string
	if sequences == nil {
		return ""
	} else if len(sequences) > 1 {
		sequence = fmt.Sprintf(`
	____propertySequence__ = %s`, strings.Join(sequences, ", "))
	}

	parent := "typy.Object"
	if t.HasParent() {
		parent = t.Parent.Name
	}

	return fmt.Sprintf(`
class %s(%s):%s%s
`, t.Name, parent, sequence, strings.Join(fields, ""))
}

func (t UnknownType) Python() string {
	return ""
}

func (t SimpleType) Python() string {
	switch t {
	case SimpleType_INT32:
		fallthrough
	case SimpleType_INT64:
		fallthrough
	case SimpleType_UINT32:
		fallthrough
	case SimpleType_UINT64:
		return "typy.Integer"
	case SimpleType_FLOAT32:
		fallthrough
	case SimpleType_FLOAT64:
		return "typy.Float"
	case SimpleType_BYTES:
		return "typy.Bytes"
	case SimpleType_STRING:
		return "typy.String"
	case SimpleType_BOOL:
		return "typy.Boolean"
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for Python: %d", t)
		return "Unknown"
	}
}

func (t *EnumType) Python() string {
	PY_WRITER.Write([]byte(t.Enum.Python()))
	return fmt.Sprintf("typy.Enum(%s)", t.Name)
}

func (t *InstanceType) Python() string {
	if object, ok := TS_OBJECTS[t.Name]; ok {
		PY_WRITER.Write([]byte(object.Python()))
		return fmt.Sprintf("typy.Instance(%s)", t.Name)
	} else {
		return fmt.Sprintf("typy.Python(%s)", t.Name)
	}
}

func (t *FixedPointType) Python() string {
	return fmt.Sprintf("typy.FixedPoint(%d, %d)", t.Precision, t.Floor)
}

func (t *ListType) Python() string {
	return fmt.Sprintf("typy.List(%s)", t.E.Python())
}

func (t *DictType) Python() string {
	return fmt.Sprintf("typy.Dict(%s, %s)", t.K.Python(), t.V.Python())
}

func (t *VariantType) Python() string {
	var variants []string
	for _, v := range t.Ts {
		if inst, ok := v.(*InstanceType); ok {
			if object, ok := TS_OBJECTS[inst.Name]; ok {
				PY_WRITER.Write([]byte(object.Python()))
				variants = append(variants, inst.Name)
				continue
			}
		}
		variants = append(variants, v.Python())
	}
	return fmt.Sprintf("typy.Instance(%s)", strings.Join(variants, ", "))
}

const (
	TYPYD_ENUM       = 0
	TYPYD_INT32      = 1
	TYPYD_INT64      = 2
	TYPYD_UINT32     = 3
	TYPYD_UINT64     = 4
	TYPYD_FIXEDPOINT = 5
	TYPYD_DOUBLE     = 6
	TYPYD_FLOAT      = 7
	TYPYD_BOOL       = 8
	TYPYD_BYTES      = 9
	TYPYD_STRING     = 10
	TYPYD_OBJECT     = 11
	TYPYD_VARIANT    = 12
	TYPYD_LIST       = 13
	TYPYD_DICT       = 14
	TYPYD_PYTHON     = 15
	MAX_TYPYD_TYPE   = 16
)

func Typyd(types []Type) []byte {
	var buffer bytes.Buffer
	buffer.Write([]byte(`
from typy import _typyd
`))

	PY_WRITER = &buffer
	PY_OBJECTS = ObjectMap(types)
	PY_TYPES = make(map[string]bool)
	var codes []string
	for _, t := range types {
		codes = append(codes, t.Typyd())
	}
	PY_TYPES = nil
	PY_WRITER = nil
	PY_OBJECTS = nil

	buffer.Write([]byte(strings.Join(codes, "")))
	return buffer.Bytes()
}

func (t *Enum) Typyd() string {
	if ok, exist := PY_TYPES[t.Name]; exist && ok {
		return ""
	} else {
		PY_TYPES[t.Name] = true
		return fmt.Sprintf(`
%s = _typyd.Enum('%s')`, t.Name, t.Name)
	}
}

func (t *Object) Typyd() string {
	if ok, exist := PY_TYPES[t.Name]; exist && ok {
		return ""
	}

	PY_TYPES[t.Name] = true
	var fields []string
	for i, f := range t.AllFields(PY_OBJECTS, true) {
		wiretype := f.WireType()
		fields = append(fields, fmt.Sprintf(`
	('%s', %d, %d, %d, %s),`, f.Name, _MAKE_TAG(i+1, wiretype), TAG_SIZE(i+1), wiretype, f.Typyd()))
	}

	if len(fields) <= 0 {
		return ""
	}

	return fmt.Sprintf(`
%s = _typyd.Object('%s', (%s
))`, t.Name, t.Name, strings.Join(fields, ""))
}

func (t UnknownType) Typyd() string {
	return ""
}

func (t SimpleType) Typyd() string {
	switch t {
	case SimpleType_INT32:
		return strconv.Itoa(TYPYD_INT32)
	case SimpleType_INT64:
		return strconv.Itoa(TYPYD_INT64)
	case SimpleType_UINT32:
		return strconv.Itoa(TYPYD_UINT32)
	case SimpleType_UINT64:
		return strconv.Itoa(TYPYD_UINT64)
	case SimpleType_FLOAT32:
		return strconv.Itoa(TYPYD_FLOAT)
	case SimpleType_FLOAT64:
		return strconv.Itoa(TYPYD_DOUBLE)
	case SimpleType_BYTES:
		return strconv.Itoa(TYPYD_BYTES)
	case SimpleType_STRING:
		return strconv.Itoa(TYPYD_STRING)
	case SimpleType_BOOL:
		return strconv.Itoa(TYPYD_BOOL)
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for Typyd: %d", t)
		return "unknown"
	}
}

func (t *EnumType) Typyd() string {
	PY_WRITER.Write([]byte(t.Enum.Typyd()))
	return fmt.Sprintf("%d, %s", TYPYD_ENUM, t.Name)
}

func (t *InstanceType) Typyd() string {
	if object, ok := TS_OBJECTS[t.Name]; ok {
		PY_WRITER.Write([]byte(object.Typyd()))
		return fmt.Sprintf("%d, %s", TYPYD_OBJECT, t.Name)
	} else {
		if ok, exist := PY_TYPES[t.Name]; !exist || !ok {
			PY_WRITER.Write([]byte(fmt.Sprintf(`
%s = _typyd.Python('%s')`, t.Name, t.Name)))
			PY_TYPES[t.Name] = true
		}
		return fmt.Sprintf("%d, %s", TYPYD_PYTHON, t.Name)
	}
}

func (t *FixedPointType) Typyd() string {
	identifier := t.Identifier()
	if ok, exist := PY_TYPES[identifier]; !exist || !ok {
		PY_WRITER.Write([]byte(fmt.Sprintf(`
%s = _typyd.FixedPoint(%d, %d)`, identifier, t.Floor, t.Precision)))
		PY_TYPES[identifier] = true
	}
	return fmt.Sprintf("%d, %s", TYPYD_FIXEDPOINT, identifier)
}

func (t *ListType) Typyd() string {
	identifier := t.Identifier()
	if ok, exist := PY_TYPES[identifier]; !exist || !ok {
		PY_WRITER.Write([]byte(fmt.Sprintf(`
%s = _typyd.List('%s', (%d, %s))`, identifier, identifier, t.E.WireType(), t.E.Typyd())))
		PY_TYPES[identifier] = true
	}
	return fmt.Sprintf("%d, %s", TYPYD_LIST, identifier)
}

func (t *DictType) Typyd() string {
	identifier := t.Identifier()
	if ok, exist := PY_TYPES[identifier]; !exist || !ok {
		PY_WRITER.Write([]byte(fmt.Sprintf(`
%s = _typyd.Dict('%s', (%d, %s), (%d, %s))`, identifier, identifier,
			t.K.WireType(), t.K.Typyd(), t.V.WireType(), t.V.Typyd())))
		PY_TYPES[identifier] = true
	}
	return fmt.Sprintf("%d, %s", TYPYD_DICT, identifier)
}

func (t *VariantType) Typyd() string {
	identifier := t.Identifier()
	if ok, exist := PY_TYPES[identifier]; !exist || !ok {
		properties := make(map[string]Type)
		for _, typ := range t.Ts {
			switch v := typ.(type) {
			case SimpleType:
				switch v {
				case SimpleType_INT32:
					fallthrough
				case SimpleType_INT64:
					fallthrough
				case SimpleType_UINT32:
					fallthrough
				case SimpleType_UINT64:
					properties["Integer"] = v
				case SimpleType_FLOAT32:
					properties["Float"] = v
				case SimpleType_FLOAT64:
					properties["Double"] = v
				case SimpleType_BYTES:
					properties["Bytes"] = v
				case SimpleType_STRING:
					properties["String"] = v
				case SimpleType_BOOL:
					properties["Boolean"] = v
				default:
					log.Fatalf("[Tygo][VariantType] Unexpect enum value for SimpleType: %d", v)
				}
			case *EnumType:
				properties["Enum"] = v
			case *InstanceType:
				properties[v.Name] = v
			case *FixedPointType:
				properties["Dict"] = v
			case *ListType:
				properties["Dict"] = v
			case *DictType:
				properties["Dict"] = v
			default:
				log.Fatalf("[Tygo][VariantType] Unexpect type for Typyd: %v", v)
			}
		}

		var sortedProperties []string
		for n, _ := range properties {
			sortedProperties = append(sortedProperties, n)
		}
		sort.Strings(sortedProperties)

		var fields []string
		for i, n := range sortedProperties {
			wiretype := properties[n].WireType()
			fields = append(fields, fmt.Sprintf(`
	('%s', %d, %d, %d, %s),`, n, _MAKE_TAG(i+1, wiretype), TAG_SIZE(i+1), wiretype, properties[n].Typyd()))
		}

		PY_WRITER.Write([]byte(fmt.Sprintf(`
%s = _typyd.Variant('%s', (%s
))`, identifier, identifier, strings.Join(fields, ""))))
		PY_TYPES[identifier] = true
	}
	return fmt.Sprintf("%d, %s", TYPYD_VARIANT, identifier)
}
