// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"go/build"
	"io/ioutil"
)

const TEMP_PREFIX = "tmp"

var (
	LOG_PKG  = map[string]string{"log": ""}
	MATH_PKG = map[string]string{"math": ""}
	SRC_PATH = path.Join(os.Getenv("GOPATH"), "src")
	PROP_PRE []Type
	DELEGATE string
)

func Inject(dir string, filename string, pkgname string, types []Type, propPre []Type, delegate string) {
	injectfile := path.Join(SRC_PATH, dir, strings.Replace(filename, ".go", ".ty.go", 1))
	if types == nil {
		os.Remove(injectfile)
		return
	}

	desVarCount = 0
	var head bytes.Buffer
	var body bytes.Buffer
	head.Write([]byte(fmt.Sprintf(`// Generated by tygo.  DO NOT EDIT!

package %s
`, pkgname)))
	body.Write([]byte(`
`))

	PROP_PRE = propPre
	DELEGATE = delegate
	var pkgs map[string]string
	for _, t := range types {
		type_s, type_p := t.Go()
		pkgs = update(pkgs, type_p)
		body.Write([]byte(type_s))
	}
	PROP_PRE = nil
	DELEGATE = ""

	var sortedPkg []string
	for pkg, _ := range pkgs {
		sortedPkg = append(sortedPkg, pkg)
	}
	sort.Strings(sortedPkg)
	for _, pkg := range sortedPkg {
		head.Write([]byte(fmt.Sprintf(`
import %s"%s"`, pkgs[pkg], pkg)))
	}

	head.Write(body.Bytes())
	ioutil.WriteFile(injectfile, head.Bytes(), 0666)
}

func (t *Enum) Go() (string, map[string]string) {
	var names []string
	var values []string
	pkgs := updateTygo(nil)
	pkgs = update(pkgs, LOG_PKG)
	for _, name := range t.Sorted() {
		names = append(names, fmt.Sprintf(`
	case %s_%s:
		return "%s"`, t.Name, name, name))
		values = append(values, fmt.Sprintf(`
	%s_%s %s%s = %d`, t.Name, name, strings.Repeat(" ", t.nameMax-len(name)), t.Name, t.Values[name]))
	}
	bytesize_s, bytesize_p := t.ByteSizeGo("size", "i", "", 0, true)
	serialize_s, serialize_p := t.SerializeGo("size", "i", "", 0, true)
	deserialize_s, _, deserialize_p := t.DeserializeGo("", "input", "i", "", 0, false)
	pkgs = update(pkgs, bytesize_p)
	pkgs = update(pkgs, serialize_p)
	pkgs = update(pkgs, deserialize_p)
	return fmt.Sprintf(`
type %s uint

const (%s
)

func (i %s) String() string {
	switch i {%s
	default:
		log.Panicf("[Tygo][%s] Unexpect enum value: %%d", i)
		return "UNKNOWN"
	}
}

func (i %s) ByteSize() (size int) {%s
	return
}

func (i %s) CachedSize() int {
	return i.ByteSize()
}

func (i %s) Serialize(output *tygo.ProtoBuf) {%s
}

func (i *%s) Deserialize(input *tygo.ProtoBuf) (err error) {%s
	return
}
`, t.Name, strings.Join(values, ""), t.Name, strings.Join(names, ""),
		t.Name, t.Name, bytesize_s, t.Name, t.Name, serialize_s, t.Name, deserialize_s), pkgs
}

func TypeListSerialize(owner string, name string, typ string, ts []Type) (string, map[string]string) {
	if ts == nil {
		return "", nil
	}

	var pkgs map[string]string
	var items []string
	var itemsComment []string
	var itemsByteSize []string
	var itemsSerialize []string

	if typ == "" {
		for i, t := range PROP_PRE {
			n := fmt.Sprintf("p%d", i)
			item_s, item_p := t.Go()
			bytesize_s, bytesize_p := t.ByteSizeGo("size", n, "", 0, false)
			serialize_s, serialize_p := t.SerializeGo("size", n, "", 0, false)
			pkgs = update(pkgs, item_p)
			pkgs = update(pkgs, bytesize_p)
			pkgs = update(pkgs, serialize_p)

			items = append(items, fmt.Sprintf("p%d %s", i, item_s))
			itemsComment = append(itemsComment, fmt.Sprintf("p%d: %s", i, t))
			itemsByteSize = append(itemsByteSize, fmt.Sprintf(`
	// %s size: p%d%s
`, typ, i, bytesize_s))
			itemsSerialize = append(itemsSerialize, fmt.Sprintf(`
	// %s serialize: p%d%s
`, typ, i, serialize_s))
		}
	}

	for i, t := range ts {
		n := fmt.Sprintf("a%d", i)
		item_s, item_p := t.Go()
		bytesize_s, bytesize_p := t.ByteSizeGo("size", n, "", i+1, true)
		serialize_s, serialize_p := t.SerializeGo("size", n, "", i+1, true)
		pkgs = update(pkgs, item_p)
		pkgs = update(pkgs, bytesize_p)
		pkgs = update(pkgs, serialize_p)

		items = append(items, fmt.Sprintf("a%d %s", i, item_s))
		itemsComment = append(itemsComment, fmt.Sprintf("a%d: %s", i, t))
		itemsByteSize = append(itemsByteSize, fmt.Sprintf(`
	// %s size: a%d%s
`, typ, i, bytesize_s))
		itemsSerialize = append(itemsSerialize, fmt.Sprintf(`
	// %s serialize: a%d%s
`, typ, i, serialize_s))
	}

	Typ := strings.Title(typ)
	itemComment := strings.Join(itemsComment, ", ")

	return fmt.Sprintf(`
// %s %s(%s)
func (s *%s) Serialize%s%s(%s) (data []byte) {
	size := 0%s
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	output := &tygo.ProtoBuf{Buffer: data}
%s
	return
}
`, name, Typ, itemComment, owner, name, Typ, strings.Join(items, ", "),
		strings.Join(itemsByteSize, ""), strings.Join(itemsSerialize, "")), updateTygo(pkgs)
}

