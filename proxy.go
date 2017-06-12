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
	"go/parser"
	"go/token"
	"io/ioutil"
)

const (
	goHeader = `// Generated by tygo.  DO NOT EDIT!

package %s

`
	goImport = `import %s "%s"
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

const (
	goEnum = `
type %s int

const (%s
)

func (i %s) String() {
	switch i {%s
	default:
		log.Fatalf("[Tygo][%s] Unexpect enum value: %%d", i)
	}
}
`
	goObject = `
type %s struct {%s
}
`
)

func inject(filename string, doc string, file *ast.File) {
	imports := make(map[string]string)
	for _, importSpec := range file.Imports {
		pkg := strings.Trim(importSpec.Path.Value, "\"")
		if importSpec.Name == nil {
			p := strings.Split(pkg, "/")
			imports[p[len(p)-1]] = pkg
		} else {
			imports[importSpec.Name.Name] = pkg
		}
	}
	enums, objects := Parse(doc, imports)

	var head bytes.Buffer
	var body bytes.Buffer
	head.Write([]byte(fmt.Sprintf(goHeader, file.Name)))

	for _, enum := range enums {
		var values []string
		var names []string
		for _, name := range enum.Sorted() {
			values = append(values, fmt.Sprintf(`
	%s_%s %s%s = %d`, enum.Name, name, strings.Repeat(" ", enum.NameMax()-len(name)), enum.Name, enum.Values[name]))
			names = append(names, fmt.Sprintf(`
	case %s_%s:
		return "%s"`, enum.Name, name, name))
		}
		body.Write([]byte(fmt.Sprintf(goEnum, enum.Name, strings.Join(values, ""), enum.Name, strings.Join(names, ""), enum.Name)))
	}

	imported := map[string]bool{}
	for _, object := range objects {
		var fields []string
		var sortedField []string
		nameMax := 0
		for name, _ := range object.Fields {
			if nameMax < len(name) {
				nameMax = len(name)
			}
			sortedField = append(sortedField, name)
		}
		sort.Strings(sortedField)
		var sortedParent []string
		for _, parent := range object.Parents {
			spec, pkgs := parent.Go()
			sortedParent = append(sortedParent, spec)
			for _, pkg := range pkgs {
				if i, ok := imported[pkg[0]]; !ok || !i {
					head.Write([]byte(fmt.Sprintf(goImport, pkg[0], pkg[1])))
					imported[pkg[0]] = true
				}
			}
		}
		sort.Strings(sortedParent)
		for _, parent := range sortedParent {
			fields = append(fields, fmt.Sprintf(`
	%s`, parent))
		}
		for _, name := range sortedField {
			spec, pkgs := object.Fields[name].Go()
			for _, pkg := range pkgs {
				if i, ok := imported[pkg[0]]; !ok || !i {
					head.Write([]byte(fmt.Sprintf(goImport, pkg[0], pkg[1])))
					imported[pkg[0]] = true
				}
			}
			fields = append(fields, fmt.Sprintf(`
	%s %s%s`, name, strings.Repeat(" ", nameMax-len(name)), spec))
		}
		body.Write([]byte(fmt.Sprintf(goObject, object.Name, strings.Join(fields, ""))))
	}

	head.Write(body.Bytes())
	ioutil.WriteFile(filename, head.Bytes(), 0666)
}
