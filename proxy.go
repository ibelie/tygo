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
				if !ok || spec.Path.Value != "\"github.com/ibelie/tygo\"" {
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

func (t *Object) Go() (string, [][2]string) {
	pkgs := [][2]string{[2]string{"", "io"}}
	var fields []string
	var sortedParent []string
	for _, parent := range t.Parents {
		s, p := parent.Go()
		sortedParent = append(sortedParent, s)
		pkgs = append(pkgs, p...)
	}
	sort.Strings(sortedParent)
	for _, parent := range sortedParent {
		fields = append(fields, fmt.Sprintf(`
	%s`, parent))
	}

	nameMax := 0
	typeMax := 0
	var sortedField []string
	fieldMap := make(map[string][2]string)
	for name, field := range t.Fields {
		s, p := field.Go()
		pkgs = append(pkgs, p...)
		if nameMax < len(name) {
			nameMax = len(name)
		}
		if typeMax < len(s) {
			typeMax = len(s)
		}
		fieldMap[name] = [2]string{s, field.String()}
		sortedField = append(sortedField, name)
	}
	sort.Strings(sortedField)
	for _, name := range sortedField {
		f := fieldMap[name]
		fields = append(fields, fmt.Sprintf(`
	%s %s%s %s// %s`, name, strings.Repeat(" ", nameMax-len(name)),
			f[0], strings.Repeat(" ", typeMax-len(f[0])), f[1]))
	}

	var methods []string
	for _, method := range t.Methods {
		var params []string
		for i, param := range method.Params {
			s, p := param.Go()
			pkgs = append(pkgs, p...)
			params = append(params, fmt.Sprintf("a%d %s", i, s))
		}
		if params != nil {
			methods = append(methods, fmt.Sprintf(`
func (s *%s) Serialize%sParam(%s) (data string, err error) {
	return
}

func (s *%s) Deserialize%sParam(data string) (%s, err error) {
	return
}
`, t.Name, method.Name, strings.Join(params, ", "), t.Name, method.Name, strings.Join(params, ", ")))
		}

		var results []string
		for i, result := range method.Results {
			s, p := result.Go()
			pkgs = append(pkgs, p...)
			results = append(results, fmt.Sprintf("a%d %s", i, s))
		}
		if results != nil {
			methods = append(methods, fmt.Sprintf(`
func (s *%s) Serialize%sResult(%s) (data string, err error) {
	return
}

func (s *%s) Deserialize%sResult(data string) (%s, err error) {
	return
}
`, t.Name, method.Name, strings.Join(results, ", "), t.Name, method.Name, strings.Join(results, ", ")))
		}
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

func (s *%s) ByteSize() (int, error) {
	return 0, nil
}

func (s *%s) Serialize(w io.Writer) error {
	return nil
}

func (s *%s) Deserialize(r io.Reader) error {
	return nil
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
