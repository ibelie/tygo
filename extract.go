// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"log"
	"path"
	"strings"

	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
)

type Extracter func(string, string, string, []Type)

func Extract(dir string, extracter Extracter) (types []Type) {
	buildPackage, err := build.Import(dir, "", build.ImportComment)
	if err != nil {
		log.Fatalf("[Tygo][Extract] Cannot import package:\n>>>> %v", err)
		return
	}
	fs := token.NewFileSet()
	for _, filename := range buildPackage.GoFiles {
		file, err := parser.ParseFile(fs, path.Join(buildPackage.Dir, filename), nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("[Tygo][Extract] Cannot parse file:\n>>>> %v", err)
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
				imports, typePkg := extractPkgs(file)
				var ts []Type
				if strings.TrimSpace(decl.Doc.Text()) != "" {
					ts = Parse(decl.Doc.Text(), imports, typePkg)
					types = append(types, ts...)
				}
				if extracter != nil {
					extracter(dir, filename, file.Name.Name, ts)
				}
			}
		}
	}
	return
}

func extractPkgs(file *ast.File) (map[string]string, map[string][2]string) {
	imports := make(map[string]string)
	typePkg := make(map[string][2]string)
	for _, importSpec := range file.Imports {
		pkg := strings.Trim(importSpec.Path.Value, "\"")
		if importSpec.Name == nil {
			if p, err := build.Import(pkg, "", build.AllowBinary); err != nil {
				log.Fatalf("[Tygo][Inject] Cannot import package:\n>>>> %v", err)
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
	return imports, typePkg
}