func TypeListDeserialize(owner string, name string, typ string, ts []Type) (string, map[string]string) {
	if ts == nil {
		return "", nil
	}

	var pkgs map[string]string
	var items []string
	var itemsComment []string
	var itemsDeserialize []string

	l := desVar()
	var deserialize_s string
	var deserialize_w WireType
	var deserialize_p map[string]string

	for i, t := range ts {
		n := fmt.Sprintf("a%d", i)
		item_s, item_p := t.Go()
		pkgs = update(pkgs, item_p)

		if i == 0 {
			deserialize_s, deserialize_w, deserialize_p = t.DeserializeGo("tag", "input", n, "", i+1, false)
			pkgs = update(pkgs, deserialize_p)
		}

		var next string
		var fall string
		var next_s string
		var next_w WireType
		var next_p map[string]string
		if i < len(ts)-1 {
			next_s, next_w, next_p = ts[i+1].DeserializeGo("tag", "input", fmt.Sprintf("a%d", i+1), "", i+2, false)
			pkgs = update(pkgs, next_p)
			tag_i, tag_ic := tagInt("", i+2, next_w)
			tag_s, tag_sc := expectTag("", i+2, next_w)
			next = fmt.Sprintf(`
					if !input.%s {%s
						continue method_%s // next tag for %s
					}
					tag = %s%s // fallthrough case %d`, tag_s, tag_sc, l, typ, tag_i, tag_ic, i+2)
			fall = fmt.Sprintf(` else {
					break switch_%s // skip tag
				}
				fallthrough`, l)
		} else {
			next = fmt.Sprintf(`
					if input.ExpectEnd() {
						break method_%s // end for %s
					}
					continue method_%s // next tag for %s`, l, typ, l, typ)
		}

		var listTag string
		var listComment string
		if l, ok := t.(*ListType); ok && l.E.IsPrimitive() {
			listTag = fmt.Sprintf(" || tag == %d", _MAKE_TAG(i+1, WireBytes))
			listComment = fmt.Sprintf(" || MAKE_TAG(%d, %s=%d)", i+1, WireBytes, WireBytes)
		}

		items = append(items, fmt.Sprintf("a%d %s", i, item_s))
		itemsComment = append(itemsComment, fmt.Sprintf("a%d: %s", i, t))
		itemsDeserialize = append(itemsDeserialize, fmt.Sprintf(`
			// %s deserialize: a%d
			case %d:
				if tag == %d%s { // MAKE_TAG(%d, %s=%d)%s%s%s
				}%s`, typ, i, i+1, _MAKE_TAG(i+1, deserialize_w), listTag, i+1, deserialize_w,
			deserialize_w, listComment, addIndent(deserialize_s, 4), next, fall))
		if i < len(ts)-1 {
			deserialize_s, deserialize_w, deserialize_p = next_s, next_w, next_p
		}
	}

	Typ := strings.Title(typ)
	itemComment := strings.Join(itemsComment, ", ")
	var switchLabel string
	if len(ts) > 1 {
		switchLabel = fmt.Sprintf(`
		switch_%s:`, l)
	}

	return fmt.Sprintf(`
// %s %s(%s)
func (s *%s) Deserialize%s%s(data []byte) (%s, err error) {
	input := &tygo.ProtoBuf{Buffer: data}
method_%s:
	for !input.ExpectEnd() {
		var tag int
		var cutoff bool
		if tag, cutoff, err = input.ReadTag(%s); err != nil {
			return
		} else if cutoff {%s
			switch %s {%s
			}
		}
		if err = input.SkipField(tag); err != nil {
			return
		}
	}
	return
}
`, name, Typ, itemComment, owner, name, Typ, strings.Join(items, ", "),
		l, _MAKE_CUTOFF_STR(strconv.Itoa(len(ts))), switchLabel, _TAG_FIELD_STR("tag"),
		strings.Join(itemsDeserialize, "")), updateTygo(pkgs)
}

