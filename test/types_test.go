// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package test

import (
	"bytes"
	"github.com/ibelie/tygo"
	"runtime/debug"
	"testing"
)

func compareGoType(t *testing.T, g1 *GoType, g2 *GoType) {
	if g1 == g2 {
		return
	} else if g1 == nil || g2 == nil {
		debug.PrintStack()
		t.Errorf("GoType %v %v", g1, g2)
	} else if g1.PP != g2.PP {
		debug.PrintStack()
		t.Errorf("GoType.PP: %v %v", g1.PP, g2.PP)
	} else if g1.AP != g2.AP {
		debug.PrintStack()
		t.Errorf("GoType.AP: %v %v", g1.AP, g2.AP)
	}
}

func compareVector2(t *testing.T, v1 *Vector2, v2 *Vector2) {
	if v1 == v2 {
		return
	} else if v1 == nil || v2 == nil {
		debug.PrintStack()
		t.Errorf("Vector2 %v %v", v1, v2)
	} else if v1.X != v2.X {
		debug.PrintStack()
		t.Errorf("Vector2.X: %v %v", v1.X, v2.X)
	} else if v1.Y != v2.Y {
		debug.PrintStack()
		t.Errorf("Vector2.Y: %v %v", v1.Y, v2.Y)
	} else if v1.S != v2.S {
		debug.PrintStack()
		t.Errorf("Vector2.S: %v %v", v1.S, v2.S)
	} else if bytes.Compare(v1.B, v2.B) != 0 {
		debug.PrintStack()
		t.Errorf("Vector2.B: %v %v", v1.B, v2.B)
	} else if v1.E != v2.E {
		debug.PrintStack()
		t.Errorf("Vector2.E: %v %v", v1.E, v2.E)
	} else {
		compareGoType(t, v1.P, v2.P)
	}
}

