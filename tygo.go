// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
)

var TYGO_PATH = reflect.TypeOf(Tygo{}).PkgPath()

type Tygo struct {
	cachedSize int
}

func (t *Tygo) SetCachedSize(cachedSize int) {
	t.cachedSize = cachedSize
}

func (t *Tygo) CachedSize() int {
	return t.cachedSize
}

type Type interface {
	WireType() WireType
	Identifier() string
	String() string
	IsPrimitive() bool
	IsIterative() bool
	Go() (string, map[string]string)
	Typyd() string
	Python() string
	Typescript() string
	Javascript() (string, map[string]string)
	ByteSizeGo(string, string, string, int, bool) (string, map[string]string)
	CachedSizeGo(string, string, string, int, bool) (string, map[string]string)
	SerializeGo(string, string, string, int, bool) (string, map[string]string)
	DeserializeGo(string, string, string, string, int, bool) (string, WireType, map[string]string)
}

type Enum struct {
	nameMax int
	sorted  []string
	Name    string
	Values  map[string]int
}

type Field struct {
	Type
	Name string
	Hide bool
}

type Method struct {
	Name    string
	Params  []Type
	Results []Type
}

type Object struct {
	Name    string
	Parent  *InstanceType
	Fields  []*Field
	Methods []*Method
}

type UnknownType string

type SimpleType uint

const (
	SimpleType_UNKNOWN SimpleType = iota
	SimpleType_NIL
	SimpleType_INT32
	SimpleType_INT64
	SimpleType_UINT32
	SimpleType_UINT64
	SimpleType_BYTES
	SimpleType_STRING
	SimpleType_SYMBOL
	SimpleType_BOOL
	SimpleType_FLOAT32
	SimpleType_FLOAT64
)

type FixedPointType struct {
	Precision uint
	Floor     int
}

type EnumType struct {
	*Enum
	Name string
}

type InstanceType struct {
	*Object
	IsPtr   bool
	Name    string
	PkgName string
	PkgPath string
}

type ListType struct {
	E Type
}

type DictType struct {
	K Type
	V Type
}

type VariantType struct {
	Ts []Type
}

func (t *Enum) IsPrimitive() bool {
	return false
}

func (t *Object) IsPrimitive() bool {
	return false
}

func (t UnknownType) IsPrimitive() bool {
	return false
}

func (t SimpleType) IsPrimitive() bool {
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
		fallthrough
	case SimpleType_FLOAT32:
		fallthrough
	case SimpleType_FLOAT64:
		return true
	default:
		return false
	}
}

func (t *FixedPointType) IsPrimitive() bool {
	return true
}

func (t *EnumType) IsPrimitive() bool {
	return true
}

func (t *InstanceType) IsPrimitive() bool {
	return false
}

func (t *ListType) IsPrimitive() bool {
	return false
}

func (t *DictType) IsPrimitive() bool {
	return false
}

func (t *VariantType) IsPrimitive() bool {
	return false
}

func (t *Enum) IsIterative() bool {
	return false
}

func (t *Object) IsIterative() bool {
	return false
}

func (t UnknownType) IsIterative() bool {
	return false
}

func (t SimpleType) IsIterative() bool {
	return false
}

func (t *FixedPointType) IsIterative() bool {
	return false
}

func (t *EnumType) IsIterative() bool {
	return false
}

func (t *InstanceType) IsIterative() bool {
	return false
}

func (t *ListType) IsIterative() bool {
	return !t.E.IsPrimitive()
}

func (t *DictType) IsIterative() bool {
	return true
}

func (t *VariantType) IsIterative() bool {
	return false
}

func (t *Enum) Sorted() []string {
	if t.sorted == nil || len(t.sorted) != len(t.Values) {
		t.sorted = nil
		for name, _ := range t.Values {
			if t.nameMax < len(name) {
				t.nameMax = len(name)
			}
			t.sorted = append(t.sorted, name)
		}
	} else if t.nameMax <= 0 {
		for name, _ := range t.Values {
			if t.nameMax < len(name) {
				t.nameMax = len(name)
			}
		}
	}
	if !sort.IsSorted(t) {
		sort.Sort(t)
	}
	return t.sorted
}

func (t *Enum) Len() int {
	return len(t.sorted)
}

