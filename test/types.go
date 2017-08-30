// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package test

/*
type Corpus enum {
	UNIVERSAL = 0
	WEB       = 1
	IMAGES    = 2
	LOCAL     = 3
	NEWS      = 4
	PRODUCTS  = 5
	VIDEO     = 6
}

type Vector2 object {
	X float32
	Y fixedpoint<1, -10>
	B bytes
	S string
	M symbol
	E Corpus
	P *GoType
}

type Fighter_Part1 object {
	Pos     *Vector2
	IsAwake bool
	Hp      float32
	Poss    map[int32]*Vector2
	Posi    map[int32]float32
	Posl    []*Vector2
	Posll   [][]*Vector2
	Pyl     []*GoType
	Pyd     map[int32]*GoType
	Pyv1    variant<int32, *GoType>
	Pyv2    variant<int32, *GoType>
}

type Fighter_Part2 object {
	Fighter_Part1
	Fl []float32
	Bl []bytes
	Sl []string
	Bd map[string]bytes
	Sd map[int32]string
	Ml []symbol
	Mbd map[symbol]bytes
	Md map[int32]symbol
	El []Corpus
	Ed map[int32]Corpus
	Ll [][]float32
}

type Fighter object {
	Fighter_Part2
	V0  variant<int32, float32, bytes, *Vector2>
	V1  variant<int32, float32, bytes, *Vector2>
	V2  variant<int32, float32, bytes, *Vector2>
	V3  variant<int32, float32, bytes, *Vector2>
	V4  variant<int32, float32, bytes, *Vector2>
	Vl  []variant<int32, fixedpoint<3, 0>, string, *Vector2>
	Vd  map[int32]variant<Corpus, float64, string, *Vector2>
	Ld  map[int32][]variant<Corpus, float64, string, *Vector2>
	Fld map[int32][]float32
	Dd  map[int32]map[int32]variant<int32, Corpus, float64, string, *Vector2>
	Fdd map[int32]map[int32]float32
	Nv  variant<nil, int32>
	Lv  variant<int32, []variant<float32, string>>
	Flv variant<int32, []float32>
	Dv  variant<int32, map[int32]variant<float32, string>>
	Fdv variant<int32, map[int32]float32>
	Poslll  [][][]*Vector2
	Posdl   []map[string]*Vector2
	RPG(*Fighter, variant<nil, int32>, fixedpoint<3, 0>) *Vector2
	GPR(map[int32]variant<Corpus, float64, string, *Vector2>) (*Fighter, int32)
}

*/
import _ "github.com/ibelie/tygo"

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"runtime/debug"
)

type Errorf func(string, ...interface{})

type GoType struct {
	PP int32
	AP string
}

func (g *GoType) ByteSize() int {
	var w bytes.Buffer
	if err := gob.NewEncoder(&w).Encode(g); err == nil {
		return w.Len()
	} else {
		return 0
	}
}

func (g *GoType) CachedSize() int {
	return g.ByteSize()
}

func (g *GoType) Serialize(w io.Writer) {
	gob.NewEncoder(w).Encode(g)
}

func (g *GoType) Deserialize(r io.Reader) error {
	return gob.NewDecoder(r).Decode(g)
}

var v *Vector2 = &Vector2{
	X: 123,
	Y: 45.6,
	B: []byte("asdf 1234"),
	S: "哈哈哈哈",
	M: "Symbol_AB",
	E: Corpus_LOCAL,
	P: &GoType{PP: 123, AP: "asdf"},
}

var v2 *Vector2 = &Vector2{
	X: 1234,
	Y: 345.6,
	B: []byte("xxx 1234"),
	S: "哈哈 吼吼吼",
	M: "Symbol_ABC",
	E: Corpus_PRODUCTS,
	P: &GoType{PP: 321, AP: "qwer"},
}

