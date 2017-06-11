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
}

type Enum struct {
	Values map[string]int
}

type Method struct {
	Params  []Type
	Results []Type
}

type Object struct {
	Fields  map[string]Type
	Methods map[string]*Method
}

type SimpleType string

func (t SimpleType) String() string {
	return string(t)
}

type ObjectType struct {
	T     Type
	isPtr bool
	pkg   string
}

func (t *ObjectType) String() string {
	s := ""
	if t.isPtr {
		s += "*"
	}
	if t.pkg != "" {
		s += t.pkg + "."
	}
	s += t.T.String()
	return s
}

type FixedPointType struct {
	precision int
	floor     int
}

func (t *FixedPointType) String() string {
	return fmt.Sprintf("fixedpoint<%d, %d>", t.precision, t.floor)
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
	return fmt.Sprintf("variant<%s>", strings.Join(t.Ts, ", "))
}