func (t *Enum) Swap(i, j int) {
	t.sorted[i], t.sorted[j] = t.sorted[j], t.sorted[i]
}

func (t *Enum) Less(i, j int) bool {
	return t.Values[t.sorted[i]] < t.Values[t.sorted[j]]
}

func (t *Enum) String() string {
	var values []string
	for _, name := range t.Sorted() {
		values = append(values, fmt.Sprintf("%s: %d", name, t.Values[name]))
	}
	return fmt.Sprintf("%s[%s]", t.Name, strings.Join(values, ", "))
}

func (t *Method) String() string {
	var params []string
	for _, param := range t.Params {
		params = append(params, param.String())
	}
	s := fmt.Sprintf("%s(%s)", t.Name, strings.Join(params, ", "))

	var results []string
	for _, result := range t.Results {
		results = append(results, result.String())
	}
	if len(results) == 1 {
		s = fmt.Sprintf("%s %s", s, results[0])
	} else if len(results) > 1 {
		s = fmt.Sprintf("%s (%s)", s, strings.Join(results, ", "))
	}

	return s
}

func (t *Object) HasParent() bool {
	if t.Parent.Name == "Tygo" && t.Parent.PkgName == "tygo" && t.Parent.PkgPath == TYGO_PATH {
		return false
	} else if t.Parent.PkgPath != "" {
		return true
	} else {
		name, num := t.Parent.Object.MaxFieldNum()
		return name != "" || num != 0
	}
}

func (t *Object) MaxFieldNum() (string, int) {
	if !t.HasParent() {
		return "", len(t.VisibleFields())
	} else if t.Parent.PkgPath == "" {
		name, num := t.Parent.Object.MaxFieldNum()
		return name, num + len(t.VisibleFields())
	} else {
		return t.Parent.Name, len(t.VisibleFields())
	}
}

func (t *Object) VisibleFields() (fields []*Field) {
	for _, field := range t.Fields {
		if !field.Hide {
			fields = append(fields, field)
		}
	}
	return
}

func (t *Object) AllFields(objects map[string]*Object, visible bool) []*Field {
	fields := t.Fields
	if visible {
		fields = t.VisibleFields()
	}
	if !t.HasParent() {
		return fields
	} else if t.Parent.PkgPath == "" {
		return append(t.Parent.Object.AllFields(objects, visible), fields...)
	} else if object, ok := objects[t.Parent.Name]; ok {
		return append(object.AllFields(objects, visible), fields...)
	} else {
		log.Fatalf("[Tygo][Object] Unsolved parent type: %s", t.Parent.Name)
		return fields
	}
}

func (t *Object) String() string {
	var fields []string

	fields = append(fields, fmt.Sprintf(`
	%s`, t.Parent))

	nameMax := 0
	for _, field := range t.Fields {
		length := len(field.Name)
		if field.Hide {
			length++
		}
		if nameMax < length {
			nameMax = length
		}
	}
	for _, field := range t.Fields {
		name := field.Name
		if field.Hide {
			name = "." + field.Name
		}
		fields = append(fields, fmt.Sprintf(`
	%s %s%s`, name, strings.Repeat(" ", nameMax-len(name)), field))
	}
	for _, method := range t.Methods {
		fields = append(fields, fmt.Sprintf(`
	%s`, method))
	}

	return fmt.Sprintf(`
%s{%s
}
`, t.Name, strings.Join(fields, ""))
}

func (t UnknownType) String() string {
	return string(t)
}

func SimpleType_FromString(s string) Type {
	switch s {
	case STR_NIL:
		return SimpleType_NIL
	case STR_INT32:
		return SimpleType_INT32
	case STR_INT64:
		return SimpleType_INT64
	case STR_UINT32:
		return SimpleType_UINT32
	case STR_UINT64:
		return SimpleType_UINT64
	case STR_BYTES:
		return SimpleType_BYTES
	case STR_STRING:
		return SimpleType_STRING
	case STR_SYMBOL:
		return SimpleType_SYMBOL
	case STR_BOOL:
		return SimpleType_BOOL
	case STR_FLOAT32:
		return SimpleType_FLOAT32
	case STR_FLOAT64:
		return SimpleType_FLOAT64
	default:
		return UnknownType(s)
	}
}

