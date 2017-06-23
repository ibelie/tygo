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
	"io"
)

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
	E: Corpus_LOCAL,
	P: &GoType{PP: 123, AP: "asdf"},
}

var v2 *Vector2 = &Vector2{
	X: 1234,
	Y: 345.6,
	B: []byte("xxx 1234"),
	S: "哈哈 吼吼吼",
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
		Fl: []float32{0.123, 456, 7.89},
		Sl: []string{"哈哈", "吼吼", "嘿嘿"},
		Bl: [][]byte{[]byte("aaa 0.123"), []byte("bbb 456"), []byte("ccc 7.89")},
		Bd: map[string][]byte{"哈哈": []byte("aaa 0.123"), "asdf": []byte("bbb 456")},
		Sd: map[int32]string{321: "哈哈 3", 231: "吼吼 2"},
		El: []Corpus{Corpus_LOCAL, Corpus_NEWS, Corpus_VIDEO},
		Ed: map[int32]Corpus{789: Corpus_WEB, 567: Corpus_IMAGES},
		Ll: [][]float32{[]float32{12.3, 1.23}, []float32{1.234, 12.34, 123.4}},
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
