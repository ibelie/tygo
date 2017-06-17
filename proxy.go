// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
)

const (
	TEMP_PREFIX = "tmp"
	goHeader    = `// Generated by tygo.  DO NOT EDIT!

package %s

`
	goImport = `import %s"%s"
`
)

var (
	FMT_PKG  = map[string]string{"fmt": ""}
	MATH_PKG = map[string]string{"math": ""}
	SRC_PATH = os.Getenv("GOPATH") + "/src/"
)

func Inject(path string) {
	buildPackage, err := build.Import(path, "", build.ImportComment)
	if err != nil {
		log.Fatalf("[Tygo][Inject] Cannot import package:\n>>>>%v", err)
		return
	}
	fs := token.NewFileSet()
	for _, filename := range buildPackage.GoFiles {
		if strings.HasSuffix(filename, ".ty.go") {
			continue
		}
		file, err := parser.ParseFile(fs, buildPackage.Dir+"/"+filename, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("[Tygo][Inject] Cannot parse file:\n>>>>%v", err)
		}
		for _, d := range file.Decls {
			decl, ok := d.(*ast.GenDecl)
			if !ok || decl.Tok != token.IMPORT {
				continue
			}
			for _, s := range decl.Specs {
				spec, ok := s.(*ast.ImportSpec)
				if !ok || strings.Trim(spec.Path.Value, "\"") != TYGO_PATH {
					continue
				}
				injectfile := SRC_PATH + path + "/" + strings.Replace(filename, ".go", ".ty.go", 1)
				if strings.TrimSpace(decl.Doc.Text()) == "" {
					os.Remove(injectfile)
				} else {
					inject(injectfile, decl.Doc.Text(), file)
				}
			}
		}
	}
}

func inject(filename string, doc string, file *ast.File) {
	desVarCount = 0
	imports := make(map[string]string)
	typePkg := make(map[string][2]string)
	for _, importSpec := range file.Imports {
		pkg := strings.Trim(importSpec.Path.Value, "\"")
		if importSpec.Name == nil {
			if p, err := build.Import(pkg, "", build.AllowBinary); err != nil {
				log.Fatalf("[Tygo][Inject] Cannot import package:\n>>>>%v", err)
			} else {
				imports[p.Name] = p.ImportPath
			}
		} else if importSpec.Name.Name == "." {
			if doc := packageDoc(pkg); doc != nil {
				for _, t := range doc.Types {
					typePkg[t.Name] = [2]string{doc.Name, pkg}
				}
			}
		} else {
			imports[importSpec.Name.Name] = pkg
		}
	}
	types := Parse(doc, imports, typePkg)

	var head bytes.Buffer
	var body bytes.Buffer
	head.Write([]byte(fmt.Sprintf(goHeader, file.Name)))

	var pkgs map[string]string
	for _, t := range types {
		type_s, type_p := t.Go()
		pkgs = update(pkgs, type_p)
		body.Write([]byte(type_s))
	}
	var sortedPkg []string
	for path, _ := range pkgs {
		sortedPkg = append(sortedPkg, path)
	}
	sort.Strings(sortedPkg)
	for _, path := range sortedPkg {
		head.Write([]byte(fmt.Sprintf(goImport, pkgs[path], path)))
	}

	head.Write(body.Bytes())
	ioutil.WriteFile(filename, head.Bytes(), 0666)
}

func (t *Enum) Go() (string, map[string]string) {
	var names []string
	var values []string
	pkgs := map[string]string{"fmt": "", TYGO_PATH: ""}
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
		panic(fmt.Sprintf("[Tygo][%s] Unexpect enum value: %%d", i))
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

func (t *Method) Go() (string, map[string]string) {
	return "", nil
}

func typeListGo(owner string, name string, typ string, ts []Type) (string, map[string]string) {
	if ts == nil {
		return "", nil
	}

	var pkgs map[string]string
	var items []string
	var itemsComment []string
	var itemsByteSize []string
	var itemsSerialize []string
	var itemsDeserialize []string

	l := desVar()
	var deserialize_s string
	var deserialize_w WireType
	var deserialize_p map[string]string

	for i, t := range ts {
		n := fmt.Sprintf("a%d", i)
		item_s, item_p := t.Go()
		bytesize_s, bytesize_p := t.ByteSizeGo("size", n, "", i+1, true)
		serialize_s, serialize_p := t.SerializeGo("size", n, "", i+1, true)
		pkgs = update(pkgs, item_p)
		pkgs = update(pkgs, bytesize_p)
		pkgs = update(pkgs, serialize_p)

		if i == 0 {
			deserialize_s, deserialize_w, deserialize_p = t.DeserializeGo("tag", "input", n, "", i+1, false)
			pkgs = update(pkgs, deserialize_p)
		}

		var next string
		var fall string
		if i < len(ts)-1 {
			next_s, next_w, next_p := ts[i+1].DeserializeGo("tag", "input", fmt.Sprintf("a%d", i+1), "", i+2, false)
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
			deserialize_s, deserialize_w, deserialize_p = next_s, next_w, next_p
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
		itemsByteSize = append(itemsByteSize, fmt.Sprintf(`
	// %s size: a%d%s
`, typ, i, bytesize_s))
		itemsSerialize = append(itemsSerialize, fmt.Sprintf(`
	// %s serialize: a%d%s
`, typ, i, serialize_s))
		itemsDeserialize = append(itemsDeserialize, fmt.Sprintf(`
			// %s deserialize: a%d
			case %d:
				if tag == %d%s { // MAKE_TAG(%d, %s=%d)%s%s%s
				}%s`, typ, i, i+1, _MAKE_TAG(i+1, deserialize_w), listTag, i+1, deserialize_w,
			deserialize_w, listComment, addIndent(deserialize_s, 3), next, fall))
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
func %sSerialize%s%s(%s) (data []byte) {
	size := 0%s
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	output := &tygo.ProtoBuf{Buffer: data}
%s
	return
}

// %s %s(%s)
func %sDeserialize%s%s(data []byte) (%s, err error) {
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
		} else if err = input.SkipField(tag); err != nil {
			return
		}
	}
	return
}
`, name, Typ, itemComment, owner, name, Typ, strings.Join(items, ", "),
		strings.Join(itemsByteSize, ""), strings.Join(itemsSerialize, ""),
		name, Typ, itemComment, owner, name, Typ, strings.Join(items, ", "),
		l, _MAKE_CUTOFF_STR(strconv.Itoa(len(ts))), switchLabel, _TAG_FIELD_STR("tag"),
		strings.Join(itemsDeserialize, "")), updateTygo(pkgs)
}

func (t *Object) Go() (string, map[string]string) {
	var fields []string
	pkgs := map[string]string{TYGO_PATH: ""}
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
	for _, method := range t.Methods {
		param_s, param_p := typeListGo(fmt.Sprintf("(s *%s) ", t.Name), method.Name, "param", method.Params)
		result_s, result_p := typeListGo(fmt.Sprintf("(s *%s) ", t.Name), method.Name, "result", method.Params)
		pkgs = update(pkgs, param_p)
		pkgs = update(pkgs, result_p)
		methods = append(methods, param_s)
		methods = append(methods, result_s)
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
	s.SetCachedSize(size)
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
