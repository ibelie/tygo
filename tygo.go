// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"sort"
	"strings"
)

const TYGO_PATH = "github.com/ibelie/tygo"

type Tygo struct {
	cachedSize int
}

type Type interface {
	String() string
	Go() (string, map[string]string)
	BsGo(string, string) (string, map[string]string)
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

type SimpleType string

type ObjectType struct {
	IsPtr   bool
	Name    string
	PkgName string
	PkgPath string
}

type FixedPointType struct {
	Precision int
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

func (t SimpleType) String() string {
	return string(t)
}

func (t *ObjectType) String() string {
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

func update(a map[string]string, b map[string]string) map[string]string {
	if b == nil {
		return a
	} else if a == nil {
		return b
	}
	for k, v := range b {
		a[k] = v
	}
	return a
}