func (t *Object) Go() (string, map[string]string) {
	var fields []string
	pkgs := updateTygo(nil)
	parent_s, parent_p := t.Parent.Go()
	pkgs = update(pkgs, parent_p)
	fields = append(fields, fmt.Sprintf(`
	%s`, parent_s))

	nameMax := 0
	typeMax := 0
	var preparedFields [][3]string
	for _, field := range t.Fields {
		field_s, field_p := field.Go()
		pkgs = update(pkgs, field_p)
		if nameMax < len(field.Name) {
			nameMax = len(field.Name)
		}
		if typeMax < len(field_s) {
			typeMax = len(field_s)
		}
		preparedFields = append(preparedFields, [3]string{field.Name, field_s, field.String()})
	}
	for _, field := range preparedFields {
		fields = append(fields, fmt.Sprintf(`
	%s %s%s %s// %s`, field[0], strings.Repeat(" ", nameMax-len(field[0])),
			field[1], strings.Repeat(" ", typeMax-len(field[1])), field[2]))
	}

	var methods []string
	if PROP_PRE != nil {
		for _, field := range t.Fields {
			field_s, field_p := TypeListSerialize(t.Name, field.Name, "", []Type{field})
			pkgs = update(pkgs, field_p)
			methods = append(methods, field_s)
		}
	}

	for _, method := range t.Methods {
		if len(method.Params) > 0 {
			param_d_s, param_d_p := TypeListDeserialize(t.Name, method.Name, "param", method.Params)
			param_s_s, param_s_p := TypeListSerialize(t.Name+DELEGATE, method.Name, "param", method.Params)
			pkgs = update(pkgs, param_d_p)
			pkgs = update(pkgs, param_s_p)
			methods = append(methods, param_d_s)
			methods = append(methods, param_s_s)
		}
		if len(method.Results) > 0 {
			result_d_s, result_d_p := TypeListDeserialize(t.Name+DELEGATE, method.Name, "result", method.Results)
			result_s_s, result_s_p := TypeListSerialize(t.Name, method.Name, "result", method.Results)
			pkgs = update(pkgs, result_d_p)
			pkgs = update(pkgs, result_s_p)
			methods = append(methods, result_d_s)
			methods = append(methods, result_s_s)
		}
	}

	mfn_n, mfn_i := t.MaxFieldNum()
	if mfn_n != "" {
		mfn_n = fmt.Sprintf("s.%s.MaxFieldNum() + ", mfn_n)
	}

	bytesize_s, bytesize_p := t.ByteSizeGo("size", "s", "", 0, true)
	serialize_s, serialize_p := t.SerializeGo("size", "s", "", 0, true)
	deserialize_s, _, deserialize_p := t.DeserializeGo("", "input", "s", "", 0, false)
	pkgs = update(pkgs, bytesize_p)
	pkgs = update(pkgs, serialize_p)
	pkgs = update(pkgs, deserialize_p)

	return fmt.Sprintf(`
type %s struct {%s
}

func (s *%s) MaxFieldNum() int {
	return %s%d
}

func (s *%s) ByteSize() (size int) {%s
	return
}

func (s *%s) Serialize(output *tygo.ProtoBuf) {%s
}

func (s *%s) Deserialize(input *tygo.ProtoBuf) (err error) {%s
	return
}
%s`, t.Name, strings.Join(fields, ""), t.Name, mfn_n, mfn_i, t.Name, bytesize_s,
		t.Name, serialize_s, t.Name, deserialize_s, strings.Join(methods, "")), pkgs
}

func (t UnknownType) Go() (string, map[string]string) {
	return string(t), nil
}

func (t SimpleType) Go() (string, map[string]string) {
	switch t {
	case SimpleType_BYTES:
		return "[]byte", nil
	default:
		return t.String(), nil
	}
}

func (t *EnumType) Go() (string, map[string]string) {
	return t.String(), nil
}

func (t *InstanceType) Go() (string, map[string]string) {
	if t.PkgPath == "" {
		return t.String(), nil
	} else {
		s := ""
		if t.IsPtr {
			s += "*"
		}
		var a string
		if p, err := build.Import(t.PkgPath, "", build.AllowBinary); err == nil && p.Name == t.PkgName {
			a = ""
		} else {
			a = t.PkgName + " "
		}
		return s + t.PkgName + "." + t.Name, map[string]string{t.PkgPath: a}
	}
}

func (t *FixedPointType) Go() (string, map[string]string) {
	return "float64", nil
}

func (t *ListType) Go() (string, map[string]string) {
	s, p := t.E.Go()
	return fmt.Sprintf("[]%s", s), p
}

func (t *DictType) Go() (string, map[string]string) {
	ks, kp := t.K.Go()
	vs, vp := t.V.Go()
	return fmt.Sprintf("map[%s]%s", ks, vs), update(kp, vp)
}

func (t *VariantType) Go() (string, map[string]string) {
	var p map[string]string
	for _, vt := range t.Ts {
		_, vp := vt.Go()
		p = update(p, vp)
	}
	return "interface{}", p
}
