// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"sort"
	"strings"
)

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

func (e *Enum) NameMax() int {
	if e.nameMax <= 0 {
		for name, _ := range e.Values {
			if e.nameMax < len(name) {
				e.nameMax = len(name)
			}
		}
	}
	return e.nameMax
}

func (e *Enum) Sorted() []string {
	if e.sorted != nil && len(e.sorted) == len(e.Values) && sort.IsSorted(e) {
		return e.sorted
	}
	e.sorted = nil
	e.nameMax = 0
	for name, _ := range e.Values {
		if e.nameMax < len(name) {
			e.nameMax = len(name)
		}
		e.sorted = append(e.sorted, name)
	}
	sort.Sort(e)
	return e.sorted
}

func (e *Enum) Len() int {
	return len(e.sorted)
}

func (e *Enum) Swap(i, j int) {
	e.sorted[i], e.sorted[j] = e.sorted[j], e.sorted[i]
}

func (e *Enum) Less(i, j int) bool {
	return e.Values[e.sorted[i]] < e.Values[e.sorted[j]]
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

func (t SimpleType) String() string {
	return string(t)
}

func (t SimpleType) Go() (string, [][2]string) {
	return string(t), nil
}

type ObjectType struct {
	IsPtr   bool
	Name    string
	PkgName string
	PkgPath string
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

func (t *ObjectType) Go() (string, [][2]string) {
	if t.PkgPath == "" {
		return t.String(), nil
	} else {
		s := ""
		if t.IsPtr {
			s += "*"
		}
		s += t.PkgName + "." + t.Name
		p := strings.Split(t.PkgPath, "/")
		var a string
		if t.PkgName == p[len(p)-1] {
			a = ""
		} else {
			a = t.PkgName + " "
		}
		return s, [][2]string{[2]string{a, t.PkgPath}}
	}
}

type FixedPointType struct {
	Precision int
	Floor     int
}

func (t *FixedPointType) String() string {
	return fmt.Sprintf("fixedpoint<%d, %d>", t.Precision, t.Floor)
}

func (t *FixedPointType) Go() (string, [][2]string) {
	return "float64", nil
}

type ListType struct {
	E Type
}

func (t *ListType) String() string {
	return fmt.Sprintf("[]%s", t.E)
}

func (t *ListType) Go() (string, [][2]string) {
	s, p := t.E.Go()
	return fmt.Sprintf("[]%s", s), p
}

type DictType struct {
	K Type
	V Type
}

func (t *DictType) String() string {
	return fmt.Sprintf("map[%s]%s", t.K, t.V)
}

func (t *DictType) Go() (string, [][2]string) {
	ks, kp := t.K.Go()
	vs, vp := t.V.Go()
	return fmt.Sprintf("map[%s]%s", ks, vs), append(kp, vp...)
}

type VariantType struct {
	Ts []Type
}

func (t *VariantType) String() string {
	var ts []string
	for _, t := range t.Ts {
		ts = append(ts, t.String())
	}
	return fmt.Sprintf("variant<%s>", strings.Join(ts, ", "))
}

func (t *VariantType) Go() (string, [][2]string) {
	var p [][2]string
	for _, vt := range t.Ts {
		_, vp := vt.Go()
		p = append(p, vp...)
	}
	return "interface{}", p
}
