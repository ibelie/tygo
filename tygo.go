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
	Go() (string, [][2]string)
}

type Enum struct {
	nameMax int
	sorted  []string
	Name    string
	Values  map[string]int
}

type Method struct {
	Name    string
	Params  []Type
	Results []Type
}

type Object struct {
	Name    string
	Fields  map[string]Type
	Parents []Type
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

func (t *Object) String() string {
	var fields []string

	var sortedParent []string
	for _, parent := range t.Parents {
		sortedParent = append(sortedParent, parent.String())
	}
	sort.Strings(sortedParent)
	for _, parent := range sortedParent {
		fields = append(fields, fmt.Sprintf(`
	%s`, parent))
	}

	nameMax := 0
	var sortedField []string
	for name, _ := range t.Fields {
		if nameMax < len(name) {
			nameMax = len(name)
		}
		sortedField = append(sortedField, name)
	}
	sort.Strings(sortedField)
	for _, name := range sortedField {
		fields = append(fields, fmt.Sprintf(`
	%s %s%s`, name, strings.Repeat(" ", nameMax-len(name)), t.Fields[name].String()))
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
