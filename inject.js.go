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

var (
	JS_MODULE  = ""
	JS_WRITER  io.Writer
	JS_TYPES   map[string]Type
	JS_OBJECTS map[string]*Object
)

func Javascript(dir string, name string, module string, types []Type, propPre []Type) {
	var head bytes.Buffer
	var body bytes.Buffer
	head.Write([]byte(`// Generated for tyts by tygo.  DO NOT EDIT!
`))
	body.Write([]byte(`
`))

	JS_OBJECTS = ObjectMap(types, false)
	var sortedObjects []string
	for n, _ := range JS_OBJECTS {
		sortedObjects = append(sortedObjects, n)
	}
	sort.Strings(sortedObjects)

	PROP_PRE = propPre
	JS_MODULE = module
	JS_WRITER = &body
	JS_TYPES = make(map[string]Type)
	var requires map[string]string
	for _, name := range sortedObjects {
		object := JS_OBJECTS[name]
		m := JS_MODULE
		if m == "" {
			m = object.Package
		}
		js, rs := JS_OBJECTS[name].Javascript()
		requires = update(requires, rs)
		head.Write([]byte(fmt.Sprintf(`
goog.provide('%s.%s');`, strings.Replace(m, "/", ".", -1), name)))
		body.Write([]byte(js))
	}
	PROP_PRE = nil
	JS_TYPES = nil
	JS_WRITER = nil
	JS_OBJECTS = nil
	JS_MODULE = ""

	var sortedRequires []string
	for require, _ := range requires {
		if strings.HasPrefix(require, "goog.require('ibelie.tyts.") {
			sortedRequires = append(sortedRequires, require)
		}
	}
	sort.Strings(sortedRequires)
	head.Write([]byte("\n\n" + strings.Join(sortedRequires, "\n")))

	sortedRequires = nil
	for require, _ := range requires {
		if !strings.HasPrefix(require, "goog.require('ibelie.tyts.") {
			sortedRequires = append(sortedRequires, require)
		}
	}
	if len(sortedRequires) > 0 {
		sort.Strings(sortedRequires)
		head.Write([]byte("\n\n" + strings.Join(sortedRequires, "\n")))
	}

	if name == "" {
		name = module
	}
	head.Write(body.Bytes())
	ioutil.WriteFile(path.Join(dir, name+".js"), head.Bytes(), 0666)
}

func (t *Enum) Javascript() (string, map[string]string) {
	return "", nil
}

func typeListJavascript(name string, ts []Type) (string, map[string]string) {
	requires := map[string]string{"goog.require('ibelie.tyts.Method');": ""}
	var items []string
	for i, t := range ts {
		js, rs := t.Javascript()
		update(requires, rs)
		items = append(items, fmt.Sprintf(`
	{tag: %d, tagsize: %d, type: %s}`, _MAKE_TAG(i+1, t.WireType()), TAG_SIZE(i+1), js))
	}
	return fmt.Sprintf(`new ibelie.tyts.Method('%s', %d, [%s
]);
`, name, _MAKE_CUTOFF(len(items)), strings.Join(items, ",")), requires
}

func (t *Object) Javascript() (string, map[string]string) {
	if _, exist := JS_TYPES[t.Name]; exist {
		return "", nil
	}

	var fields []string
	requires := map[string]string{"goog.require('ibelie.tyts.Object');": ""}
	for i, field := range t.AllFields(JS_OBJECTS, true) {
		wiretype := field.WireType()
		js, rs := field.Javascript()
		requires = update(requires, rs)
		fields = append(fields, fmt.Sprintf(`
	{name: '%s', tag: %d, tagsize: %d, type: %s}`, field.Name, _MAKE_TAG(i+1, wiretype), TAG_SIZE(i+1), js))
	}

	JS_TYPES[t.Name] = t
	var method_props []string
	var method_types []string
	method_index := 0

	if PROP_PRE != nil {
		for _, field := range t.VisibleFields() {
			js, rs := typeListJavascript(t.Name+field.Name, []Type{field})
			update(requires, rs)
			method_props = append(method_props, fmt.Sprintf(`
	{name: '%s', type: null}`, field.Name))
			method_types = append(method_types, fmt.Sprintf(`
%s.methods[%d].type = %s`, t.Name, method_index, js))
			method_index++
		}
	}

	for _, method := range t.Methods {
		if len(method.Params) > 0 {
			js_p, rs_p := typeListJavascript(t.Name+method.Name+"Param", method.Params)
			update(requires, rs_p)
			method_props = append(method_props, fmt.Sprintf(`
	{name: '%s%s', type: null}`, method.Name, "Param"))
			method_types = append(method_types, fmt.Sprintf(`
%s.methods[%d].type = %s`, t.Name, method_index, js_p))
			method_index++
		}
		if len(method.Results) > 0 {
			js_r, rs_r := typeListJavascript(t.Name+method.Name+"Result", method.Results)
			update(requires, rs_r)
			method_props = append(method_props, fmt.Sprintf(`
	{name: '%s%s', type: null}`, method.Name, "Result"))
			method_types = append(method_types, fmt.Sprintf(`
%s.methods[%d].type = %s`, t.Name, method_index, js_r))
			method_index++
		}
	}

	m := JS_MODULE
	if m == "" {
		m = t.Package
	}

	return fmt.Sprintf(`
var %s = new ibelie.tyts.Object('%s', %d, [%s
], [%s
]);%s
%s['%s'] = %s.Type;
`, t.Name, t.Name, _MAKE_CUTOFF(len(fields)), strings.Join(fields, ","),
		strings.Join(method_props, ","), strings.Join(method_types, ""),
		strings.Replace(m, "/", ".", -1), t.Name, t.Name), requires
}

