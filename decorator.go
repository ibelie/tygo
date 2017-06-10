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

	"go/ast"
	"go/build"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
)

type Decorator struct {
	Name   string
	Params []*Decorator
}

type Decorators map[string]*Decorator

func (s Decorators) CheckDecorator(t string) bool {
	_, ok := s[t]
	return ok
}

func (s Decorators) sorted() (decorators []string) {
	for d := range s {
		decorators = append(decorators, d)
	}
	sort.Strings(decorators)
	return
}

type Parent struct {
	Type
	Decorators
	Tag string
}

type Parents map[string]*Parent

func (s Parents) CheckParent(t string) bool {
	_, ok := s[t]
	return ok
}

func (s Parents) sorted() (parents []string) {
	for parent := range s {
		parents = append(parents, parent)
	}
	sort.Strings(parents)
	return
}

type Field struct {
	Type
	Decorators
	Tag  string
	Name string
}

type Fields map[string]*Field

func (s Fields) sorted() (fields []string) {
	for field := range s {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return
}

type Method struct {
	Decorators
	Name   string
	Params []Type
	Result []Type
}

type Methods map[string]*Method

func (s Methods) sorted() (methods []string) {
	for method := range s {
		methods = append(methods, method)
	}
	sort.Strings(methods)
	return
}

type Object struct {
	Name string
	Path string
	Parents
	Fields
	Methods
}

func (object *Object) MethodsOf(t string) Methods {
	methods := make(Methods)
	for _, method := range object.Methods {
		if method.CheckDecorator(t) {
			methods[method.Name] = method
		}
	}
	return methods
}

func (object *Object) FieldsOf(t string) Message {
	properties := make(Message)
	for _, field := range object.Fields {
		if field.CheckDecorator(t) {
			properties[field.Name] = field.Type
		}
	}
	return properties
}

type Objects map[string]*Object

func (s Objects) sorted() (objects []string) {
	for object := range s {
		objects = append(objects, object)
	}
	sort.Strings(objects)
	return
}

type Package struct {
	Name    string
	Path    string
	Imports []string
	Objects
}

func parse(path string) (*token.FileSet, *doc.Package) {
	buildPackage, err := build.Import(path, "", build.ImportComment)
	if err != nil {
		log.Printf("[Decorator] Cannot import package:\n>>>>%v", err)
		return nil, nil
	}
	fs := token.NewFileSet()
	// include tells parser.ParseDir which files to include.
	// That means the file must be in the build package's GoFiles or CgoFiles
	// list only (no tag-ignored files, tests, swig or other non-Go files).
	include := func(info os.FileInfo) bool {
		for _, name := range buildPackage.GoFiles {
			if name == info.Name() {
				return true
			}
		}
		return false
	}
	pkgs, err := parser.ParseDir(fs, buildPackage.Dir, include, parser.ParseComments)
	if err != nil {
		panic(fmt.Sprintf("[Tygo][Decorator] Cannot parse package:\n>>>>%v", err))
	}
	// Make sure they are all in one package.
	if len(pkgs) != 1 {
		panic(fmt.Sprintf("[Tygo][Decorator] Multiple packages in directory %s", buildPackage.Dir))
	}
	return fs, doc.New(pkgs[buildPackage.Name], buildPackage.ImportPath, doc.AllDecls)
}

func copyDecorators(from Decorators) Decorators {
	to := make(Decorators)
	for k, v := range from {
		to[k] = v
	}
	return to
}

func processType(fs *token.FileSet, typ ast.Expr) Type {
	switch t := typ.(type) {
	case *ast.Ident:
		return SimpleType(t.Name)
	case *ast.SelectorExpr:
		return SimpleType(t.Sel.Name)
	case *ast.StarExpr:
		return processType(fs, t.X)
	case *ast.ArrayType:
		return &ListType{E: processType(fs, t.Elt)}
	case *ast.MapType:
		return &DictType{K: processType(fs, t.Key), V: processType(fs, t.Value)}
	case *ast.InterfaceType:
	case *ast.FuncType:
	case *ast.StructType:
	case *ast.Ellipsis:
	case *ast.ChanType:
	default:
		ast.Print(fs, typ)
		log.Println("[Decorator] Type warning: ", typ)
	}
	return SimpleType("Unknown")
}

func processFields(fs *token.FileSet, fieldList []*ast.Field) (fields []*Field) {
	for _, field := range fieldList {
		f := &Field{Decorators: make(Decorators)}
		if field.Tag != nil {
			f.Tag = field.Tag.Value
		}
		if field.Comment != nil {
			for _, comment := range field.Comment.List {
				decorate(comment.Text, f.Decorators)
			}
		}
		f.Type = processType(fs, field.Type)
		if field.Names == nil {
			fields = append(fields, f)
			continue
		}
		for _, name := range field.Names {
			fields = append(fields, &Field{
				Name:       name.Name,
				Type:       f.Type,
				Tag:        f.Tag,
				Decorators: copyDecorators(f.Decorators),
			})
		}
	}
	return
}

func Decorate(path string) *Package {
	fs, docPkg := parse(path)
	if fs == nil || docPkg == nil {
		return nil
	}
	pkg := &Package{
		Name:    docPkg.Name,
		Path:    docPkg.ImportPath,
		Imports: docPkg.Imports,
		Objects: make(map[string]*Object),
	}
	for _, obj := range docPkg.Types {
		if typeSpec, ok := obj.Decl.Specs[0].(*ast.TypeSpec); ok {
			object := &Object{
				Path:    path,
				Parents: make(Parents),
				Fields:  make(Fields),
				Methods: make(Methods),
			}
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				object.Name = obj.Name
				for _, field := range processFields(fs, structType.Fields.List) {
					if field.Name == "" {
						object.Parents[field.Type.String()] = &Parent{
							Tag:        field.Tag,
							Type:       field.Type,
							Decorators: field.Decorators,
						}
					} else {
						object.Fields[field.Name] = field
					}
				}
			}
			for _, meth := range obj.Methods {
				m := &Method{
					Name:       meth.Name,
					Decorators: make(Decorators),
				}
				decorate(meth.Doc, m.Decorators)
				for _, param := range processFields(fs, meth.Decl.Type.Params.List) {
					m.Params = append(m.Params, param.Type)
				}
				if meth.Decl.Type.Results != nil {
					for _, result := range processFields(fs, meth.Decl.Type.Results.List) {
						m.Result = append(m.Result, result.Type)
					}
				}
				object.Methods[meth.Name] = m
			}
			pkg.Objects[object.Name] = object
		}
	}
	return pkg
}

var formatBuf bytes.Buffer

func toString(fset *token.FileSet, node interface{}) string {
	formatBuf.Reset()
	err := format.Node(&formatBuf, fset, node)
	if err != nil {
		panic(fmt.Sprintf("[Tygo][Decorator] Cannot format node:\n>>>>%v", err))
	}
	return fmt.Sprintf("%s", formatBuf.Bytes())
}

func decorate(text string, decorators Decorators) {
	parser := &decoratorParserImpl{}
	parser.Parse(&decoratorLex{line: []byte(text)})
	for _, dec := range parser.lval.decorators {
		decorators[dec.Name] = dec
	}
}
