// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package test

import (
	"os"
	"reflect"
	"testing"

	"github.com/ibelie/tygo"
)

func TestTypescript(t *testing.T) {
	TEST_PATH := os.Getenv("TYPESCRIPT_PATH") + "/tyts/test"
	THIS_PATH := reflect.TypeOf(GoType{}).PkgPath()
	types := tygo.Extract(THIS_PATH, nil)
	tygo.Typescript(TEST_PATH, "types", types, nil)
}