func compareFighter_Part1(t *testing.T, f1 *Fighter_Part1, f2 *Fighter_Part1) {
	if f1 == f2 {
		return
	} else if f1 == nil || f2 == nil {
		debug.PrintStack()
		t.Errorf("Fighter_Part1 %v %v", f1, f2)
		return
	}
	compareVector2(t, f1.Pos, f2.Pos)
	if f1.IsAwake != f2.IsAwake {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.IsAwake: %v %v", f1.IsAwake, f2.IsAwake)
	}
	if f1.Hp != f2.Hp {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Hp: %v %v", f1.Hp, f2.Hp)
	}
	if len(f1.Poss) != len(f2.Poss) {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Poss: %v %v", f1.Poss, f2.Poss)
	} else {
		for k, v1 := range f1.Poss {
			if v2, ok := f2.Poss[k]; !ok {
				debug.PrintStack()
				t.Errorf("Fighter_Part1.Poss: %v %v %v", k, v1, v2)
			} else {
				compareVector2(t, v1, v2)
			}
		}
	}
	if len(f1.Posi) != len(f2.Posi) {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Posi: %v %v", f1.Posi, f2.Posi)
	} else {
		for k, v1 := range f1.Posi {
			if v2, ok := f2.Posi[k]; !ok || v1 != v2 {
				debug.PrintStack()
				t.Errorf("Fighter_Part1.Posi: %v %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Posl) != len(f2.Posl) {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Posl: %v %v", f1.Posl, f2.Posl)
	} else {
		for k, v1 := range f1.Posl {
			v2 := f2.Posl[k]
			compareVector2(t, v1, v2)
		}
	}
	if len(f1.Posll) != len(f2.Posll) {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Posll: %v %v", f1.Posll, f2.Posll)
	} else {
		for k1, l1 := range f1.Posll {
			l2 := f2.Posll[k1]
			if len(l1) != len(l2) {
				debug.PrintStack()
				t.Errorf("Fighter_Part1.Posll[%v]: %v %v", k1, l1, l2)
			} else {
				for k2, v1 := range l1 {
					v2 := l2[k2]
					compareVector2(t, v1, v2)
				}
			}
		}
	}
	if len(f1.Pyl) != len(f2.Pyl) {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Pyl: %v %v", f1.Pyl, f2.Pyl)
	} else {
		for k, v1 := range f1.Pyl {
			v2 := f2.Pyl[k]
			compareGoType(t, v1, v2)
		}
	}
	if len(f1.Pyd) != len(f2.Pyd) {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Pyd: %v %v", f1.Pyd, f2.Pyd)
	} else {
		for k, v1 := range f1.Pyd {
			if v2, ok := f2.Pyd[k]; !ok {
				debug.PrintStack()
				t.Errorf("Fighter_Part1.Pyd: %v %v %v", k, v1, v2)
			} else {
				compareGoType(t, v1, v2)
			}
		}
	}
	if int32(f1.Pyv1.(int)) != f2.Pyv1.(int32) {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Pyv1: %v %v", f1.Pyv1, f2.Pyv1)
	}
	compareGoType(t, f1.Pyv2.(*GoType), f2.Pyv2.(*GoType))
}

func compareFighter_Part2(t *testing.T, f1 *Fighter_Part2, f2 *Fighter_Part2) {
	if f1 == f2 {
		return
	} else if f1 == nil || f2 == nil {
		debug.PrintStack()
		t.Errorf("Fighter_Part2 %v %v", f1, f2)
		return
	}
	compareFighter_Part1(t, &f1.Fighter_Part1, &f2.Fighter_Part1)
	if len(f1.Fl) != len(f2.Fl) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Fl: %v %v", f1.Fl, f2.Fl)
	} else {
		for k, v1 := range f1.Fl {
			v2 := f2.Fl[k]
			if v1 != v2 {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Fl: %v %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Bl) != len(f2.Bl) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Bl: %v %v", f1.Bl, f2.Bl)
	} else {
		for k, v1 := range f1.Bl {
			v2 := f2.Bl[k]
			if bytes.Compare(v1, v2) != 0 {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Bl: %v %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Sl) != len(f2.Sl) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Sl: %v %v", f1.Sl, f2.Sl)
	} else {
		for k, v1 := range f1.Sl {
			v2 := f2.Sl[k]
			if v1 != v2 {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Sl: %v %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Bd) != len(f2.Bd) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Bd: %v %v", f1.Bd, f2.Bd)
	} else {
		for k, v1 := range f1.Bd {
			if v2, ok := f2.Bd[k]; !ok || bytes.Compare(v1, v2) != 0 {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Bd: %v %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Sd) != len(f2.Sd) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Sd: %v %v", f1.Sd, f2.Sd)
	} else {
		for k, v1 := range f1.Sd {
			if v2, ok := f2.Sd[k]; !ok || v1 != v2 {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Sd: %v %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.El) != len(f2.El) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.El: %v %v", f1.El, f2.El)
	} else {
		for k, v1 := range f1.El {
			v2 := f2.El[k]
			if v1 != v2 {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.El: %v %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Ed) != len(f2.Ed) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Ed: %v %v", f1.Ed, f2.Ed)
	} else {
		for k, v1 := range f1.Ed {
			if v2, ok := f2.Ed[k]; !ok || v1 != v2 {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Ed: %v %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Ll) != len(f2.Ll) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Ll: %v %v", f1.Ll, f2.Ll)
	} else {
		for k1, l1 := range f1.Ll {
			l2 := f2.Ll[k1]
			if len(l1) != len(l2) {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Ll[%v]: %v %v", k1, l1, l2)
			} else {
				for k2, v1 := range l1 {
					v2 := l2[k2]
					if v1 != v2 {
						debug.PrintStack()
						t.Errorf("Fighter_Part2.Ll[%v]: %v %v %v", k1, k2, v1, v2)
					}
				}
			}
		}
	}
}

func compareFighter(t *testing.T, f1 *Fighter, f2 *Fighter) {
	if f1 == f2 {
		return
	} else if f1 == nil || f2 == nil {
		debug.PrintStack()
		t.Errorf("Fighter %v %v", f1, f2)
		return
	}
	compareFighter_Part2(t, &f1.Fighter_Part2, &f2.Fighter_Part2)
	if f1.V0 != f2.V0 {
		debug.PrintStack()
		t.Errorf("Fighter.V0: %v %v", f1.V0, f2.V0)
	}
	if int32(f1.V1.(int)) != f2.V1.(int32) {
		debug.PrintStack()
		t.Errorf("Fighter.V1: %v %v", f1.V1, f2.V1)
	}
	if bytes.Compare(f1.V2.([]byte), f2.V2.([]byte)) != 0 {
		debug.PrintStack()
		t.Errorf("Fighter.V2: %v %v", f1.V1, f2.V1)
	}
	compareVector2(t, f1.V3.(*Vector2), f2.V3.(*Vector2))
	if float32(f1.V4.(float64)) != f2.V4.(float32) {
		debug.PrintStack()
		t.Errorf("Fighter.V4: %v %v", f1.V4, f2.V4)
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
