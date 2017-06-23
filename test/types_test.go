// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package test

import (
	"bytes"
	"fmt"
	"github.com/ibelie/tygo"
	"runtime/debug"
	"testing"
)

func compareGoType(t *testing.T, g1 *GoType, g2 *GoType, prefix string) {
	if g1 == g2 {
		return
	} else if g1 == nil || g2 == nil {
		debug.PrintStack()
		t.Errorf("%s GoType %v %v", prefix, g1, g2)
	} else if g1.PP != g2.PP {
		debug.PrintStack()
		t.Errorf("%s GoType.PP: %v %v", prefix, g1.PP, g2.PP)
	} else if g1.AP != g2.AP {
		debug.PrintStack()
		t.Errorf("%s GoType.AP: %v %v", prefix, g1.AP, g2.AP)
	}
}

func compareVector2(t *testing.T, v1 *Vector2, v2 *Vector2, prefix string) {
	if v1 == v2 {
		return
	} else if v1 == nil || v2 == nil {
		debug.PrintStack()
		t.Errorf("%s Vector2 %v %v", prefix, v1, v2)
	} else if v1.X != v2.X {
		debug.PrintStack()
		t.Errorf("%s Vector2.X: %v %v", prefix, v1.X, v2.X)
	} else if v1.Y != v2.Y {
		debug.PrintStack()
		t.Errorf("%s Vector2.Y: %v %v", prefix, v1.Y, v2.Y)
	} else if v1.S != v2.S {
		debug.PrintStack()
		t.Errorf("%s Vector2.S: %v %v", prefix, v1.S, v2.S)
	} else if bytes.Compare(v1.B, v2.B) != 0 {
		debug.PrintStack()
		t.Errorf("%s Vector2.B: %v %v", prefix, v1.B, v2.B)
	} else if v1.E != v2.E {
		debug.PrintStack()
		t.Errorf("%s Vector2.E: %v %v", prefix, v1.E, v2.E)
	} else {
		compareGoType(t, v1.P, v2.P, prefix+".P")
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
	compareVector2(t, f1.Pos, f2.Pos, "Fighter_Part1.Pos")
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
				compareVector2(t, v1, v2, fmt.Sprintf("Fighter_Part1.Poss[%v]", k))
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
			compareVector2(t, v1, v2, fmt.Sprintf("Fighter_Part1.Posl[%v]", k))
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
					compareVector2(t, v1, v2, fmt.Sprintf("Fighter_Part1.Posll[%v][%v]", k1, k2))
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
			compareGoType(t, v1, v2, fmt.Sprintf("Fighter_Part1.Pyl[%v]", k))
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
				compareGoType(t, v1, v2, fmt.Sprintf("Fighter_Part1.Pyd[%v]", k))
			}
		}
	}
	if int32(f1.Pyv1.(int)) != f2.Pyv1.(int32) {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Pyv1: %v %v", f1.Pyv1, f2.Pyv1)
	}
	compareGoType(t, f1.Pyv2.(*GoType), f2.Pyv2.(*GoType), "Fighter_Part1.Pyv1")
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
	compareVector2(t, f1.V3.(*Vector2), f2.V3.(*Vector2), "Fighter.V3")
	if float32(f1.V4.(float64)) != f2.V4.(float32) {
		debug.PrintStack()
		t.Errorf("Fighter.V4: %v %v", f1.V4, f2.V4)
	}
	if len(f1.Vl) != len(f2.Vl) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Vl: %v %v", f1.Vl, f2.Vl)
	} else {
		if int32(f1.Vl[0].(int)) != f2.Vl[0].(int32) {
			debug.PrintStack()
			t.Errorf("Fighter.Vl[0]: %v %v", f1.Vl[0], f2.Vl[0])
		}
		if f1.Vl[1].(string) != f2.Vl[1].(string) {
			debug.PrintStack()
			t.Errorf("Fighter.Vl[1]: %v %v", f1.Vl[1], f2.Vl[1])
		}
		if f1.Vl[2] != nil || f2.Vl[2] != nil {
			debug.PrintStack()
			t.Errorf("Fighter.Vl[2]: %v %v", f1.Vl[2], f2.Vl[2])
		}
		compareVector2(t, f1.Vl[3].(*Vector2), f2.Vl[3].(*Vector2), "Fighter.Vl[3]")
		if f1.Vl[4].(float64) != f2.Vl[4].(float64) {
			debug.PrintStack()
			t.Errorf("Fighter.Vl[4]: %v %v", f1.Vl[4], f2.Vl[4])
		}
		if f1.Vl[5] != nil || f2.Vl[5] != nil {
			debug.PrintStack()
			t.Errorf("Fighter.Vl[5]: %v %v", f1.Vl[5], f2.Vl[5])
		}
	}
	if len(f1.Vd) != len(f2.Vd) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Vd: %v %v", f1.Vd, f2.Vd)
	} else {
		if f1.Vd[0] != nil || f2.Vd[0] != nil {
			debug.PrintStack()
			t.Errorf("Fighter.Vd[0]: %v %v", f1.Vd[0], f2.Vd[0])
		}
		if f1.Vd[12].(Corpus) != f2.Vd[12].(Corpus) {
			debug.PrintStack()
			t.Errorf("Fighter.Vd[12]: %v %v", f1.Vd[12], f2.Vd[12])
		}
		if f1.Vd[23].(string) != f2.Vd[23].(string) {
			debug.PrintStack()
			t.Errorf("Fighter.Vd[23]: %v %v", f1.Vd[23], f2.Vd[23])
		}
		compareVector2(t, f1.Vd[34].(*Vector2), f2.Vd[34].(*Vector2), "Fighter.Vd[34]")
		if f1.Vd[45].(float64) != f2.Vd[45].(float64) {
			debug.PrintStack()
			t.Errorf("Fighter.Vd[45]: %v %v", f1.Vd[45], f2.Vd[45])
		}
	}
	if len(f1.Ld) != len(f2.Ld) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Ld: %v %v", f1.Ld, f2.Ld)
	} else {
		for k1, l1 := range f1.Ld {
			if l2, ok := f2.Ld[k1]; !ok || len(l1) != len(l2) {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Ld[%v]: %v %v", k1, l1, l2)
			} else if k1 == 12 && len(l1) == 2 && len(l2) == 2 {
				if f1.Ld[12][0].(Corpus) != f2.Ld[12][0].(Corpus) {
					debug.PrintStack()
					t.Errorf("Fighter.Ld[12][0]: %v %v", f1.Ld[12][0], f2.Ld[12][0])
				}
				if f1.Ld[12][1].(string) != f2.Ld[12][1].(string) {
					debug.PrintStack()
					t.Errorf("Fighter.Ld[12][1]: %v %v", f1.Ld[12][1], f2.Ld[12][1])
				}
			} else if k1 == 34 && len(l1) == 2 && len(l2) == 2 {
				compareVector2(t, f1.Ld[34][0].(*Vector2), f2.Ld[34][0].(*Vector2), "Fighter.Ld[34][0]")
				if f1.Ld[34][1].(float64) != f2.Ld[34][1].(float64) {
					debug.PrintStack()
					t.Errorf("Fighter.Ld[34][1]: %v %v", f1.Ld[34][1], f2.Ld[34][1])
				}
			} else {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Ld[%v]: %v %v", k1, l1, l2)
			}
		}
	}
	if len(f1.Fld) != len(f2.Fld) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Fld: %v %v", f1.Fld, f2.Fld)
	} else {
		for k1, l1 := range f1.Fld {
			if l2, ok := f2.Fld[k1]; !ok || len(l1) != len(l2) {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Fld[%v]: %v %v", k1, l1, l2)
			} else {
				for k2, v1 := range l1 {
					v2 := l2[k2]
					if v1 != v2 {
						debug.PrintStack()
						t.Errorf("Fighter_Part2.Fld[%v][%v]: %v %v", k1, k2, v1, v2)
					}
				}
			}
		}
	}
	if len(f1.Dd) != len(f2.Dd) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Dd: %v %v", f1.Dd, f2.Dd)
	} else {
		for k1, l1 := range f1.Dd {
			if l2, ok := f2.Dd[k1]; !ok || len(l1) != len(l2) {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Dd[%v]: %v %v", k1, l1, l2)
			} else if k1 == 12 && len(l1) == 2 && len(l2) == 2 {
				if f1.Dd[12][111].(Corpus) != f2.Dd[12][111].(Corpus) {
					debug.PrintStack()
					t.Errorf("Fighter.Dd[12][111]: %v %v", f1.Dd[12][111], f2.Dd[12][111])
				}
				if f1.Dd[12][222].(string) != f2.Dd[12][222].(string) {
					debug.PrintStack()
					t.Errorf("Fighter.Dd[12][222]: %v %v", f1.Dd[12][222], f2.Dd[12][222])
				}
			} else if k1 == 34 && len(l1) == 2 && len(l2) == 2 {
				compareVector2(t, f1.Dd[34][333].(*Vector2), f2.Dd[34][333].(*Vector2), "Fighter.Dd[34][333]")
				if f1.Dd[34][444].(float64) != f2.Dd[34][444].(float64) {
					debug.PrintStack()
					t.Errorf("Fighter.Dd[34][444]: %v %v", f1.Dd[34][444], f2.Dd[34][444])
				}
			} else {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Dd[%v]: %v %v", k1, l1, l2)
			}
		}
	}
	if len(f1.Fdd) != len(f2.Fdd) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Fdd: %v %v", f1.Fdd, f2.Fdd)
	} else {
		for k1, l1 := range f1.Fdd {
			if l2, ok := f2.Fdd[k1]; !ok || len(l1) != len(l2) {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Fdd[%v]: %v %v", k1, l1, l2)
			} else {
				for k2, v1 := range l1 {
					if v2, ok := l2[k2]; !ok || v1 != v2 {
						debug.PrintStack()
						t.Errorf("Fighter_Part2.Fdd[%v][%v]: %v %v", k1, k2, v1, v2)
					}
				}
			}
		}
	}
	if int32(f1.Nv.(int)) != f2.Nv.(int32) {
		debug.PrintStack()
		t.Errorf("Fighter.Nv: %v %v", f1.Nv, f2.Nv)
	}
	if len(f1.Lv.([]interface{})) != len(f2.Lv.([]interface{})) || len(f1.Lv.([]interface{})) != 2 {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Lv: %v %v", f1.Lv, f2.Lv)
	} else {
		if float32(f1.Lv.([]interface{})[0].(int)) != f2.Lv.([]interface{})[0].(float32) {
			debug.PrintStack()
			t.Errorf("Fighter.Lv[0]: %v %v", f1.Lv.([]interface{})[0], f2.Lv.([]interface{})[0])
		}
		if f1.Lv.([]interface{})[1].(string) != f2.Lv.([]interface{})[1].(string) {
			debug.PrintStack()
			t.Errorf("Fighter.Lv[1]: %v %v", f1.Lv.([]interface{})[1], f2.Lv.([]interface{})[1])
		}
	}
	if len(f1.Flv.([]float32)) != len(f2.Flv.([]float32)) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Flv: %v %v", f1.Flv, f2.Flv)
	} else {
		for k, v1 := range f1.Flv.([]float32) {
			v2 := f2.Flv.([]float32)[k]
			if v1 != v2 {
				debug.PrintStack()
				t.Errorf("Fighter.Flv[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Dv.(map[int32]interface{})) != len(f2.Dv.(map[int32]interface{})) || len(f1.Dv.(map[int32]interface{})) != 2 {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Dv: %v %v", f1.Dv, f2.Dv)
	} else {
		if float32(f1.Dv.(map[int32]interface{})[333].(int)) != f2.Dv.(map[int32]interface{})[333].(float32) {
			debug.PrintStack()
			t.Errorf("Fighter.Dv[333]: %v %v", f1.Dv.(map[int32]interface{})[333], f2.Dv.(map[int32]interface{})[333])
		}
		if f1.Dv.(map[int32]interface{})[444].(string) != f2.Dv.(map[int32]interface{})[444].(string) {
			debug.PrintStack()
			t.Errorf("Fighter.Dv[444]: %v %v", f1.Dv.(map[int32]interface{})[444], f2.Dv.(map[int32]interface{})[444])
		}
	}
	if len(f1.Fdv.(map[int32]float32)) != len(f2.Fdv.(map[int32]float32)) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Fdv: %v %v", f1.Fdv, f2.Fdv)
	} else {
		for k, v1 := range f1.Fdv.(map[int32]float32) {
			if v2, ok := f2.Fdv.(map[int32]float32)[k]; !ok || v1 != v2 {
				debug.PrintStack()
				t.Errorf("Fighter.Fdv[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Poslll) != len(f2.Poslll) {
		debug.PrintStack()
		t.Errorf("Fighter_Part1.Poslll: %v %v", f1.Poslll, f2.Poslll)
	} else {
		for k1, ll1 := range f1.Poslll {
			ll2 := f2.Poslll[k1]
			if len(ll1) != len(ll2) {
				debug.PrintStack()
				t.Errorf("Fighter_Part1.Poslll[%v]: %v %v", k1, ll1, ll2)
			} else {
				for k2, l1 := range ll1 {
					l2 := ll2[k2]
					if len(l1) != len(l2) {
						debug.PrintStack()
						t.Errorf("Fighter_Part1.Poslll[%v][%v]: %v %v", k1, k2, l1, l2)
					} else {
						for k3, v1 := range l1 {
							v2 := l2[k3]
							compareVector2(t, v1, v2, fmt.Sprintf("Fighter_Part1.Poslll[%v][%v][%v]", k1, k2, k3))
						}
					}
				}
			}
		}
	}
	if len(f1.Posdl) != len(f2.Posdl) {
		debug.PrintStack()
		t.Errorf("Fighter_Part2.Posdl: %v %v", f1.Posdl, f2.Posdl)
	} else {
		for k1, d1 := range f1.Posdl {
			d2 := f2.Posdl[k1]
			if len(d1) != len(d2) {
				debug.PrintStack()
				t.Errorf("Fighter_Part2.Posdl[%v]: %v %v", k1, d1, d2)
			} else {
				for k2, v1 := range d1 {
					if v2, ok := d2[k2]; !ok {
						debug.PrintStack()
						t.Errorf("Fighter_Part2.Posdl[%v][%v]: %v %v", k1, k2, v1, v2)
					} else {
						compareVector2(t, v1, v2, fmt.Sprintf("Fighter_Part1.Posdl[%v][%v]", k1, k2))
					}
				}
			}
		}
	}
}

func TestVector2(t *testing.T) {
	vd := &tygo.ProtoBuf{Buffer: make([]byte, v.ByteSize())}
	v.Serialize(vd)
	vd.Reset()
	v3 := &Vector2{}
	if err := v3.Deserialize(vd); err == nil {
		compareVector2(t, v, v3, "")
	} else {
		t.Errorf("TestVector2 Deserialize error: %v", err)
	}
}

func TestFighter(t *testing.T) {
	fd := &tygo.ProtoBuf{Buffer: make([]byte, fighter.ByteSize())}
	fighter.Serialize(fd)
	fd.Reset()
	fighter2 := &Fighter{}
	if err := fighter2.Deserialize(fd); err == nil {
		compareFighter(t, fighter, fighter2)
	} else {
		t.Errorf("TestFighter Deserialize error: %v", err)
	}
}
