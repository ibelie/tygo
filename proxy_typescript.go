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

func Typescript(dir string, name string, types []Type) {
	var buffer bytes.Buffer

	var codes []string
	for _, t := range types {
		codes = append(codes, t.Typescript())
	}
	buffer.Write([]byte(fmt.Sprintf(`// Generated for tyts by tygo.  DO NOT EDIT!

declare module tyts.tygo {%s
}`, strings.Join(codes, ""))))

	ioutil.WriteFile(path.Join(dir, name+".d.ts"), buffer.Bytes(), 0666)
	Javascript(dir, name, types)
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

func (t *Method) Typescript() string {
	return ""
}

func (t *Object) Typescript() string {
	var fields []string
	for _, field := range t.Fields {
		fields = append(fields, fmt.Sprintf(`
		%s: %s;`, field.Name, field.Typescript()))
	}
	return fmt.Sprintf(`

	class %s {%s
	}`, t.Name, strings.Join(fields, ""))
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
	return t.Name
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
