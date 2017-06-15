// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"os"
	"strings"

	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
)

func pow10(x uint) uint {
	if x == 0 {
		return 1
	} else {
		return pow10(x-1) * 10
	}
}

func update(a map[string]string, b map[string]string) map[string]string {
	if b == nil {
		return a
	} else if a == nil {
		return b
	}
	for k, v := range b {
		a[k] = v
	}
	return a
}

func updateTygo(a map[string]string) map[string]string {
	if a == nil {
		return map[string]string{TYGO_PATH: ""}
	}
	a[TYGO_PATH] = ""
	return a
}

func addIndent(text string, indent int) string {
	if indent <= 0 || strings.TrimSpace(text) == "" {
		return text
	}
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) != "" {
			lines[i] = strings.Repeat("\t", indent) + line
		}
	}
	return strings.Join(lines, "\n")
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
