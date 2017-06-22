// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"path"

	"io/ioutil"
)

func Typescript(dir string, name string, types []Type) {
	var buffer bytes.Buffer

	ioutil.WriteFile(path.Join(dir, name+".d.ts"), buffer.Bytes(), 0666)
	Javascript(dir, name, types)
}

func (t *Enum) Typescript() string {
	return ""
}

func (t *Method) Typescript() string {
	return ""
}

func (t *Object) Typescript() string {
	return ""
}

func (t UnknownType) Typescript() string {
	return ""
}

func (t SimpleType) Typescript() string {
	return ""
}

func (t *EnumType) Typescript() string {
	return ""
}

func (t *InstanceType) Typescript() string {
	return ""
}

func (t *FixedPointType) Typescript() string {
	return ""
}

func (t *ListType) Typescript() string {
	return ""
}

func (t *DictType) Typescript() string {
	return ""
}

func (t *VariantType) Typescript() string {
	return ""
}
