// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"

	. "tygo"
)

func main() {
	path := flag.String("path", "", "target file")
	flag.Parse()
	fmt.Println(path)
}
