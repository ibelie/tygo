// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/ibelie/tygo"
)

func main() {
	vd := &tygo.ProtoBuf{Buffer: make([]byte, v.ByteSize())}
	v.Serialize(vd)
	fmt.Println(len(vd.Buffer), vd.Buffer)
	fd := &tygo.ProtoBuf{Buffer: make([]byte, fighter.ByteSize())}
	fighter.Serialize(fd)
	fd.Reset()
	fighter2 := &Fighter{}
	fighter2.Deserialize(fd)
	fmt.Println(fighter2)
	fmt.Println(len(fd.Buffer))
}
