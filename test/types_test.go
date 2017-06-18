// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package test

import (
	"bytes"
	"testing"

	"github.com/ibelie/tygo"
)

func compareGoType(t *testing.T, g1 *GoType, g2 *GoType) {
	if g1.PP != g2.PP {
		t.Errorf("GoType.PP: %v %v", g1.PP, g2.PP)
	}
	if g1.AP != g2.AP {
		t.Errorf("GoType.AP: %v %v", g1.AP, g2.AP)
	}
}

func compareVector2(t *testing.T, v1 *Vector2, v2 *Vector2) {
	if v1.X != v2.X {
		t.Errorf("Vector2.X: %v %v", v1.X, v2.X)
	}
	if v1.Y != v2.Y {
		t.Errorf("Vector2.Y: %v %v", v1.Y, v2.Y)
	}
	if v1.S != v2.S {
		t.Errorf("Vector2.S: %v %v", v1.S, v2.S)
	}
	if bytes.Compare(v1.B, v2.B) != 0 {
		t.Errorf("Vector2.B: %v %v", v1.B, v2.B)
	}
	if v1.E != v2.E {
		t.Errorf("Vector2.E: %v %v", v1.E, v2.E)
	}
	compareGoType(t, v1.P, v2.P)
}

func TestVector2(t *testing.T) {
	vd := &tygo.ProtoBuf{Buffer: make([]byte, v.ByteSize())}
	v.Serialize(vd)
	vd.Reset()
	v3 := &Vector2{}
	v3.Deserialize(vd)
	compareVector2(t, v, v3)
}

func TestFighter(t *testing.T) {
	fd := &tygo.ProtoBuf{Buffer: make([]byte, fighter.ByteSize())}
	fighter.Serialize(fd)
	fd.Reset()
	fighter2 := &Fighter{}
	fighter2.Deserialize(fd)
	t.Log(fighter2)
	t.Log(len(fd.Buffer))
}
