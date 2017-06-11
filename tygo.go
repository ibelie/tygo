// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"strings"
)

type Type interface {
	String() string
	Go() (string, []string)
}

type Enum struct {
	Name   string
	Values map[string]int
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

func (t SimpleType) Go() (string, []string) {
	return string(t), nil
}

type ObjectType struct {
	IsPtr bool
	Pkg   string
	Name  string
}

func (t *ObjectType) String() string {
	s := ""
	if t.IsPtr {
		s += "*"
	}
	if t.Pkg != "" {
		s += t.Pkg + "."
	}
	s += t.Name
	return s
}

func (t *ObjectType) Go() (string, []string) {
	return t.String(), []string{t.Pkg}
}

type FixedPointType struct {
	Precision int
	Floor     int
}

func (t *FixedPointType) String() string {
	return fmt.Sprintf("fixedpoint<%d, %d>", t.Precision, t.Floor)
}

func (t *FixedPointType) Go() (string, []string) {
	return "float64", nil
}

type ListType struct {
	E Type
}

func (t *ListType) String() string {
	return fmt.Sprintf("[]%s", t.E)
}

func (t *ListType) Go() (string, []string) {
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

func (t *DictType) Go() (string, []string) {
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

func (t *VariantType) Go() (string, []string) {
	var p []string
	for _, vt := range t.Ts {
		_, vp := vt.Go()
		p = append(p, vp...)
	}
	return "interface{}", p
}
