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
	if g1 == g2 {
		return
	} else if g1 == nil || g2 == nil {
		t.Errorf("GoType %v %v", g1, g2)
	} else if g1.PP != g2.PP {
		t.Errorf("GoType.PP: %v %v", g1.PP, g2.PP)
	} else if g1.AP != g2.AP {
		t.Errorf("GoType.AP: %v %v", g1.AP, g2.AP)
	}
}

func compareVector2(t *testing.T, v1 *Vector2, v2 *Vector2) {
	if v1 == v2 {
		return
	} else if v1 == nil || v2 == nil {
		t.Errorf("Vector2 %v %v", v1, v2)
	} else if v1.X != v2.X {
		t.Errorf("Vector2.X: %v %v", v1.X, v2.X)
	} else if v1.Y != v2.Y {
		t.Errorf("Vector2.Y: %v %v", v1.Y, v2.Y)
	} else if v1.S != v2.S {
		t.Errorf("Vector2.S: %v %v", v1.S, v2.S)
	} else if bytes.Compare(v1.B, v2.B) != 0 {
		t.Errorf("Vector2.B: %v %v", v1.B, v2.B)
	} else if v1.E != v2.E {
		t.Errorf("Vector2.E: %v %v", v1.E, v2.E)
	} else {
		compareGoType(t, v1.P, v2.P)
	}
}

func compareFighter(t *testing.T, f1 *Fighter, f2 *Fighter) {
	if f1 == f2 {
		return
	} else if f1 == nil || f2 == nil {
		t.Errorf("Fighter %v %v", f1, f2)
		return
	}
	compareVector2(t, f1.Pos, f2.Pos)
	if f1.IsAwake != f2.IsAwake {
		t.Errorf("Fighter.IsAwake: %v %v", f1.IsAwake, f2.IsAwake)
	}
	if f1.Hp != f2.Hp {
		t.Errorf("Fighter.Hp: %v %v", f1.Hp, f2.Hp)
	}
	if len(f1.Poss) != len(f2.Poss) {
		t.Errorf("Fighter.Poss: %v %v", f1.Poss, f2.Poss)
	} else {
		for k, v1 := range f1.Poss {
			if v2, ok := f2.Poss[k]; !ok {
				t.Errorf("Fighter.Poss: %v %v %v", k, v1, v2)
			} else {
				compareVector2(t, v1, v2)
			}
		}
	}
	if len(f1.Posi) != len(f2.Posi) {
		t.Errorf("Fighter.Posi: %v %v", f1.Posi, f2.Posi)
	} else {
		for k, v1 := range f1.Posi {
			if v2, ok := f2.Posi[k]; !ok || v1 != v2 {
				t.Errorf("Fighter.Posi: %v %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Posl) != len(f2.Posl) {
		t.Errorf("Fighter.Posl: %v %v", f1.Posl, f2.Posl)
	} else {
		for k, v1 := range f1.Posl {
			v2 := f2.Posl[k]
			compareVector2(t, v1, v2)
		}
	}
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
	compareFighter(t, fighter, fighter2)
}
