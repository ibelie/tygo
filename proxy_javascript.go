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

	var requires map[string]string
	genTypes := make(map[string]Type)
	for _, name := range sortedObjects {
		js, rs := objects[name].Javascript(&body, genTypes, objects)
		requires = update(requires, rs)
		head.Write([]byte(fmt.Sprintf(`
goog.provide('tyts.tygo.%s');`, name)))
		body.Write([]byte(js))
	}

	var sortedRequires []string
	for require, _ := range requires {
		if strings.HasPrefix(require, "goog.require('tyts.") {
			sortedRequires = append(sortedRequires, require)
		}
	}
	sort.Strings(sortedRequires)
	head.Write([]byte("\n\n" + strings.Join(sortedRequires, "\n")))

	sortedRequires = nil
	for require, _ := range requires {
		if !strings.HasPrefix(require, "goog.require('tyts.") {
			sortedRequires = append(sortedRequires, require)
		}
	}
	sort.Strings(sortedRequires)
	head.Write([]byte("\n\n" + strings.Join(sortedRequires, "\n")))

	head.Write(body.Bytes())
	ioutil.WriteFile(path.Join(dir, name+".js"), head.Bytes(), 0666)
}

func (t *Enum) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	return "", nil
}

func (t *Method) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	return "", nil
}

func typeListJavascript(name string, ts []Type, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	requires := map[string]string{"goog.require('tyts.Method');": ""}
	if _, exist := types[name]; !exist {
		var items []string
		for i, t := range ts {
			js, rs := t.Javascript(writer, types, objects)
			update(requires, rs)
			items = append(items, fmt.Sprintf(`
	{tag: %d, tagsize: %d, type: %s}`, _MAKE_TAG(i+1, t.WireType()), TAG_SIZE(i+1), js))
		}
		writer.Write([]byte(fmt.Sprintf(`
%s = tyts.Method('%s', %d, [%s
]);
`, name, name, _MAKE_CUTOFF(len(items)), strings.Join(items, ", "))))
		types[name] = SimpleType_NIL
	}
	return name, requires
}

func (t *Object) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	if _, exist := types[t.Name]; exist {
		return "", nil
	}
	types[t.Name] = t

	var fields []string
	requires := map[string]string{"goog.require('tyts.Object');": ""}
	for i, field := range t.AllFields(objects) {
		wiretype := field.WireType()
		js, rs := field.Javascript(writer, types, objects)
		requires = update(requires, rs)
		fields = append(fields, fmt.Sprintf(`
	{name: '%s', tag: %d, tagsize: %d, type: %s}`, field.Name, _MAKE_TAG(i+1, wiretype), TAG_SIZE(i+1), js))
	}

	var methods []string
	for _, method := range t.Methods {
		js_p, rs_p := typeListJavascript(t.Name+method.Name+"Param", method.Params, writer, types, objects)
		js_r, rs_r := typeListJavascript(t.Name+method.Name+"Result", method.Results, writer, types, objects)
		update(requires, rs_p)
		update(requires, rs_r)
		methods = append(methods, fmt.Sprintf(`
	{name: '%s%s', type: %s}`, method.Name, "Param", js_p))
		methods = append(methods, fmt.Sprintf(`
	{name: '%s%s', type: %s}`, method.Name, "Result", js_r))
	}

	return fmt.Sprintf(`
%s = tyts.Object('%s', %d, [%s
], [%s
]);
tyts.tygo.%s = %s.Type;
`, t.Name, t.Name, _MAKE_CUTOFF(len(fields)), strings.Join(fields, ","), strings.Join(methods, ","), t.Name, t.Name), requires
}

func (t UnknownType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	return "", nil
}

