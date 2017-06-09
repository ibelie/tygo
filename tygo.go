// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"sort"
)

type FixedPoint float64

type Type interface {
	String() string
	TypeStr() *TypeStr
}

type Message map[string]Type

func (message Message) sorted() (fields []string) {
	for field := range message {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return
}

type SimpleType string

func (t SimpleType) String() string {
	return string(t)
}

func (t SimpleType) TypeStr() *TypeStr {
	return &TypeStr{Simple: t.String()}
}

type ListType struct {
	E Type
}

func (t *ListType) String() string {
	return "List(" + t.E.String() + ")"
}

func (t *ListType) TypeStr() *TypeStr {
	return &TypeStr{List: t.E.TypeStr()}
}

type DictType struct {
	K Type
	V Type
}

func (t *DictType) String() string {
	return "Dict(" + t.K.String() + ", " + t.V.String() + ")"
}

func (t *DictType) TypeStr() *TypeStr {
	return &TypeStr{Key: t.K.TypeStr(), Value: t.V.TypeStr()}
}