var fighter *Fighter = &Fighter{
	Fighter_Part2: Fighter_Part2{
		Fighter_Part1: Fighter_Part1{
			Pos:     v,
			Hp:      12,
			IsAwake: true,
			Poss:    map[int32]*Vector2{321: v, 320: nil, 231: v2},
			Posi:    map[int32]float32{123: 0.456},
			Posl:    []*Vector2{v, &Vector2{}, v2},
			Posll:   [][]*Vector2{[]*Vector2{v, &Vector2{}, v2}, nil, []*Vector2{v2, v}},
			Pyl:     []*GoType{&GoType{PP: 123, AP: "adsf"}, &GoType{}, &GoType{PP: 456, AP: "xxxx"}},
			Pyd:     map[int32]*GoType{321: &GoType{PP: 123, AP: "adsf"}, 654: &GoType{PP: 456, AP: "xxxx"}, 320: nil},
			Pyv1:    123,
			Pyv2:    &GoType{PP: 123, AP: "adsf"},
		},
		Fl:  []float32{0.123, 456, 7.89},
		Sl:  []string{"哈哈", "吼吼", "嘿嘿"},
		Bl:  [][]byte{[]byte("aaa 0.123"), []byte("bbb 456"), []byte("ccc 7.89")},
		Bd:  map[string][]byte{"哈哈": []byte("aaa 0.123"), "asdf": []byte("bbb 456")},
		Sd:  map[int32]string{321: "哈哈 3", 231: "吼吼 2"},
		Ml:  []string{"Symbol_ABCD", "", "Sym"},
		Mbd: map[string][]byte{"Symbol_aaa": []byte("aaa 0.123"), "Symbol_bbb": []byte("bbb 456")},
		Md:  map[int32]string{321: "S", 231: "Symbol_231"},
		El:  []Corpus{Corpus_LOCAL, Corpus_NEWS, Corpus_VIDEO},
		Ed:  map[int32]Corpus{789: Corpus_WEB, 567: Corpus_IMAGES},
		Ll:  [][]float32{[]float32{12.3, 1.23}, []float32{1.234, 12.34, 123.4}},
	},
	V1:     98765,
	V2:     []byte("adsf"),
	V3:     v,
	V4:     345.123,
	Vl:     []interface{}{123, "adsf", nil, v, 345.123, nil},
	Vd:     map[int32]interface{}{0: nil, 12: Corpus_IMAGES, 23: "adsf", 34: v2, 45: 345.123},
	Ld:     map[int32][]interface{}{12: []interface{}{Corpus_IMAGES, "adsf"}, 34: []interface{}{v2, 345.123}},
	Fld:    map[int32][]float32{123: []float32{222.111, 345.123}},
	Dd:     map[int32]map[int32]interface{}{12: map[int32]interface{}{111: Corpus_IMAGES, 222: "adsf"}, 34: map[int32]interface{}{333: v2, 444: 345.123}},
	Fdd:    map[int32]map[int32]float32{123: map[int32]float32{12: 222.111, 23: 345.123}},
	Nv:     123456,
	Lv:     []interface{}{123, "adsf"},
	Flv:    []float32{222.111, 345.123},
	Dv:     map[int32]interface{}{333: 123, 444: "adsf"},
	Fdv:    map[int32]float32{333: 222.111, 444: 345.123},
	Poslll: [][][]*Vector2{[][]*Vector2{[]*Vector2{v, &Vector2{}, v2}, nil}, nil, [][]*Vector2{[]*Vector2{v, &Vector2{}, v2}, []*Vector2{v2, v}}},
	Posdl:  []map[string]*Vector2{map[string]*Vector2{"231": v, "320": nil, "321": v2}, nil, map[string]*Vector2{"321": v, "320": nil, "231": v2}},
}

func CompareGoType(p Errorf, g1 *GoType, g2 *GoType, prefix string) {
	if g1 == g2 {
		return
	} else if g1 == nil || g2 == nil {
		debug.PrintStack()
		p("%s GoType %v %v", prefix, g1, g2)
	} else if g1.PP != g2.PP {
		debug.PrintStack()
		p("%s GoType.PP: %v %v", prefix, g1.PP, g2.PP)
	} else if g1.AP != g2.AP {
		debug.PrintStack()
		p("%s GoType.AP: %v %v", prefix, g1.AP, g2.AP)
	}
}

