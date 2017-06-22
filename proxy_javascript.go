// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path"
	"sort"
	"strings"

	"io/ioutil"
)

func Javascript(dir string, name string, types []Type) {
	var buffer bytes.Buffer
	var sortedObjects []string
	objects := make(map[string]*Object)
	for _, t := range types {
		if object, ok := t.(*Object); ok {
			if o, exist := objects[object.Name]; exist {
				log.Fatalf("[Tygo][Javascript] Object already exists: %v %v", o, object)
			}
			objects[object.Name] = object
			sortedObjects = append(sortedObjects, object.Name)
		}
	}
	sort.Strings(sortedObjects)

	genTypes := make(map[string]Type)
	for _, name := range sortedObjects {
		objects[name].Javascript(&buffer, genTypes, objects)
	}

	ioutil.WriteFile(path.Join(dir, name+".js"), buffer.Bytes(), 0666)
}

func (t *Enum) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}

func (t *Method) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}

func (t *Object) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	if _, exist := types[t.Name]; exist {
		return WireVarint, ""
	}
	fields := t.AllFields(objects)
	var codes []string
	for i, field := range fields {
		wiretype, typename := field.Javascript(writer, types, objects)
		codes = append(codes, fmt.Sprintf(`
	{name: %s, tag: %d, tagsize: %d, wiretype: %d, type: %s}`, field.Name, _MAKE_TAG(i+1, wiretype), TAG_SIZE(i+1), wiretype, typename))
	}
	types[t.Name] = t
	writer.Write([]byte(fmt.Sprintf(`
%s = tyts.Object('%s', %d, [%s
]);`, t.Name, t.Name, _MAKE_CUTOFF(len(fields)), strings.Join(codes, ","))))
	return WireVarint, ""
}

func (t UnknownType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}

func (t SimpleType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}

func (t *EnumType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}

func (t *InstanceType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}

func (t *FixedPointType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}

func (t *ListType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}

func (t *DictType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}

func (t *VariantType) Javascript(writer io.Writer, types map[string]Type, objects map[string]*Object) (WireType, string) {
	return WireVarint, ""
}
