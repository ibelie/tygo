// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/ibelie/tygo"
)

func main() {
	input := flag.String("in", "", "input package")
	output := flag.String("out", "", "output file")
	name := flag.String("name", "", "file name")
	module := flag.String("module", "", "module name")
	flag.Parse()
	tygo.Typescript(*output, *name, *module, tygo.Extract(*input, nil), nil)
}
