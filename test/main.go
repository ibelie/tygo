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
	vd.Reset()
	v3 := &Vector2{}
	v3.Deserialize(vd)
	fmt.Println(v3)
	fmt.Println(len(vd.Buffer), vd.Buffer)
	fd := &tygo.ProtoBuf{Buffer: make([]byte, fighter.ByteSize())}
	fighter.Serialize(fd)
	fmt.Println(len(fd.Buffer))
}