func CompareVector2(p Errorf, v1 *Vector2, v2 *Vector2, prefix string) {
	if v1 == v2 {
		return
	} else if v1 == nil || v2 == nil {
		debug.PrintStack()
		p("%s Vector2 %v %v", prefix, v1, v2)
	} else if v1.X != v2.X {
		debug.PrintStack()
		p("%s Vector2.X: %v %v", prefix, v1.X, v2.X)
	} else if v1.Y != v2.Y {
		debug.PrintStack()
		p("%s Vector2.Y: %v %v", prefix, v1.Y, v2.Y)
	} else if v1.S != v2.S {
		debug.PrintStack()
		p("%s Vector2.S: %v %v", prefix, v1.S, v2.S)
	} else if v1.M != v2.M {
		debug.PrintStack()
		p("%s Vector2.M: %v %v", prefix, v1.M, v2.M)
	} else if bytes.Compare(v1.B, v2.B) != 0 {
		debug.PrintStack()
		p("%s Vector2.B: %v %v", prefix, v1.B, v2.B)
	} else if v1.E != v2.E {
		debug.PrintStack()
		p("%s Vector2.E: %v %v", prefix, v1.E, v2.E)
	} else {
		CompareGoType(p, v1.P, v2.P, prefix+".P")
	}
}