func (t UnknownType) Javascript() (string, map[string]string) {
	return "", nil
}

func (t SimpleType) Javascript() (string, map[string]string) {
	switch t {
	case SimpleType_INT32:
		fallthrough
	case SimpleType_INT64:
		fallthrough
	case SimpleType_UINT32:
		fallthrough
	case SimpleType_UINT64:
		return "ibelie.tyts.Integer", map[string]string{"goog.require('ibelie.tyts.Integer');": ""}
	case SimpleType_BYTES:
		return "ibelie.tyts.Bytes", map[string]string{"goog.require('ibelie.tyts.Bytes');": ""}
	case SimpleType_STRING:
		return "ibelie.tyts.String", map[string]string{"goog.require('ibelie.tyts.String');": ""}
	case SimpleType_SYMBOL:
		return "ibelie.tyts.Symbol", map[string]string{"goog.require('ibelie.tyts.Symbol');": ""}
	case SimpleType_BOOL:
		return "ibelie.tyts.Bool", map[string]string{"goog.require('ibelie.tyts.Bool');": ""}
	case SimpleType_FLOAT32:
		return "ibelie.tyts.Float32", map[string]string{"goog.require('ibelie.tyts.Float32');": ""}
	case SimpleType_FLOAT64:
		return "ibelie.tyts.Float64", map[string]string{"goog.require('ibelie.tyts.Float64');": ""}
	default:
		log.Fatalf("[Tygo][SimpleType] Unexpect enum value for Javascript: %d", t)
		return "", nil
	}
}

func (t *EnumType) Javascript() (string, map[string]string) {
	return "ibelie.tyts.Integer", map[string]string{"goog.require('ibelie.tyts.Integer');": ""}
}

func (t *InstanceType) Javascript() (string, map[string]string) {
	if object, ok := JS_OBJECTS[t.Name]; ok {
		js, rs := object.Javascript()
		JS_WRITER.Write([]byte(js))
		return t.Name, rs
	} else {
		identifier := t.Name + "Delegate"
		if _, exist := JS_TYPES[identifier]; !exist {
			JS_WRITER.Write([]byte(fmt.Sprintf(`
var %s = new ibelie.tyts.Extension('%s', %s)
`, identifier, identifier, t.Name)))
			JS_TYPES[identifier] = t
		}
		return identifier, map[string]string{
			"goog.require('ibelie.tyts.Extension');":   "",
			fmt.Sprintf("goog.require('%s');", t.Name): "",
		}
	}
}

func (t *FixedPointType) Javascript() (string, map[string]string) {
	identifier := t.Identifier()
	if _, exist := JS_TYPES[identifier]; !exist {
		JS_WRITER.Write([]byte(fmt.Sprintf(`
var %s = new ibelie.tyts.FixedPoint(%d, %d)
`, identifier, t.Floor, t.Precision)))
		JS_TYPES[identifier] = t
	}
	return identifier, map[string]string{"goog.require('ibelie.tyts.FixedPoint');": ""}
}

func (t *ListType) Javascript() (string, map[string]string) {
	identifier := t.Identifier()
	requires := map[string]string{"goog.require('ibelie.tyts.List');": ""}
	if _, exist := JS_TYPES[identifier]; !exist {
		js, rs := t.E.Javascript()
		requires = update(requires, rs)
		JS_WRITER.Write([]byte(fmt.Sprintf(`
var %s = new ibelie.tyts.List('%s', %s)
`, identifier, identifier, js)))
		JS_TYPES[identifier] = t
	}
	return identifier, requires
}

func (t *DictType) Javascript() (string, map[string]string) {
	identifier := t.Identifier()
	requires := map[string]string{"goog.require('ibelie.tyts.Dict');": ""}
	if _, exist := JS_TYPES[identifier]; !exist {
		js_k, rs_k := t.K.Javascript()
		js_v, rs_v := t.V.Javascript()
		requires = update(requires, rs_k)
		requires = update(requires, rs_v)
		JS_WRITER.Write([]byte(fmt.Sprintf(`
var %s = new ibelie.tyts.Dict('%s', %s, %s)
`, identifier, identifier, js_k, js_v)))
		JS_TYPES[identifier] = t
	}
	return identifier, requires
}

func (t *VariantType) Javascript() (string, map[string]string) {
	identifier := t.Identifier()
	requires := map[string]string{"goog.require('ibelie.tyts.Variant');": ""}
	if _, exist := JS_TYPES[identifier]; !exist {
		var codes []string
		variantNum := 0
		for _, st := range t.Ts {
			if s, ok := st.(SimpleType); ok && s == SimpleType_NIL {
				continue
			}
			variantNum++
			js, rs := st.Javascript()
			wiretype := st.WireType()
			requires = update(requires, rs)
			codes = append(codes, fmt.Sprintf(`
	{tag: %d, tagsize: %d, type: %s}`, _MAKE_TAG(variantNum, wiretype), TAG_SIZE(variantNum), js))
		}
		JS_WRITER.Write([]byte(fmt.Sprintf(`
var %s = new ibelie.tyts.Variant('%s', %d, [%s
])
`, identifier, identifier, _MAKE_CUTOFF(len(codes)), strings.Join(codes, ","))))
		JS_TYPES[identifier] = t
	}
	return identifier, requires
}
