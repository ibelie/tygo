// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"fmt"
	"log"
	"strings"
)

var PY_OBJECTS map[string]*Object

func (t *Enum) Python() string {
	var enums []string
	for _, name := range t.Sorted() {
		enums = append(enums, fmt.Sprintf(`
		%s = %d`, name, t.Values[name]))
	}
	return fmt.Sprintf(`

	const enum %s {%s
	}`, t.Name, strings.Join(enums, ","))
}

func (t *Object) Python() string {
	var parent string
	if t.HasParent() {
		parent = fmt.Sprintf(" extends %s", t.Parent.Python())
	}
	var members []string
	for _, field := range t.VisibleFields() {
		members = append(members, fmt.Sprintf(`
		%s: %s;`, field.Name, field.Python()))
	}

	if PROP_PRE != nil {
		for _, field := range t.VisibleFields() {
			members = append(members, field.Python())
			// members = append(members, typeListPython(field.Name, "", []Type{field}))
		}
	}

	for _, method := range t.Methods {
		if len(method.Params) > 0 {
			// members = append(members, typeListPython(method.Name, "Param", method.Params))
		}
		if len(method.Results) > 0 {
			// members = append(members, typeListPython(method.Name, "Result", method.Results))
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

func (t UnknownType) Python() string {
	return ""
}

func (t SimpleType) Python() string {
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
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for Python: %d", t)
		return "unknown"
	}
}

func (t *EnumType) Python() string {
	return t.Name
}

func (t *InstanceType) Python() string {
	if _, ok := TS_OBJECTS[t.Name]; ok {
		return t.Name
	} else {
		return "Type"
	}
}

func (t *FixedPointType) Python() string {
	return "number"
}

func (t *ListType) Python() string {
	return fmt.Sprintf("%s[]", t.E.Python())
}

func (t *DictType) Python() string {
	return fmt.Sprintf("{[index: %s]: %s}", t.K.Python(), t.V.Python())
}

func (t *VariantType) Python() string {
	return "any"
}

func (t *Enum) Typyd() string {
	var enums []string
	for _, name := range t.Sorted() {
		enums = append(enums, fmt.Sprintf(`
		%s = %d`, name, t.Values[name]))
	}
	return fmt.Sprintf(`

	const enum %s {%s
	}`, t.Name, strings.Join(enums, ","))
}

func (t *Object) Typyd() string {
	var parent string
	if t.HasParent() {
		parent = fmt.Sprintf(" extends %s", t.Parent.Typyd())
	}
	var members []string
	for _, field := range t.VisibleFields() {
		members = append(members, fmt.Sprintf(`
		%s: %s;`, field.Name, field.Typyd()))
	}

	if PROP_PRE != nil {
		for _, field := range t.VisibleFields() {
			members = append(members, field.Python())
			// members = append(members, typeListPython(field.Name, "", []Type{field}))
		}
	}

	for _, method := range t.Methods {
		if len(method.Params) > 0 {
			// members = append(members, typeListPython(method.Name, "Param", method.Params))
		}
		if len(method.Results) > 0 {
			// members = append(members, typeListPython(method.Name, "Result", method.Results))
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

func (t UnknownType) Typyd() string {
	return ""
}

func (t SimpleType) Typyd() string {
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
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for Typyd: %d", t)
		return "unknown"
	}
}

func (t *EnumType) Typyd() string {
	return t.Name
}

func (t *InstanceType) Typyd() string {
	if _, ok := TS_OBJECTS[t.Name]; ok {
		return t.Name
	} else {
		return "Type"
	}
}

func (t *FixedPointType) Typyd() string {
	return "number"
}

func (t *ListType) Typyd() string {
	return fmt.Sprintf("%s[]", t.E.Typyd())
}

func (t *DictType) Typyd() string {
	return fmt.Sprintf("{[index: %s]: %s}", t.K.Typyd(), t.V.Typyd())
}

func (t *VariantType) Typyd() string {
	return "any"
}
