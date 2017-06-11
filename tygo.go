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

type FixedPointType struct {
	Precision int
	Floor     int
}

func (t *FixedPointType) String() string {
	return fmt.Sprintf("fixedpoint<%d, %d>", t.Precision, t.Floor)
}

type ListType struct {
	E Type
}

func (t *ListType) String() string {
	return fmt.Sprintf("[]%s", t.E)
}

type DictType struct {
	K Type
	V Type
}

func (t *DictType) String() string {
	return fmt.Sprintf("map[%s]%s", t.K, t.V)
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
