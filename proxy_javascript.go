// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path"
	"sort"
	"strings"

	"io/ioutil"
)

func Javascript(dir string, name string, types []Type) {
	var head bytes.Buffer
	var body bytes.Buffer
	head.Write([]byte(`// Generated for tyts by tygo.  DO NOT EDIT!
`))
	body.Write([]byte(`
`))

	var sortedObjects []string
	objects := make(map[string]*Object)
	for _, t := range types {
		if object, ok := t.(*Object); ok {
			if o, exist := objects[object.Name]; exist {
				log.Fatalf("[Tygo][Javascript] Object already exists: %v %v", o, object)
			}
			objects[object.Name] = object
			sortedObjects = append(sortedObjects, object.Name)
		}
	}
	sort.Strings(sortedObjects)
	for _, name := range sortedObjects {
		head.Write([]byte(fmt.Sprintf(`
goog.provide('tyts.tygo.%s');`, name)))
	}
	head.Write([]byte(`

goog.require('tyts.Integer');
goog.require('tyts.FixedPoint');
goog.require('tyts.Float64');
goog.require('tyts.Float32');
goog.require('tyts.Bool');
goog.require('tyts.Bytes');
goog.require('tyts.String');
goog.require('tyts.Object');
goog.require('tyts.Variant');
goog.require('tyts.List');
goog.require('tyts.Dict');
goog.require('tyts.Extension');
`))

	genTypes := make(map[string]Type)
	for _, name := range sortedObjects {
		body.Write([]byte(objects[name].Javascript(&head, &body, genTypes, objects)))
	}

	head.Write(body.Bytes())
	ioutil.WriteFile(path.Join(dir, name+".js"), head.Bytes(), 0666)
}

func (t *Enum) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	return ""
}

func (t *Method) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	return ""
}

func (t *Object) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	if _, exist := types[t.Name]; exist {
		return ""
	}
	fields := t.AllFields(objects)
	var codes []string
	for i, field := range fields {
		wiretype := field.WireType()
		codes = append(codes, fmt.Sprintf(`
	{name: %s, tag: %d, tagsize: %d, type: %s}`, field.Name, _MAKE_TAG(i+1, wiretype), TAG_SIZE(i+1), field.Javascript(head, body, types, objects)))
	}
	types[t.Name] = t
	return fmt.Sprintf(`
tyts.tygo.%s = tyts.Object('%s', %d, [%s
]);
`, t.Name, t.Name, _MAKE_CUTOFF(len(fields)), strings.Join(codes, ","))
}

func (t UnknownType) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	return ""
}

func (t SimpleType) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	switch t {
	case SimpleType_INT32:
		fallthrough
	case SimpleType_INT64:
		fallthrough
	case SimpleType_UINT32:
		fallthrough
	case SimpleType_UINT64:
		return "tyts.Integer"
	case SimpleType_BYTES:
		return "tyts.Bytes"
	case SimpleType_STRING:
		return "tyts.String"
	case SimpleType_BOOL:
		return "tyts.Bool"
	case SimpleType_FLOAT32:
		return "tyts.Float32"
	case SimpleType_FLOAT64:
		return "tyts.Float64"
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for Javascript: %d", t)
		return ""
	}
}

func (t *EnumType) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	return "tyts.Integer"
}

func (t *InstanceType) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	identifier := t.Name
	if object, ok := objects[identifier]; ok {
		body.Write([]byte(object.Javascript(head, body, types, objects)))
	} else {
		identifier = t.Name + "Delegate"
		if _, exist := types[identifier]; !exist {
			head.Write([]byte(fmt.Sprintf(`
goog.require('%s');`, t.Name)))
			body.Write([]byte(fmt.Sprintf(`
%s = tyts.Extension('%s', %s)
`, identifier, identifier, t.Name)))
			types[identifier] = t
		}
	}
	return identifier
}

func (t *FixedPointType) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	identifier := t.Identifier()
	if _, exist := types[identifier]; !exist {
		body.Write([]byte(fmt.Sprintf(`
%s = tyts.FixedPoint(%d, %d)
`, identifier, t.Floor, t.Precision)))
		types[identifier] = t
	}
	return identifier
}

func (t *ListType) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	identifier := t.Identifier()
	if _, exist := types[identifier]; !exist {
		body.Write([]byte(fmt.Sprintf(`
%s = tyts.List('%s', %s)
`, identifier, identifier, t.E.Javascript(head, body, types, objects))))
		types[identifier] = t
	}
	return identifier
}

func (t *DictType) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	identifier := t.Identifier()
	if _, exist := types[identifier]; !exist {
		body.Write([]byte(fmt.Sprintf(`
%s = tyts.Dict('%s', %s, %s)
`, identifier, identifier, t.K.Javascript(head, body, types, objects), t.V.Javascript(head, body, types, objects))))
		types[identifier] = t
	}
	return identifier
}

func (t *VariantType) Javascript(head io.Writer, body io.Writer, types map[string]Type, objects map[string]*Object) string {
	identifier := t.Identifier()
	if _, exist := types[identifier]; !exist {
		var codes []string
		for i, st := range t.Ts {
			if s, ok := st.(SimpleType); ok && s == SimpleType_NIL {
				continue
			}
			wiretype := st.WireType()
			codes = append(codes, fmt.Sprintf(`
	{tag: %d, tagsize: %d, type: %s}`, _MAKE_TAG(i+1, wiretype), TAG_SIZE(i+1), st.Javascript(head, body, types, objects)))
		}
		body.Write([]byte(fmt.Sprintf(`
%s = tyts.Variant('%s', [%s
])
`, identifier, identifier, strings.Join(codes, ","))))
		types[identifier] = t
	}
	return identifier
}
