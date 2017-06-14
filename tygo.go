// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

const TYGO_PATH = "github.com/ibelie/tygo"

type Tygo struct {
	cachedSize int
}

type Type interface {
	String() string
	IsPrimitive() bool
	Go() (string, map[string]string)
	ByteSizeGo(string, string, string) (string, map[string]string)
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
}

type Method struct {
	Name    string
	Params  []Type
	Results []Type
}

type Object struct {
	Name    string
	Parent  Type
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
	SimpleType_BOOL
	SimpleType_FLOAT32
	SimpleType_FLOAT64
)

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

type FixedPointType struct {
	Precision uint
	Floor     int
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

func (t *Method) IsPrimitive() bool {
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

func (t *EnumType) IsPrimitive() bool {
	return true
}

func (t *InstanceType) IsPrimitive() bool {
	return false
}

func (t *FixedPointType) IsPrimitive() bool {
	return true
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
		values = append(values, fmt.Sprintf(`%s: %d`, name, t.Values[name]))
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

func (t *Object) String() string {
	var fields []string

	fields = append(fields, fmt.Sprintf(`
	%s`, t.Parent))

	nameMax := 0
	for _, field := range t.Fields {
		if nameMax < len(field.Name) {
			nameMax = len(field.Name)
		}
	}
	for _, field := range t.Fields {
		fields = append(fields, fmt.Sprintf(`
	%s %s%s`, field.Name, strings.Repeat(" ", nameMax-len(field.Name)), field))
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
	case "nil":
		return SimpleType_NIL
	case "int32":
		return SimpleType_INT32
	case "int64":
		return SimpleType_INT64
	case "uint32":
		return SimpleType_UINT32
	case "uint64":
		return SimpleType_UINT64
	case "bytes":
		return SimpleType_BYTES
	case "string":
		return SimpleType_STRING
	case "bool":
		return SimpleType_BOOL
	case "float32":
		return SimpleType_FLOAT32
	case "float64":
		return SimpleType_FLOAT64
	default:
		return UnknownType(s)
	}
}

func (t SimpleType) String() string {
	switch t {
	case SimpleType_NIL:
		return "nil"
	case SimpleType_INT32:
		return "int32"
	case SimpleType_INT64:
		return "int64"
	case SimpleType_UINT32:
		return "uint32"
	case SimpleType_UINT64:
		return "uint64"
	case SimpleType_BYTES:
		return "bytes"
	case SimpleType_STRING:
		return "string"
	case SimpleType_BOOL:
		return "bool"
	case SimpleType_FLOAT32:
		return "float32"
	case SimpleType_FLOAT64:
		return "float64"
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value: %d", t)
		return "unknown"
	}
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

func (t *FixedPointType) String() string {
	return fmt.Sprintf("fixedpoint<%d, %d>", t.Precision, t.Floor)
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
