// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"strings"

	"io/ioutil"
)

var TS_OBJECTS map[string]*Object

func Typescript(dir string, name string, module string, types []Type, propPre []Type) {
	var buffer bytes.Buffer

	TS_OBJECTS = make(map[string]*Object)
	for _, t := range types {
		if object, ok := t.(*Object); ok {
			if o, exist := TS_OBJECTS[object.Name]; exist {
				log.Fatalf("[Tygo][Typescript] Object already exists: %v %v", o, object)
			}
			TS_OBJECTS[object.Name] = object
		}
	}

	PROP_PRE = propPre
	var codes []string
	for _, t := range types {
		codes = append(codes, t.Typescript())
	}
	PROP_PRE = nil
	TS_OBJECTS = nil

	buffer.Write([]byte(fmt.Sprintf(`// Generated for tyts by tygo.  DO NOT EDIT!

declare module %s {
	interface Type {
		__class__: string;
		ByteSize(): number;
		Serialize(): Uint8Array;
		Deserialize(data: Uint8Array): void;
	}%s
}
`, module, strings.Join(codes, ""))))

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
		static Serialize%s%s(%s): Uint8Array;
		static Deserialize%s%s(data: Uint8Array): any;`, name, typ, strings.Join(items, ", "), name, typ)
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
		for _, field := range t.Fields {
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
		__class__: string;
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
	if _, ok := TS_OBJECTS[t.Name]; ok {
		return t.Name
	} else {
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
