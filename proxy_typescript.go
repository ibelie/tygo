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

func Typescript(dir string, module string, types []Type) {
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

	var codes []string
	for _, t := range types {
		codes = append(codes, t.Typescript(objects))
	}
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

	ioutil.WriteFile(path.Join(dir, module+".d.ts"), buffer.Bytes(), 0666)
	Javascript(dir, module, types)
}

func (t *Enum) Typescript(objects map[string]*Object) string {
	var enums []string
	for _, name := range t.Sorted() {
		enums = append(enums, fmt.Sprintf(`
		%s = %d`, name, t.Values[name]))
	}
	return fmt.Sprintf(`

	export const enum %s {%s
	}`, t.Name, strings.Join(enums, ","))
}

func (t *Method) Typescript(objects map[string]*Object) string {
	return ""
}

func typeListTypescript(name string, typ string, ts []Type, objects map[string]*Object) string {
	var items []string
	for i, t := range ts {
		items = append(items, fmt.Sprintf("a%d: %s", i, t.Typescript(objects)))
	}
	return fmt.Sprintf(`
		Serialize%s%s(%s): Uint8Array;
		Deserialize%s%s(data: Uint8Array): any;`, name, typ, strings.Join(items, ", "), name, typ)
}

func (t *Object) Typescript(objects map[string]*Object) string {
	var members []string
	for _, field := range t.Fields {
		members = append(members, fmt.Sprintf(`
		%s: %s;`, field.Name, field.Typescript(objects)))
	}

	for _, method := range t.Methods {
		members = append(members, typeListTypescript(method.Name, "Param", method.Params, objects))
		members = append(members, typeListTypescript(method.Name, "Result", method.Results, objects))
	}

	return fmt.Sprintf(`

	export class %s {
		__class__: string;
		constructor();
		ByteSize(): number;
		Serialize(): Uint8Array;
		Deserialize(data: Uint8Array): void;
%s
	}

	export namespace %s {
		function Deserialize(data: Uint8Array): %s;
	}`, t.Name, strings.Join(members, ""), t.Name, t.Name)
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
