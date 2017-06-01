// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"

	"encoding/json"
	"io/ioutil"
)

type TypeStr struct {
	Simple string
	List   *TypeStr
	Key    *TypeStr
	Value  *TypeStr
}

type ParentStr struct {
	Decorators
	Tag  string
	Type *TypeStr
}

type FieldStr struct {
	Decorators
	Name string
	Tag  string
	Type *TypeStr
}

type MethodStr struct {
	Decorators
	Name   string
	Params []*TypeStr
	Result []*TypeStr
}

type ObjectStr struct {
	Name    string
	Path    string
	Parents map[string]*ParentStr
	Fields  map[string]*FieldStr
	Methods map[string]*MethodStr
}

type PackageStr struct {
	Name    string
	Path    string
	Imports []string
	Objects map[string]*ObjectStr
}

func Load(path string) *Package {
	pkgStr := &PackageStr{}
	bytes, err := ioutil.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(bytes, pkgStr)
	}
	if err != nil {
		panic(fmt.Sprintf("[Tygo][Loader] Cannot read proto:\n>>>>%v", err))
	}

	pkg := &Package{
		Name:    pkgStr.Name,
		Path:    pkgStr.Path,
		Imports: pkgStr.Imports,
		Objects: make(map[string]*Object),
	}
	for _, object := range pkgStr.Objects {
		o := &Object{
			Name:    object.Name,
			Path:    object.Path,
			Parents: make(Parents),
			Fields:  make(Fields),
			Methods: make(Methods),
		}
		for _, parent := range object.Parents {
			o.Parents[parent.Type.Type().String()] = &Parent{
				Tag:        parent.Tag,
				Type:       parent.Type.Type(),
				Decorators: copyDecorators(parent.Decorators),
			}
		}
		for _, field := range object.Fields {
			o.Fields[field.Name] = &Field{
				Name:       field.Name,
				Tag:        field.Tag,
				Type:       field.Type.Type(),
				Decorators: copyDecorators(field.Decorators),
			}
		}
		for _, method := range object.Methods {
			m := &Method{
				Name:       method.Name,
				Decorators: copyDecorators(method.Decorators),
			}
			for _, param := range method.Params {
				m.Params = append(m.Params, param.Type())
			}
			for _, result := range method.Result {
				m.Result = append(m.Result, result.Type())
			}
			o.Methods[method.Name] = m
		}
		pkg.Objects[object.Name] = o
	}
	return pkg
}

func (pkg *Package) Save(path string) {
	pkgStr := &PackageStr{
		Name:    pkg.Name,
		Path:    pkg.Path,
		Imports: pkg.Imports,
		Objects: make(map[string]*ObjectStr),
	}
	for _, object := range pkg.Objects {
		o := &ObjectStr{
			Name:    object.Name,
			Path:    object.Path,
			Parents: make(map[string]*ParentStr),
			Fields:  make(map[string]*FieldStr),
			Methods: make(map[string]*MethodStr),
		}
		for _, parent := range object.Parents {
			o.Parents[parent.String()] = &ParentStr{
				Tag:        parent.Tag,
				Type:       parent.TypeStr(),
				Decorators: copyDecorators(parent.Decorators),
			}
		}
		for _, field := range object.Fields {
			o.Fields[field.Name] = &FieldStr{
				Name:       field.Name,
				Tag:        field.Tag,
				Type:       field.TypeStr(),
				Decorators: copyDecorators(field.Decorators),
			}
		}
		for _, method := range object.Methods {
			m := &MethodStr{
				Name:       method.Name,
				Decorators: copyDecorators(method.Decorators),
			}
			for _, param := range method.Params {
				m.Params = append(m.Params, param.TypeStr())
			}
			for _, result := range method.Result {
				m.Result = append(m.Result, result.TypeStr())
			}
			o.Methods[method.Name] = m
		}
		pkgStr.Objects[object.Name] = o
	}

	bytes, err := json.Marshal(pkgStr)
	if err == nil {
		err = ioutil.WriteFile(path, bytes, 0666)
	}
	if err != nil {
		panic(fmt.Sprintf("[Tygo][Loader] Cannot write proto:\n>>>>%v", err))
	}
}

func (t *TypeStr) Type() Type {
	if t.Key != nil && t.Value != nil {
		return &DictType{K: t.Key.Type(), V: t.Value.Type()}
	} else if t.List != nil {
		return &ListType{E: t.List.Type()}
	} else {
		return SimpleType(t.Simple)
	}
}