func (t SimpleType) String() string {
	switch t {
	case SimpleType_NIL:
		return STR_NIL
	case SimpleType_INT32:
		return STR_INT32
	case SimpleType_INT64:
		return STR_INT64
	case SimpleType_UINT32:
		return STR_UINT32
	case SimpleType_UINT64:
		return STR_UINT64
	case SimpleType_BYTES:
		return STR_BYTES
	case SimpleType_STRING:
		return STR_STRING
	case SimpleType_SYMBOL:
		return STR_SYMBOL
	case SimpleType_BOOL:
		return STR_BOOL
	case SimpleType_FLOAT32:
		return STR_FLOAT32
	case SimpleType_FLOAT64:
		return STR_FLOAT64
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value: %d", t)
		return STR_UNKNOWN
	}
}

func (t *FixedPointType) ToVarint(name string) string {
	value := name
	if t.Floor > 0 {
		value = fmt.Sprintf("(%s - %d)", value, t.Floor)
	} else if t.Floor < 0 {
		value = fmt.Sprintf("(%s + %d)", value, -t.Floor)
	}
	if precision := pow10(t.Precision); precision > 0 {
		value = fmt.Sprintf("%s * %d", value, precision)
	}
	return value
}

func (t *FixedPointType) FromVarint(name string) string {
	value := fmt.Sprintf("float64(%s)", name)
	if precision := pow10(t.Precision); precision > 0 {
		value = fmt.Sprintf("%s / %d", value, precision)
	}
	if t.Floor > 0 {
		value = fmt.Sprintf("%s + %d", value, t.Floor)
	} else if t.Floor < 0 {
		value = fmt.Sprintf("%s - %d", value, -t.Floor)
	}
	return value
}

func (t *FixedPointType) String() string {
	return fmt.Sprintf("fixedpoint<%d, %d>", t.Precision, t.Floor)
}

func (t *EnumType) String() string {
	return t.Name
}

func (t *InstanceType) String() string {
	s := ""
	if t.IsPtr {
		s += "*"
	}
	if t.PkgPath != "" {
		s += t.PkgPath + "."
	}
	s += t.Name
	return s
}

func (t *ListType) String() string {
	return fmt.Sprintf("[]%s", t.E)
}

func (t *DictType) String() string {
	return fmt.Sprintf("map[%s]%s", t.K, t.V)
}

func (t *VariantType) String() string {
	var ts []string
	for _, t := range t.Ts {
		ts = append(ts, t.String())
	}
	return fmt.Sprintf("variant<%s>", strings.Join(ts, ", "))
}

func (t *Enum) Identifier() string {
	return t.Name
}

func (t *Object) Identifier() string {
	return t.Name
}

func (t UnknownType) Identifier() string {
	return string(t)
}

func (t SimpleType) Identifier() string {
	switch t {
	case SimpleType_NIL:
		return "n"
	case SimpleType_INT32:
		return "i"
	case SimpleType_INT64:
		return "i2"
	case SimpleType_UINT32:
		return "ui"
	case SimpleType_UINT64:
		return "ui2"
	case SimpleType_BYTES:
		return "by"
	case SimpleType_STRING:
		return "s"
	case SimpleType_SYMBOL:
		return "sy"
	case SimpleType_BOOL:
		return "b"
	case SimpleType_FLOAT32:
		return "f"
	case SimpleType_FLOAT64:
		return "d"
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for Identifier: %d", t)
		return ""
	}
}

func (t *EnumType) Identifier() string {
	return t.Name
}

func (t *InstanceType) Identifier() string {
	return t.Name
}

func (t *FixedPointType) Identifier() string {
	if t.Floor >= 0 {
		return fmt.Sprintf("FixedPoint_%d_%d", t.Precision, t.Floor)
	} else {
		return fmt.Sprintf("FixedPoint_%d_s%d", t.Precision, -t.Floor)
	}
}

func (t *ListType) Identifier() string {
	return "L" + t.E.Identifier()
}

func (t *DictType) Identifier() string {
	return "D" + t.K.Identifier() + t.V.Identifier()
}

func (t *VariantType) Identifier() string {
	var names []string
	for _, v := range t.Ts {
		names = append(names, v.Identifier())
	}
	sort.Strings(names)
	return "V" + shortName(strings.Join(names, ""))
}
