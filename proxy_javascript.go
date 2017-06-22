// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"bytes"
	"path"

	"io/ioutil"
)

func Javascript(dir string, name string, types []Type) {
	var head bytes.Buffer
	var body bytes.Buffer

	head.Write(body.Bytes())
	ioutil.WriteFile(path.Join(dir, name+".js"), head.Bytes(), 0666)
}

func (t *Enum) Javascript() (string, []string) {
	return "", nil
}

func (t *Method) Javascript() (string, []string) {
	return "", nil
}

func (t *Object) Javascript() (string, []string) {
	return "", nil
}

func (t UnknownType) Javascript() (string, []string) {
	return "", nil
}

func (t SimpleType) Javascript() (string, []string) {
	return "", nil
}

func (t *EnumType) Javascript() (string, []string) {
	return "", nil
}

func (t *InstanceType) Javascript() (string, []string) {
	return "", nil
}

func (t *FixedPointType) Javascript() (string, []string) {
	return "", nil
}

func (t *ListType) Javascript() (string, []string) {
	return "", nil
}

func (t *DictType) Javascript() (string, []string) {
	return "", nil
}

func (t *VariantType) Javascript() (string, []string) {
	return "", nil
}
