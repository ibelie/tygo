// Copyright 2017 - 2018 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"log"
	"os"
	"strings"

	"crypto/md5"
	"encoding/base64"

	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
)

const (
	STR_ZERO        = "0"
	STR_NIL         = "nil"
	STR_INT32       = "int32"
	STR_UINT32      = "uint32"
	STR_INT64       = "int64"
	STR_UINT64      = "uint64"
	STR_BYTES       = "bytes"
	STR_STRING      = "string"
	STR_SYMBOL      = "symbol"
	STR_BOOL        = "bool"
	STR_FLOAT32     = "float32"
	STR_FLOAT64     = "float64"
	STR_UNKNOWN     = "unknown"
	STR_PREFIELDNUM = "preFieldNum"
)

func ObjectMap(types []Type, useFullName bool) (objects map[string]*Object) {
	objects = make(map[string]*Object)
	for _, t := range types {
		if object, ok := t.(*Object); ok {
			fullName := object.Name
			if useFullName {
				fullName = object.FullName()
			}
			if o, exist := objects[fullName]; exist {
				log.Fatalf("[Tygo] Object already exists: %v %v", o, object)
			}
			objects[fullName] = object
		}
	}
	return
}

func PkgTypeMap(types []Type) (pkgs map[string][]Type) {
	pkgs = make(map[string][]Type)
	for _, t := range types {
		if enum, ok := t.(*Enum); ok {
			pkgs[enum.Package] = append(pkgs[enum.Package], enum)
		} else if object, ok := t.(*Object); ok {
			pkgs[object.Package] = append(pkgs[object.Package], object)
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
