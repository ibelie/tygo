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

func Typescript(dir string, name string, module string, types []Type, methodPre []Type) {
	var buffer bytes.Buffer

	objects := make(map[string]*Object)
	for _, t := range types {
		if object, ok := t.(*Object); ok {
			if o, exist := objects[object.Name]; exist {
				log.Fatalf("[Tygo][Typescript] Object already exists: %v %v", o, object)
			}
			objects[object.Name] = object
		}
	}

	METH_PRE = methodPre
	var codes []string
	for _, t := range types {
		codes = append(codes, t.Typescript(objects))
	}
	METH_PRE = nil

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
	Javascript(dir, name, module, types, methodPre)
}

func (t *Enum) Typescript(objects map[string]*Object) string {
	var enums []string
	for _, name := range t.Sorted() {
		enums = append(enums, fmt.Sprintf(`
		%s = %d`, name, t.Values[name]))
	}
	return fmt.Sprintf(`

	const enum %s {%s
	}`, t.Name, strings.Join(enums, ","))
}

func typeListTypescript(name string, typ string, ts []Type, objects map[string]*Object) string {
	var items []string
	for i, t := range ts {
		items = append(items, fmt.Sprintf("a%d: %s", i, t.Typescript(objects)))
	}
	return fmt.Sprintf(`
		static Serialize%s%s(%s): Uint8Array;
		static Deserialize%s%s(data: Uint8Array): any;`, name, typ, strings.Join(items, ", "), name, typ)
}

func (t *Object) Typescript(objects map[string]*Object) string {
	var parent string
	if t.HasParent() {
		parent = fmt.Sprintf(" extends %s", t.Parent.Typescript(objects))
	}
	var members []string
	for _, field := range t.Fields {
		members = append(members, fmt.Sprintf(`
		%s: %s;`, field.Name, field.Typescript(objects)))
	}

	if METH_PRE != nil {
		for _, field := range t.Fields {
			members = append(members, typeListTypescript(field.Name, "", []Type{field}, objects))
		}
	}

	for _, method := range t.Methods {
		if len(method.Params) > 0 {
			members = append(members, typeListTypescript(method.Name, "Param", method.Params, objects))
		}
		if len(method.Results) > 0 {
			members = append(members, typeListTypescript(method.Name, "Result", method.Results, objects))
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

func (t UnknownType) Typescript(objects map[string]*Object) string {
	return ""
}

func (t SimpleType) Typescript(objects map[string]*Object) string {
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

func (t *EnumType) Typescript(objects map[string]*Object) string {
	return t.Name
}

func (t *InstanceType) Typescript(objects map[string]*Object) string {
	if _, ok := objects[t.Name]; ok {
		return t.Name
	} else {
		return "Type"
	}
}

func (t *FixedPointType) Typescript(objects map[string]*Object) string {
	return "number"
}

func (t *ListType) Typescript(objects map[string]*Object) string {
	return fmt.Sprintf("%s[]", t.E.Typescript(objects))
}

func (t *DictType) Typescript(objects map[string]*Object) string {
	return fmt.Sprintf("{[index: %s]: %s}", t.K.Typescript(objects), t.V.Typescript(objects))
}

func (t *VariantType) Typescript(objects map[string]*Object) string {
	return "any"
}
