// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"os"
	"log"
	"strings"

	"crypto/md5"
	"encoding/base64"

	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
)

func ObjectMap(types []Type) (objects map[string]*Object){
	objects = make(map[string]*Object)
	for _, t := range types {
		if object, ok := t.(*Object); ok {
			if o, exist := objects[object.Name]; exist {
				log.Fatalf("[Tygo] Object already exists: %v %v", o, object)
			}
			objects[object.Name] = object
		}
	}
	return
}

func shortName(name string) string {
	if len(name) > 24 {
		bytes := md5.Sum([]byte(name))
		name = strings.Replace(base64.RawURLEncoding.EncodeToString(bytes[:]), "-", "__", -1)
	}
	return name
}

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
		a = make(map[string]string)
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