func CompareFighter_Part1(p Errorf, f1 *Fighter_Part1, f2 *Fighter_Part1) {
	if f1 == f2 {
		return
	} else if f1 == nil || f2 == nil {
		debug.PrintStack()
		p("Fighter_Part1 %v %v", f1, f2)
		return
	}
	CompareVector2(p, f1.Pos, f2.Pos, "Fighter_Part1.Pos")
	if f1.IsAwake != f2.IsAwake {
		debug.PrintStack()
		p("Fighter_Part1.IsAwake: %v %v", f1.IsAwake, f2.IsAwake)
	}
	if f1.Hp != f2.Hp {
		debug.PrintStack()
		p("Fighter_Part1.Hp: %v %v", f1.Hp, f2.Hp)
	}
	if len(f1.Poss) != len(f2.Poss) {
		debug.PrintStack()
		p("Fighter_Part1.Poss: %v %v", f1.Poss, f2.Poss)
	} else {
		for k, v1 := range f1.Poss {
			if v2, ok := f2.Poss[k]; !ok {
				debug.PrintStack()
				p("Fighter_Part1.Poss[%v]: %v %v", k, v1, v2)
			} else {
				CompareVector2(p, v1, v2, fmt.Sprintf("Fighter_Part1.Poss[%v]", k))
			}
		}
	}
	if len(f1.Posi) != len(f2.Posi) {
		debug.PrintStack()
		p("Fighter_Part1.Posi: %v %v", f1.Posi, f2.Posi)
	} else {
		for k, v1 := range f1.Posi {
			if v2, ok := f2.Posi[k]; !ok || v1 != v2 {
				debug.PrintStack()
				p("Fighter_Part1.Posi[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Posl) != len(f2.Posl) {
		debug.PrintStack()
		p("Fighter_Part1.Posl: %v %v", f1.Posl, f2.Posl)
	} else {
		for k, v1 := range f1.Posl {
			v2 := f2.Posl[k]
			CompareVector2(p, v1, v2, fmt.Sprintf("Fighter_Part1.Posl[%v]", k))
		}
	}
	if len(f1.Posll) != len(f2.Posll) {
		debug.PrintStack()
		p("Fighter_Part1.Posll: %v %v", f1.Posll, f2.Posll)
	} else {
		for k1, l1 := range f1.Posll {
			l2 := f2.Posll[k1]
			if len(l1) != len(l2) {
				debug.PrintStack()
				p("Fighter_Part1.Posll[%v]: %v %v", k1, l1, l2)
			} else {
				for k2, v1 := range l1 {
					v2 := l2[k2]
					CompareVector2(p, v1, v2, fmt.Sprintf("Fighter_Part1.Posll[%v][%v]", k1, k2))
				}
			}
		}
	}
	if len(f1.Pyl) != len(f2.Pyl) {
		debug.PrintStack()
		p("Fighter_Part1.Pyl: %v %v", f1.Pyl, f2.Pyl)
	} else {
		for k, v1 := range f1.Pyl {
			v2 := f2.Pyl[k]
			CompareGoType(p, v1, v2, fmt.Sprintf("Fighter_Part1.Pyl[%v]", k))
		}
	}
	if len(f1.Pyd) != len(f2.Pyd) {
		debug.PrintStack()
		p("Fighter_Part1.Pyd: %v %v", f1.Pyd, f2.Pyd)
	} else {
		for k, v1 := range f1.Pyd {
			if v2, ok := f2.Pyd[k]; !ok {
				debug.PrintStack()
				p("Fighter_Part1.Pyd[%v]: %v %v", k, v1, v2)
			} else {
				CompareGoType(p, v1, v2, fmt.Sprintf("Fighter_Part1.Pyd[%v]", k))
			}
		}
	}
	if int32(f1.Pyv1.(int)) != f2.Pyv1.(int32) {
		debug.PrintStack()
		p("Fighter_Part1.Pyv1: %v %v", f1.Pyv1, f2.Pyv1)
	}
	CompareGoType(p, f1.Pyv2.(*GoType), f2.Pyv2.(*GoType), "Fighter_Part1.Pyv1")
}

func CompareFighter_Part2(p Errorf, f1 *Fighter_Part2, f2 *Fighter_Part2) {
	if f1 == f2 {
		return
	} else if f1 == nil || f2 == nil {
		debug.PrintStack()
		p("Fighter_Part2 %v %v", f1, f2)
		return
	}
	CompareFighter_Part1(p, &f1.Fighter_Part1, &f2.Fighter_Part1)
	if len(f1.Fl) != len(f2.Fl) {
		debug.PrintStack()
		p("Fighter_Part2.Fl: %v %v", f1.Fl, f2.Fl)
	} else {
		for k, v1 := range f1.Fl {
			v2 := f2.Fl[k]
			if v1 != v2 {
				debug.PrintStack()
				p("Fighter_Part2.Fl[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Bl) != len(f2.Bl) {
		debug.PrintStack()
		p("Fighter_Part2.Bl: %v %v", f1.Bl, f2.Bl)
	} else {
		for k, v1 := range f1.Bl {
			v2 := f2.Bl[k]
			if bytes.Compare(v1, v2) != 0 {
				debug.PrintStack()
				p("Fighter_Part2.Bl[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Sl) != len(f2.Sl) {
		debug.PrintStack()
		p("Fighter_Part2.Sl: %v %v", f1.Sl, f2.Sl)
	} else {
		for k, v1 := range f1.Sl {
			v2 := f2.Sl[k]
			if v1 != v2 {
				debug.PrintStack()
				p("Fighter_Part2.Sl[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Bd) != len(f2.Bd) {
		debug.PrintStack()
		p("Fighter_Part2.Bd: %v %v", f1.Bd, f2.Bd)
	} else {
		for k, v1 := range f1.Bd {
			if v2, ok := f2.Bd[k]; !ok || bytes.Compare(v1, v2) != 0 {
				debug.PrintStack()
				p("Fighter_Part2.Bd[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Sd) != len(f2.Sd) {
		debug.PrintStack()
		p("Fighter_Part2.Sd: %v %v", f1.Sd, f2.Sd)
	} else {
		for k, v1 := range f1.Sd {
			if v2, ok := f2.Sd[k]; !ok || v1 != v2 {
				debug.PrintStack()
				p("Fighter_Part2.Sd[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Ml) != len(f2.Ml) {
		debug.PrintStack()
		p("Fighter_Part2.Ml: %v %v", f1.Ml, f2.Ml)
	} else {
		for k, v1 := range f1.Ml {
			v2 := f2.Ml[k]
			if v1 != v2 {
				debug.PrintStack()
				p("Fighter_Part2.Ml[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Mbd) != len(f2.Mbd) {
		debug.PrintStack()
		p("Fighter_Part2.Mbd: %v %v", f1.Mbd, f2.Mbd)
	} else {
		for k, v1 := range f1.Mbd {
			if v2, ok := f2.Mbd[k]; !ok || bytes.Compare(v1, v2) != 0 {
				debug.PrintStack()
				p("Fighter_Part2.Mbd[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Md) != len(f2.Md) {
		debug.PrintStack()
		p("Fighter_Part2.Md: %v %v", f1.Md, f2.Md)
	} else {
		for k, v1 := range f1.Md {
			if v2, ok := f2.Md[k]; !ok || v1 != v2 {
				debug.PrintStack()
				p("Fighter_Part2.Md[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.El) != len(f2.El) {
		debug.PrintStack()
		p("Fighter_Part2.El: %v %v", f1.El, f2.El)
	} else {
		for k, v1 := range f1.El {
			v2 := f2.El[k]
			if v1 != v2 {
				debug.PrintStack()
				p("Fighter_Part2.El[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Ed) != len(f2.Ed) {
		debug.PrintStack()
		p("Fighter_Part2.Ed: %v %v", f1.Ed, f2.Ed)
	} else {
		for k, v1 := range f1.Ed {
			if v2, ok := f2.Ed[k]; !ok || v1 != v2 {
				debug.PrintStack()
				p("Fighter_Part2.Ed[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Ll) != len(f2.Ll) {
		debug.PrintStack()
		p("Fighter_Part2.Ll: %v %v", f1.Ll, f2.Ll)
	} else {
		for k1, l1 := range f1.Ll {
			l2 := f2.Ll[k1]
			if len(l1) != len(l2) {
				debug.PrintStack()
				p("Fighter_Part2.Ll[%v]: %v %v", k1, l1, l2)
			} else {
				for k2, v1 := range l1 {
					v2 := l2[k2]
					if v1 != v2 {
						debug.PrintStack()
						p("Fighter_Part2.Ll[%v][%v]: %v %v", k1, k2, v1, v2)
					}
				}
			}
		}
	}
}

func CompareFighter(p Errorf, f1 *Fighter, f2 *Fighter) {
	if f1 == f2 {
		return
	} else if f1 == nil || f2 == nil {
		debug.PrintStack()
		p("Fighter %v %v", f1, f2)
		return
	}
	CompareFighter_Part2(p, &f1.Fighter_Part2, &f2.Fighter_Part2)
	if f1.V0 != f2.V0 {
		debug.PrintStack()
		p("Fighter.V0: %v %v", f1.V0, f2.V0)
	}
	if int32(f1.V1.(int)) != f2.V1.(int32) {
		debug.PrintStack()
		p("Fighter.V1: %v %v", f1.V1, f2.V1)
	}
	if bytes.Compare(f1.V2.([]byte), f2.V2.([]byte)) != 0 {
		debug.PrintStack()
		p("Fighter.V2: %v %v", f1.V2, f2.V2)
	}
	CompareVector2(p, f1.V3.(*Vector2), f2.V3.(*Vector2), "Fighter.V3")
	if float32(f1.V4.(float64)) != f2.V4.(float32) {
		debug.PrintStack()
		p("Fighter.V4: %v %v", f1.V4, f2.V4)
	}
	if len(f1.Vl) != len(f2.Vl) {
		debug.PrintStack()
		p("Fighter_Part2.Vl: %v %v", f1.Vl, f2.Vl)
	} else {
		if int32(f1.Vl[0].(int)) != f2.Vl[0].(int32) {
			debug.PrintStack()
			p("Fighter.Vl[0]: %v %v", f1.Vl[0], f2.Vl[0])
		}
		if f1.Vl[1].(string) != f2.Vl[1].(string) {
			debug.PrintStack()
			p("Fighter.Vl[1]: %v %v", f1.Vl[1], f2.Vl[1])
		}
		if f1.Vl[2] != nil || f2.Vl[2] != nil {
			debug.PrintStack()
			p("Fighter.Vl[2]: %v %v", f1.Vl[2], f2.Vl[2])
		}
		CompareVector2(p, f1.Vl[3].(*Vector2), f2.Vl[3].(*Vector2), "Fighter.Vl[3]")
		if f1.Vl[4].(float64) != f2.Vl[4].(float64) {
			debug.PrintStack()
			p("Fighter.Vl[4]: %v %v", f1.Vl[4], f2.Vl[4])
		}
		if f1.Vl[5] != nil || f2.Vl[5] != nil {
			debug.PrintStack()
			p("Fighter.Vl[5]: %v %v", f1.Vl[5], f2.Vl[5])
		}
	}
	if len(f1.Vd) != len(f2.Vd) {
		debug.PrintStack()
		p("Fighter_Part2.Vd: %v %v", f1.Vd, f2.Vd)
	} else {
		if f1.Vd[0] != nil || f2.Vd[0] != nil {
			debug.PrintStack()
			p("Fighter.Vd[0]: %v %v", f1.Vd[0], f2.Vd[0])
		}
		if f1.Vd[12].(Corpus) != f2.Vd[12].(Corpus) {
			debug.PrintStack()
			p("Fighter.Vd[12]: %v %v", f1.Vd[12], f2.Vd[12])
		}
		if f1.Vd[23].(string) != f2.Vd[23].(string) {
			debug.PrintStack()
			p("Fighter.Vd[23]: %v %v", f1.Vd[23], f2.Vd[23])
		}
		CompareVector2(p, f1.Vd[34].(*Vector2), f2.Vd[34].(*Vector2), "Fighter.Vd[34]")
		if f1.Vd[45].(float64) != f2.Vd[45].(float64) {
			debug.PrintStack()
			p("Fighter.Vd[45]: %v %v", f1.Vd[45], f2.Vd[45])
		}
	}
	if len(f1.Ld) != len(f2.Ld) {
		debug.PrintStack()
		p("Fighter_Part2.Ld: %v %v", f1.Ld, f2.Ld)
	} else {
		for k1, l1 := range f1.Ld {
			if l2, ok := f2.Ld[k1]; !ok || len(l1) != len(l2) {
				debug.PrintStack()
				p("Fighter_Part2.Ld[%v]: %v %v", k1, l1, l2)
			} else if k1 == 12 && len(l1) == 2 && len(l2) == 2 {
				if f1.Ld[12][0].(Corpus) != f2.Ld[12][0].(Corpus) {
					debug.PrintStack()
					p("Fighter.Ld[12][0]: %v %v", f1.Ld[12][0], f2.Ld[12][0])
				}
				if f1.Ld[12][1].(string) != f2.Ld[12][1].(string) {
					debug.PrintStack()
					p("Fighter.Ld[12][1]: %v %v", f1.Ld[12][1], f2.Ld[12][1])
				}
			} else if k1 == 34 && len(l1) == 2 && len(l2) == 2 {
				CompareVector2(p, f1.Ld[34][0].(*Vector2), f2.Ld[34][0].(*Vector2), "Fighter.Ld[34][0]")
				if f1.Ld[34][1].(float64) != f2.Ld[34][1].(float64) {
					debug.PrintStack()
					p("Fighter.Ld[34][1]: %v %v", f1.Ld[34][1], f2.Ld[34][1])
				}
			} else {
				debug.PrintStack()
				p("Fighter_Part2.Ld[%v]: %v %v", k1, l1, l2)
			}
		}
	}
	if len(f1.Fld) != len(f2.Fld) {
		debug.PrintStack()
		p("Fighter_Part2.Fld: %v %v", f1.Fld, f2.Fld)
	} else {
		for k1, l1 := range f1.Fld {
			if l2, ok := f2.Fld[k1]; !ok || len(l1) != len(l2) {
				debug.PrintStack()
				p("Fighter_Part2.Fld[%v]: %v %v", k1, l1, l2)
			} else {
				for k2, v1 := range l1 {
					v2 := l2[k2]
					if v1 != v2 {
						debug.PrintStack()
						p("Fighter_Part2.Fld[%v][%v]: %v %v", k1, k2, v1, v2)
					}
				}
			}
		}
	}
	if len(f1.Dd) != len(f2.Dd) {
		debug.PrintStack()
		p("Fighter_Part2.Dd: %v %v", f1.Dd, f2.Dd)
	} else {
		for k1, d1 := range f1.Dd {
			if d2, ok := f2.Dd[k1]; !ok || len(d1) != len(d2) {
				debug.PrintStack()
				p("Fighter_Part2.Dd[%v]: %v %v", k1, d1, d2)
			} else if k1 == 12 && len(d1) == 2 && len(d2) == 2 {
				if f1.Dd[12][111].(Corpus) != f2.Dd[12][111].(Corpus) {
					debug.PrintStack()
					p("Fighter.Dd[12][111]: %v %v", f1.Dd[12][111], f2.Dd[12][111])
				}
				if f1.Dd[12][222].(string) != f2.Dd[12][222].(string) {
					debug.PrintStack()
					p("Fighter.Dd[12][222]: %v %v", f1.Dd[12][222], f2.Dd[12][222])
				}
			} else if k1 == 34 && len(d1) == 2 && len(d2) == 2 {
				CompareVector2(p, f1.Dd[34][333].(*Vector2), f2.Dd[34][333].(*Vector2), "Fighter.Dd[34][333]")
				if f1.Dd[34][444].(float64) != f2.Dd[34][444].(float64) {
					debug.PrintStack()
					p("Fighter.Dd[34][444]: %v %v", f1.Dd[34][444], f2.Dd[34][444])
				}
			} else {
				debug.PrintStack()
				p("Fighter_Part2.Dd[%v]: %v %v", k1, d1, d2)
			}
		}
	}
	if len(f1.Fdd) != len(f2.Fdd) {
		debug.PrintStack()
		p("Fighter_Part2.Fdd: %v %v", f1.Fdd, f2.Fdd)
	} else {
		for k1, d1 := range f1.Fdd {
			if d2, ok := f2.Fdd[k1]; !ok || len(d1) != len(d2) {
				debug.PrintStack()
				p("Fighter_Part2.Fdd[%v]: %v %v", k1, d1, d2)
			} else {
				for k2, v1 := range d1 {
					if v2, ok := d2[k2]; !ok || v1 != v2 {
						debug.PrintStack()
						p("Fighter_Part2.Fdd[%v][%v]: %v %v", k1, k2, v1, v2)
					}
				}
			}
		}
	}
	if int32(f1.Nv.(int)) != f2.Nv.(int32) {
		debug.PrintStack()
		p("Fighter.Nv: %v %v", f1.Nv, f2.Nv)
	}
	if len(f1.Lv.([]interface{})) != len(f2.Lv.([]interface{})) || len(f1.Lv.([]interface{})) != 2 {
		debug.PrintStack()
		p("Fighter_Part2.Lv: %v %v", f1.Lv, f2.Lv)
	} else {
		if float32(f1.Lv.([]interface{})[0].(int)) != f2.Lv.([]interface{})[0].(float32) {
			debug.PrintStack()
			p("Fighter.Lv[0]: %v %v", f1.Lv.([]interface{})[0], f2.Lv.([]interface{})[0])
		}
		if f1.Lv.([]interface{})[1].(string) != f2.Lv.([]interface{})[1].(string) {
			debug.PrintStack()
			p("Fighter.Lv[1]: %v %v", f1.Lv.([]interface{})[1], f2.Lv.([]interface{})[1])
		}
	}
	if len(f1.Flv.([]float32)) != len(f2.Flv.([]float32)) {
		debug.PrintStack()
		p("Fighter_Part2.Flv: %v %v", f1.Flv, f2.Flv)
	} else {
		for k, v1 := range f1.Flv.([]float32) {
			v2 := f2.Flv.([]float32)[k]
			if v1 != v2 {
				debug.PrintStack()
				p("Fighter.Flv[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Dv.(map[int32]interface{})) != len(f2.Dv.(map[int32]interface{})) || len(f1.Dv.(map[int32]interface{})) != 2 {
		debug.PrintStack()
		p("Fighter_Part2.Dv: %v %v", f1.Dv, f2.Dv)
	} else {
		if float32(f1.Dv.(map[int32]interface{})[333].(int)) != f2.Dv.(map[int32]interface{})[333].(float32) {
			debug.PrintStack()
			p("Fighter.Dv[333]: %v %v", f1.Dv.(map[int32]interface{})[333], f2.Dv.(map[int32]interface{})[333])
		}
		if f1.Dv.(map[int32]interface{})[444].(string) != f2.Dv.(map[int32]interface{})[444].(string) {
			debug.PrintStack()
			p("Fighter.Dv[444]: %v %v", f1.Dv.(map[int32]interface{})[444], f2.Dv.(map[int32]interface{})[444])
		}
	}
	if len(f1.Fdv.(map[int32]float32)) != len(f2.Fdv.(map[int32]float32)) {
		debug.PrintStack()
		p("Fighter_Part2.Fdv: %v %v", f1.Fdv, f2.Fdv)
	} else {
		for k, v1 := range f1.Fdv.(map[int32]float32) {
			if v2, ok := f2.Fdv.(map[int32]float32)[k]; !ok || v1 != v2 {
				debug.PrintStack()
				p("Fighter.Fdv[%v]: %v %v", k, v1, v2)
			}
		}
	}
	if len(f1.Poslll) != len(f2.Poslll) {
		debug.PrintStack()
		p("Fighter_Part1.Poslll: %v %v", f1.Poslll, f2.Poslll)
	} else {
		for k1, ll1 := range f1.Poslll {
			ll2 := f2.Poslll[k1]
			if len(ll1) != len(ll2) {
				debug.PrintStack()
				p("Fighter_Part1.Poslll[%v]: %v %v", k1, ll1, ll2)
			} else {
				for k2, l1 := range ll1 {
					l2 := ll2[k2]
					if len(l1) != len(l2) {
						debug.PrintStack()
						p("Fighter_Part1.Poslll[%v][%v]: %v %v", k1, k2, l1, l2)
					} else {
						for k3, v1 := range l1 {
							v2 := l2[k3]
							CompareVector2(p, v1, v2, fmt.Sprintf("Fighter_Part1.Poslll[%v][%v][%v]", k1, k2, k3))
						}
					}
				}
			}
		}
	}
	if len(f1.Posdl) != len(f2.Posdl) {
		debug.PrintStack()
		p("Fighter_Part2.Posdl: %v %v", f1.Posdl, f2.Posdl)
	} else {
		for k1, d1 := range f1.Posdl {
			d2 := f2.Posdl[k1]
			if len(d1) != len(d2) {
				debug.PrintStack()
				p("Fighter_Part2.Posdl[%v]: %v %v", k1, d1, d2)
			} else {
				for k2, v1 := range d1 {
					if v2, ok := d2[k2]; !ok {
						debug.PrintStack()
						p("Fighter_Part2.Posdl[%v][%v]: %v %v", k1, k2, v1, v2)
					} else {
						CompareVector2(p, v1, v2, fmt.Sprintf("Fighter_Part1.Posdl[%v][%v]", k1, k2))
					}
				}
			}
		}
	}
}
