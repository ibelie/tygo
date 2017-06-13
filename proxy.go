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
	"strings"

	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
)

const (
	goHeader = `// Generated by tygo.  DO NOT EDIT!

package %s

`
	goImport = `import %s"%s"
`
)

var SRC_PATH = os.Getenv("GOPATH") + "/src/"

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

func packageDoc(path string) *doc.Package {
	p, err := build.Import(path, "", build.ImportComment)
	if err != nil {
		return nil
	}
	fs := token.NewFileSet()
	include := func(info os.FileInfo) bool {
		for _, name := range p.GoFiles {
			if name == info.Name() {
				return true
			}
		}
		return false
	}

	if pkgs, err := parser.ParseDir(fs, p.Dir, include, parser.ParseComments); err != nil || len(pkgs) != 1 {
		return nil
	} else {
		return doc.New(pkgs[p.Name], p.ImportPath, doc.AllDecls)
	}
}

func inject(filename string, doc string, file *ast.File) {
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
	var values []string
	var names []string
	for _, name := range t.Sorted() {
		values = append(values, fmt.Sprintf(`
	%s_%s %s%s = %d`, t.Name, name, strings.Repeat(" ", t.nameMax-len(name)), t.Name, t.Values[name]))
		names = append(names, fmt.Sprintf(`
	case %s_%s:
		return "%s"`, t.Name, name, name))
	}
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

func (i %s) ByteSize() int {
	return tygo.SizeVarint(uint64(i))
}

func (i %s) Serialize(output *tygo.ProtoBuf) {
	output.WriteUvarint(uint64(i))
}

func (i *%s) Deserialize(input *tygo.ProtoBuf) (err error) {
	x, err := input.ReadUvarint()
	*i = %s(x)
	return
}
`, t.Name, strings.Join(values, ""), t.Name, strings.Join(names, ""),
		t.Name, t.Name, t.Name, t.Name, t.Name), map[string]string{"fmt": "", TYGO_PATH: ""}
}

func (t *Method) Go() (string, map[string]string) {
	var s string
	var pkgs map[string]string
	var params []string
	var paramsComment []string
	for i, param := range t.Params {
		param_s, param_p := param.Go()
		pkgs = update(pkgs, param_p)
		params = append(params, fmt.Sprintf("a%d %s", i, param_s))
		paramsComment = append(paramsComment, fmt.Sprintf("a%d: %s", i, param))
	}
	paramComment := strings.Join(paramsComment, ", ")
	if params != nil {
		s += fmt.Sprintf(`
// %s Params(%s)
func Serialize%sParam(%s) (data []byte) {
	size := 0
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	return
}

// %s Params(%s)
func Deserialize%sParam(data []byte) (%s, err error) {
	return
}
`, t.Name, paramComment, t.Name, strings.Join(params, ", "),
			t.Name, paramComment, t.Name, strings.Join(params, ", "))
	}

	var results []string
	var resultsComment []string
	for i, result := range t.Results {
		result_s, result_p := result.Go()
		pkgs = update(pkgs, result_p)
		results = append(results, fmt.Sprintf("a%d %s", i, result_s))
		resultsComment = append(resultsComment, fmt.Sprintf("a%d: %s", i, result))
	}
	resultComment := strings.Join(resultsComment, ", ")
	if results != nil {
		s += fmt.Sprintf(`
// %s Results(%s)
func Serialize%sResult(%s) (data []byte) {
	size := 0
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	return
}

// %s Results(%s)
func Deserialize%sResult(data []byte) (%s, err error) {
	return
}
`, t.Name, resultComment, t.Name, strings.Join(results, ", "),
			t.Name, resultComment, t.Name, strings.Join(results, ", "))
	}
	return s, pkgs
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
		method_s, method_p := method.Go()
		pkgs = update(pkgs, method_p)
		method_s = strings.Replace(method_s, "func Serialize", fmt.Sprintf("func (s *%s) Serialize", t.Name), -1)
		method_s = strings.Replace(method_s, "func Deserialize", fmt.Sprintf("func (s *%s) Deserialize", t.Name), -1)
		methods = append(methods, method_s)
	}

	return fmt.Sprintf(`
type %s struct {%s
}

func (s *%s) ByteSize() (size int) {
	return
}

func (s *%s) Serialize(output *tygo.ProtoBuf) {
}

func (s *%s) Deserialize(input *tygo.ProtoBuf) (err error) {
	return
}
%s`, t.Name, strings.Join(fields, ""), t.Name, t.Name, t.Name, strings.Join(methods, "")), pkgs
}

func (t SimpleType) Go() (string, map[string]string) {
	if string(t) == "bytes" {
		return "[]byte", nil
	}
	return string(t), nil
}

func (t *ObjectType) Go() (string, map[string]string) {
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