func (t SimpleType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	switch t {
	case SimpleType_INT32:
		fallthrough
	case SimpleType_INT64:
		fallthrough
	case SimpleType_UINT32:
		fallthrough
	case SimpleType_UINT64:
		return "tyts.Integer", map[string]string{"goog.require('tyts.Integer');": ""}
	case SimpleType_BYTES:
		return "tyts.Bytes", map[string]string{"goog.require('tyts.Bytes');": ""}
	case SimpleType_STRING:
		return "tyts.String", map[string]string{"goog.require('tyts.String');": ""}
	case SimpleType_BOOL:
		return "tyts.Bool", map[string]string{"goog.require('tyts.Bool');": ""}
	case SimpleType_FLOAT32:
		return "tyts.Float32", map[string]string{"goog.require('tyts.Float32');": ""}
	case SimpleType_FLOAT64:
		return "tyts.Float64", map[string]string{"goog.require('tyts.Float64');": ""}
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for Javascript: %d", t)
		return "", nil
	}
}

func (t *EnumType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	return "tyts.Integer", map[string]string{"goog.require('tyts.Integer');": ""}
}

func (t *InstanceType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	if object, ok := objects[t.Name]; ok {
		js, rs := object.Javascript(writer, types, objects)
		writer.Write([]byte(js))
		return t.Name, rs
	} else {
		identifier := t.Name + "Delegate"
		if _, exist := types[identifier]; !exist {
			writer.Write([]byte(fmt.Sprintf(`
%s = tyts.Extension('%s', %s)
`, identifier, identifier, t.Name)))
			types[identifier] = t
		}
		return identifier, map[string]string{
			"goog.require('tyts.Extension');":          "",
			fmt.Sprintf("goog.require('%s');", t.Name): "",
		}
	}
}

func (t *FixedPointType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	identifier := t.Identifier()
	if _, exist := types[identifier]; !exist {
		writer.Write([]byte(fmt.Sprintf(`
%s = tyts.FixedPoint(%d, %d)
`, identifier, t.Floor, t.Precision)))
		types[identifier] = t
	}
	return identifier, map[string]string{"goog.require('tyts.FixedPoint');": ""}
}

func (t *ListType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	identifier := t.Identifier()
	requires := map[string]string{"goog.require('tyts.List');": ""}
	if _, exist := types[identifier]; !exist {
		js, rs := t.E.Javascript(writer, types, objects)
		requires = update(requires, rs)
		writer.Write([]byte(fmt.Sprintf(`
%s = tyts.List('%s', %s)
`, identifier, identifier, js)))
		types[identifier] = t
	}
	return identifier, requires
}

func (t *DictType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	identifier := t.Identifier()
	requires := map[string]string{"goog.require('tyts.Dict');": ""}
	if _, exist := types[identifier]; !exist {
		js_k, rs_k := t.K.Javascript(writer, types, objects)
		js_v, rs_v := t.V.Javascript(writer, types, objects)
		requires = update(requires, rs_k)
		requires = update(requires, rs_v)
		writer.Write([]byte(fmt.Sprintf(`
%s = tyts.Dict('%s', %s, %s)
`, identifier, identifier, js_k, js_v)))
		types[identifier] = t
	}
	return identifier, requires
}

func (t *VariantType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	identifier := t.Identifier()
	requires := map[string]string{"goog.require('tyts.Variant');": ""}
	if _, exist := types[identifier]; !exist {
		var codes []string
		for i, st := range t.Ts {
			if s, ok := st.(SimpleType); ok && s == SimpleType_NIL {
				continue
			}
			js, rs := st.Javascript(writer, types, objects)
			wiretype := st.WireType()
			requires = update(requires, rs)
			codes = append(codes, fmt.Sprintf(`
	{tag: %d, tagsize: %d, type: %s}`, _MAKE_TAG(i+1, wiretype), TAG_SIZE(i+1), js))
		}
		writer.Write([]byte(fmt.Sprintf(`
%s = tyts.Variant('%s', %d, [%s
])
`, identifier, identifier, _MAKE_CUTOFF(len(codes)), strings.Join(codes, ","))))
		types[identifier] = t
	}
	return identifier, requires
}
