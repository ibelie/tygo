// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"fmt"
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
		panic(fmt.Sprintf("[Tygo][Inject] Cannot import package:\n>>>>%v", err))
		return
	}
	fs := token.NewFileSet()
	for _, filename := range buildPackage.GoFiles {
		if strings.HasSuffix(filename, ".ty.go") {
			continue
		}
		file, err := parser.ParseFile(fs, buildPackage.Dir+"/"+filename, nil, parser.ParseComments)
		if err != nil {
			panic(fmt.Sprintf("[Tygo][Inject] Cannot parse file:\n>>>>%v", err))
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
				panic(fmt.Sprintf("[Tygo][Inject] Cannot import package:\n>>>>%v", err))
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

	imported := map[string]bool{}
	for _, t := range types {
		code, pkgs := t.Go()
		for _, pkg := range pkgs {
			if i, ok := imported[pkg[1]]; !ok || !i {
				head.Write([]byte(fmt.Sprintf(goImport, pkg[0], pkg[1])))
				imported[pkg[1]] = true
			}
		}
		body.Write([]byte(code))
	}

	head.Write(body.Bytes())
	ioutil.WriteFile(filename, head.Bytes(), 0666)
}

func (t *Enum) Go() (string, [][2]string) {
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
type %s int

const (%s
)

func (i %s) String() string {
	switch i {%s
	default:
		panic(fmt.Sprintf("[Tygo][%s] Unexpect enum value: %%d", i))
		return "UNKNOWN"
	}
}
`, t.Name, strings.Join(values, ""), t.Name, strings.Join(names, ""), t.Name), [][2]string{[2]string{"", "fmt"}}
}

func (t *Method) Go() (string, [][2]string) {
	var s string
	var pkgs [][2]string
	var params []string
	for i, param := range t.Params {
		param_s, param_p := param.Go()
		pkgs = append(pkgs, param_p...)
		params = append(params, fmt.Sprintf("a%d %s", i, param_s))
	}
	if params != nil {
		s += fmt.Sprintf(`
func Serialize%sParam(%s) (data []byte, err error) {
	size := 0
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	return
}

func Deserialize%sParam(data []byte) (%s, err error) {
	return
}
`, t.Name, strings.Join(params, ", "), t.Name, strings.Join(params, ", "))
	}

	var results []string
	for i, result := range t.Results {
		result_s, result_p := result.Go()
		pkgs = append(pkgs, result_p...)
		results = append(results, fmt.Sprintf("a%d %s", i, result_s))
	}
	if results != nil {
		s += fmt.Sprintf(`
func Serialize%sResult(%s) (data []byte, err error) {
	size := 0
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	return
}

func Deserialize%sResult(data []byte) (%s, err error) {
	return
}
`, t.Name, strings.Join(results, ", "), t.Name, strings.Join(results, ", "))
	}
	return s, pkgs
}

func (t *Object) Go() (string, [][2]string) {
	var fields []string
	pkgs := [][2]string{[2]string{"", TYGO_PATH}}
	parent_s, parent_p := t.Parent.Go()
	pkgs = append(pkgs, parent_p...)
	fields = append(fields, fmt.Sprintf(`
	%s`, parent_s))

	nameMax := 0
	typeMax := 0
	var preparedFields [][3]string
	for _, field := range t.Fields {
		field_s, field_p := field.Go()
		pkgs = append(pkgs, field_p...)
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
		pkgs = append(pkgs, method_p...)
		method_s = strings.Replace(method_s, "func Serialize", fmt.Sprintf("func (s *%s) Serialize", t.Name), -1)
		method_s = strings.Replace(method_s, "func Deserialize", fmt.Sprintf("func (s *%s) Deserialize", t.Name), -1)
		methods = append(methods, method_s)
	}

	pkgDict := make(map[string]string)
	var sortedPkg []string
	for _, pkg := range pkgs {
		sortedPkg = append(sortedPkg, pkg[1])
		pkgDict[pkg[1]] = pkg[0]
	}
	sort.Strings(sortedPkg)
	pkgs = nil
	for _, pkg := range sortedPkg {
		pkgs = append(pkgs, [2]string{pkgDict[pkg], pkg})
	}

	return fmt.Sprintf(`
type %s struct {%s
}

func (s *%s) ByteSize() (size int, err error) {
	return
}

func (s *%s) Serialize(output *tygo.ProtoBuf) (err error) {
	return
}

func (s *%s) Deserialize(input *tygo.ProtoBuf) (err error) {
	return
}
%s`, t.Name, strings.Join(fields, ""), t.Name, t.Name, t.Name, strings.Join(methods, "")), pkgs
}

func (t SimpleType) Go() (string, [][2]string) {
	if string(t) == "bytes" {
		return "[]byte", nil
	}
	return string(t), nil
}

func (t *ObjectType) Go() (string, [][2]string) {
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
		return s + t.PkgName + "." + t.Name, [][2]string{[2]string{a, t.PkgPath}}
	}
}

func (t *FixedPointType) Go() (string, [][2]string) {
	return "float64", nil
}

func (t *ListType) Go() (string, [][2]string) {
	s, p := t.E.Go()
	return fmt.Sprintf("[]%s", s), p
}

func (t *DictType) Go() (string, [][2]string) {
	ks, kp := t.K.Go()
	vs, vp := t.V.Go()
	return fmt.Sprintf("map[%s]%s", ks, vs), append(kp, vp...)
}

func (t *VariantType) Go() (string, [][2]string) {
	var p [][2]string
	for _, vt := range t.Ts {
		_, vp := vt.Go()
		p = append(p, vp...)
	}
	return "interface{}", p
}
