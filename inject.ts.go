// Copyright 2017 - 2018 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"io/ioutil"
)

var (
	TS_MODULE     string
	TS_CUR_MODULE string
	TS_EX_TYPE    bool
	TS_OBJECTS    map[string]*Object
	EXTENS_PKG    map[string]string
)

func Typescript(dir string, name string, module string, types []Type, propPre []Type) {
	var buffer bytes.Buffer

	PROP_PRE = propPre
	TS_MODULE = module
	TS_OBJECTS = ObjectMap(types, TS_MODULE == "")

	var pkgTypes map[string][]Type
	var sortedPkgs []string
	if TS_MODULE == "" {
		pkgTypes = PkgTypeMap(types)
		for pkg, _ := range pkgTypes {
			sortedPkgs = append(sortedPkgs, pkg)
		}
		sort.Strings(sortedPkgs)
	} else {
		pkgTypes = map[string][]Type{module: types}
		sortedPkgs = []string{module}
	}

	var modules []string
	for _, pkg := range sortedPkgs {
		ts := pkgTypes[pkg]
		TS_CUR_MODULE = pkg
		TS_EX_TYPE = false

		var codes []string
		for _, t := range ts {
			codes = append(codes, t.Typescript())
		}

		exType := ""
		if TS_EX_TYPE {
			exType = `
	interface Type {
		isObject: true;
		ByteSize(): number;
		Serialize(): Uint8Array;
		Deserialize(data: Uint8Array): void;
	}`
		}

		modules = append(modules, fmt.Sprintf(`
declare module %s {%s%s
}
`, strings.Replace(pkg, "/", ".", -1), exType, strings.Join(codes, "")))
	}

	PROP_PRE = nil
	TS_OBJECTS = nil
	TS_EX_TYPE = false
	TS_MODULE = ""
	TS_CUR_MODULE = ""

	buffer.Write([]byte(fmt.Sprintf(`// Generated for tyts by tygo.  DO NOT EDIT!
%s`, strings.Join(modules, ""))))

	if name == "" {
		name = module
	}
	ioutil.WriteFile(path.Join(dir, name+".d.ts"), buffer.Bytes(), 0666)
	Javascript(dir, name, module, types, propPre)
}

func (t *Enum) Typescript() string {
	var enums []string
	for _, name := range t.Sorted() {
		enums = append(enums, fmt.Sprintf(`
		%s = %d`, name, t.Values[name]))
	}
	return fmt.Sprintf(`

	const enum %s {%s
	}`, t.Name, strings.Join(enums, ","))
}

func typeListTypescript(name string, typ string, ts []Type) string {
	var items []string
	for i, t := range ts {
		items = append(items, fmt.Sprintf("a%d: %s", i, t.Typescript()))
	}
	return fmt.Sprintf(`
		static S_%s%s(%s): Uint8Array;
		static D_%s%s(data: Uint8Array): any;`, name, typ, strings.Join(items, ", "), name, typ)
}

func (t *Object) Typescript() string {
	var parent string
	if t.HasParent() {
		parent = fmt.Sprintf(" extends %s", t.Parent.Typescript())
	}
	var members []string
	for _, field := range t.Fields {
		members = append(members, fmt.Sprintf(`
		%s: %s;`, field.Name, field.Typescript()))
	}

	if PROP_PRE != nil {
		for _, field := range t.VisibleFields() {
			members = append(members, typeListTypescript(field.Name, "", []Type{field}))
		}
	}

	for _, method := range t.Methods {
		if len(method.Params) > 0 {
			members = append(members, typeListTypescript(method.Name, "Param", method.Params))
		}
		if len(method.Results) > 0 {
			members = append(members, typeListTypescript(method.Name, "Result", method.Results))
		}
	}

	return fmt.Sprintf(`

	class %s%s {
		isObject: true;
		ByteSize(): number;
		Serialize(): Uint8Array;
		Deserialize(data: Uint8Array): void;
%s
	}

	namespace %s {
		function Deserialize(data: Uint8Array): %s;
	}`, t.Name, parent, strings.Join(members, ""), t.Name, t.Name)
}

func (t UnknownType) Typescript() string {
	return ""
}

func (t SimpleType) Typescript() string {
	switch t {
	case SimpleType_INT32:
		fallthrough
	case SimpleType_INT64:
		fallthrough
	case SimpleType_UINT32:
		fallthrough
	case SimpleType_UINT64:
		fallthrough
	case SimpleType_FLOAT32:
		fallthrough
	case SimpleType_FLOAT64:
		return "number"
	case SimpleType_BYTES:
		return "Uint8Array"
	case SimpleType_SYMBOL:
		fallthrough
	case SimpleType_STRING:
		return "string"
	case SimpleType_BOOL:
		return "boolean"
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for Typescript: %d", t)
		return "unknown"
	}
}

func (t *EnumType) Typescript() string {
	return t.Name
}

func (t *InstanceType) Typescript() string {
	fullName := t.Name
	if TS_MODULE == "" && t.Object != nil {
		fullName = t.Object.FullName()
	} else if TS_MODULE == "" && t.PkgPath != "" {
		fullName = t.PkgPath + "/" + t.Name
	} else if EXTENS_PKG != nil {
		if pkg, ok := EXTENS_PKG[t.Name]; ok {
			if TS_CUR_MODULE == pkg {
				return t.Name
			} else {
				return strings.Replace(pkg, "/", ".", -1) + "." + t.Name
			}
		}
	}
	if _, ok := TS_OBJECTS[fullName]; ok {
		if TS_CUR_MODULE == t.PkgPath || t.Object != nil {
			return t.Name
		}
		return strings.Replace(fullName, "/", ".", -1)
	} else {
		TS_EX_TYPE = true
		return "Type"
	}
}

func (t *FixedPointType) Typescript() string {
	return "number"
}

func (t *ListType) Typescript() string {
	return fmt.Sprintf("%s[]", t.E.Typescript())
}

func (t *DictType) Typescript() string {
	return fmt.Sprintf("{[index: %s]: %s}", t.K.Typescript(), t.V.Typescript())
}

func (t *VariantType) Typescript() string {
	return "any"
}
