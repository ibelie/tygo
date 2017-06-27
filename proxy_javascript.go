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

func Javascript(dir string, module string, types []Type) {
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
		js, rs := objects[name].Javascript(module, &body, genTypes, objects)
		requires = update(requires, rs)
		head.Write([]byte(fmt.Sprintf(`
goog.provide('%s.%s');`, module, name)))
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
	ioutil.WriteFile(path.Join(dir, module+".js"), head.Bytes(), 0666)
}

func (t *Enum) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	return "", nil
}

func (t *Method) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	return "", nil
}

func typeListJavascript(module string, name string, ts []Type, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	requires := map[string]string{"goog.require('tyts.Method');": ""}
	var items []string
	for i, t := range ts {
		js, rs := t.Javascript(module, writer, types, objects)
		update(requires, rs)
		items = append(items, fmt.Sprintf(`
	{tag: %d, tagsize: %d, type: %s}`, _MAKE_TAG(i+1, t.WireType()), TAG_SIZE(i+1), js))
	}
	return fmt.Sprintf(`new tyts.Method('%s', %d, [%s
]);
`, name, _MAKE_CUTOFF(len(items)), strings.Join(items, ",")), requires
}

func (t *Object) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	if _, exist := types[t.Name]; exist {
		return "", nil
	}

	var fields []string
	requires := map[string]string{"goog.require('tyts.Object');": ""}
	for i, field := range t.AllFields(objects) {
		wiretype := field.WireType()
		js, rs := field.Javascript(module, writer, types, objects)
		requires = update(requires, rs)
		fields = append(fields, fmt.Sprintf(`
	{name: '%s', tag: %d, tagsize: %d, type: %s}`, field.Name, _MAKE_TAG(i+1, wiretype), TAG_SIZE(i+1), js))
	}

	types[t.Name] = t
	var method_props []string
	var method_types []string
	for i, method := range t.Methods {
		js_p, rs_p := typeListJavascript(module, t.Name+method.Name+"Param", method.Params, writer, types, objects)
		js_r, rs_r := typeListJavascript(module, t.Name+method.Name+"Result", method.Results, writer, types, objects)
		update(requires, rs_p)
		update(requires, rs_r)
		method_props = append(method_props, fmt.Sprintf(`
	{name: '%s%s', type: null}`, method.Name, "Param"))
		method_props = append(method_props, fmt.Sprintf(`
	{name: '%s%s', type: null}`, method.Name, "Result"))
		method_types = append(method_types, fmt.Sprintf(`
%s.methods[%d].type = %s`, t.Name, i*2, js_p))
		method_types = append(method_types, fmt.Sprintf(`
%s.methods[%d].type = %s`, t.Name, i*2+1, js_r))
	}

	return fmt.Sprintf(`
var %s = new tyts.Object('%s', %d, [%s
], [%s
]);%s
%s.%s = %s.Type;
`, t.Name, t.Name, _MAKE_CUTOFF(len(fields)), strings.Join(fields, ","),
		strings.Join(method_props, ","), strings.Join(method_types, ""),
		module, t.Name, t.Name), requires
}

func (t UnknownType) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	return "", nil
}

func (t SimpleType) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
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

func (t *EnumType) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	return "tyts.Integer", map[string]string{"goog.require('tyts.Integer');": ""}
}

func (t *InstanceType) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	if object, ok := objects[t.Name]; ok {
		js, rs := object.Javascript(module, writer, types, objects)
		writer.Write([]byte(js))
		return t.Name, rs
	} else {
		identifier := t.Name + "Delegate"
		if _, exist := types[identifier]; !exist {
			writer.Write([]byte(fmt.Sprintf(`
var %s = new tyts.Extension('%s', %s)
`, identifier, identifier, t.Name)))
			types[identifier] = t
		}
		return identifier, map[string]string{
			"goog.require('tyts.Extension');":          "",
			fmt.Sprintf("goog.require('%s');", t.Name): "",
		}
	}
}

func (t *FixedPointType) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	identifier := t.Identifier()
	if _, exist := types[identifier]; !exist {
		writer.Write([]byte(fmt.Sprintf(`
var %s = new tyts.FixedPoint(%d, %d)
`, identifier, t.Floor, t.Precision)))
		types[identifier] = t
	}
	return identifier, map[string]string{"goog.require('tyts.FixedPoint');": ""}
}

func (t *ListType) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	identifier := t.Identifier()
	requires := map[string]string{"goog.require('tyts.List');": ""}
	if _, exist := types[identifier]; !exist {
		js, rs := t.E.Javascript(module, writer, types, objects)
		requires = update(requires, rs)
		writer.Write([]byte(fmt.Sprintf(`
var %s = new tyts.List('%s', %s)
`, identifier, identifier, js)))
		types[identifier] = t
	}
	return identifier, requires
}

func (t *DictType) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	identifier := t.Identifier()
	requires := map[string]string{"goog.require('tyts.Dict');": ""}
	if _, exist := types[identifier]; !exist {
		js_k, rs_k := t.K.Javascript(module, writer, types, objects)
		js_v, rs_v := t.V.Javascript(module, writer, types, objects)
		requires = update(requires, rs_k)
		requires = update(requires, rs_v)
		writer.Write([]byte(fmt.Sprintf(`
var %s = new tyts.Dict('%s', %s, %s)
`, identifier, identifier, js_k, js_v)))
		types[identifier] = t
	}
	return identifier, requires
}

func (t *VariantType) Javascript(module string, writer io.Writer, types map[string]Type, objects map[string]*Object) (string, map[string]string) {
	identifier := t.Identifier()
	requires := map[string]string{"goog.require('tyts.Variant');": ""}
	if _, exist := types[identifier]; !exist {
		var codes []string
		variantNum := 0
		for _, st := range t.Ts {
			if s, ok := st.(SimpleType); ok && s == SimpleType_NIL {
				continue
			}
			variantNum++
			js, rs := st.Javascript(module, writer, types, objects)
			wiretype := st.WireType()
			requires = update(requires, rs)
			codes = append(codes, fmt.Sprintf(`
	{tag: %d, tagsize: %d, type: %s}`, _MAKE_TAG(variantNum, wiretype), TAG_SIZE(variantNum), js))
		}
		writer.Write([]byte(fmt.Sprintf(`
var %s = new tyts.Variant('%s', %d, [%s
])
`, identifier, identifier, _MAKE_CUTOFF(len(codes)), strings.Join(codes, ","))))
		types[identifier] = t
	}
	return identifier, requires
}
