// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"

	"encoding/gob"
	"github.com/ibelie/tygo"
)

type Corpus uint8

const (
	Corpus_UNIVERSAL Corpus = iota
	Corpus_WEB
	Corpus_IMAGES
	Corpus_LOCAL
	Corpus_NEWS
	Corpus_PRODUCTS
	Corpus_VIDEO
)

type GoType struct {
	PP int32
	AP string
}

func (g *GoType) ByteSize() (int, error) {
	var w bytes.Buffer
	err := gob.NewEncoder(w).Encode(g)
	return w.Len(), err
}

func (g *GoType) Serialize(w io.Writer) error {
	return gob.NewEncoder(w).Encode(g)
}

func (g *GoType) Deserialize(r io.Reader) error {
	return gob.NewDecoder().Decode(g)
}

type Vector2 struct {
	X float32         // @Property(坐标X)
	Y tygo.FixedPoint // @Property(坐标Y) @FixedPoint(1, -10)
	B []byte          // @Property
	S string          // @Property
	E Corpus          // @Property(Corpus)
	P *GoType         // @Property(GoType)
}

type Fighter_Part1 struct {
	Pos     *Vector2           // @Property(坐标)
	IsAwake bool               // @Property
	Hp      float32            // @Property(血量)
	Poss    map[int32]*Vector2 // @Property
	Posi    map[int32]float32  // @Property
	Posl    []*Vector2         // @Property
	Posll   [][]*Vector2       // @Property
	Pyl     []*GoType          // @Property
	Pyd     map[int32]*GoType  // @Property
	Pyv1    interface{}        // @Property @Variant(int32, GoType)
	Pyv2    interface{}        // @Property @Variant(int32, GoType)
}

type Fighter_Part2 struct {
	Fighter_Part1
	Fl []float32         // @Property
	Bl [][]byte          // @Property
	Sl []string          // @Property
	Bd map[string][]byte // @Property
	Sd map[int32]string  // @Property
	El []Corpus          // @Property
	Ed map[int32]Corpus  // @Property
	Ll [][]float32       // @Property
}

type Fighter struct {
	Fighter_Part2
	V0  interface{}                     // @Property @Variant(int32, float32, []byte, Vector2)
	V1  interface{}                     // @Property @Variant(int32, float32, []byte, Vector2)
	V2  interface{}                     // @Property @Variant(int32, float32, []byte, Vector2)
	V3  interface{}                     // @Property @Variant(int32, float32, []byte, Vector2)
	V4  interface{}                     // @Property @Variant(int32, float32, []byte, Vector2)
	Vl  []interface{}                   // @Property @Variant(int32, FixedPoint(3), string, Vector2)
	Vd  map[int32]interface{}           // @Property @Variant(Corpus, float64, string, Vector2)
	Ld  map[int32][]interface{}         // @Property @Variant(Corpus, float64, string, Vector2)
	Fld map[int32][]float32             // @Property
	Dd  map[int32]map[int32]interface{} // @Property @Variant(int32, Corpus, float64, string, Vector2)
	Fdd map[int32]map[int32]float32     // @Property
	Nv  interface{}                     // @Property @Variant(nil, int32)
	Lv  interface{}                     // @Property @Variant(int32, List(float32, string))
	Flv interface{}                     // @Property @Variant(int32, List(float32))
	Dv  interface{}                     // @Property @Variant(int32, Dict(int32, float32, string))
	Fdv interface{}                     // @Property @Variant(int32, Dict(int32, float32))
}

// @Procedure @Variant(nil, int32) @Variant(Corpus, float64, string, Vector2) FixedPoint(3)
func (f *Fighter) RPG(fighter *Fighter, nv interface{}, vd map[int32]interface{}, h tygo.FixedPoint) *Vector2 {
	return nil
}

var v *Vector2 = &Vector2{
	X: 123,
	Y: 45.6,
	B: []byte("asdf 1234"),
	S: "哈哈哈哈",
	E: Corpus_LOCAL,
	P: &GoType{pp: 123, ap: "asdf"},
}

var v2 *Vector2 = &Vector2{
	X: 1234,
	Y: 345.6,
	B: []byte("xxx 1234"),
	S: "哈哈 吼吼吼",
	E: Corpus_PRODUCTS,
	P: &GoType{pp: 321, ap: "qwer"},
}

var fighter Fighter = Fighter{
	Pos:     v,
	Hp:      12,
	IsAwake: true,
	Fl:      []float32{0.123, 456, 7.89},
	Sl:      []string{"哈哈", "吼吼", "嘿嘿"},
	Bl:      [][]byte{[]byte("aaa 0.123"), []byte("bbb 456"), []byte("ccc 7.89")},
	V1:      98765,
	V2:      []byte("adsf"),
	V3:      v,
	V4:      345.123,
	Poss:    map[int32]*Vector2{321: v, 320: nil, 231: v2},
	Posi:    map[int32]float32{123: 0.456},
	Posl:    []*Vector2{v, nil, v2},
	Bd:      map[string][]byte{"哈哈": []byte("aaa 0.123"), "asdf": []byte("bbb 456")},
	Sd:      map[int32]string{321: "哈哈 3", 231: "吼吼 2"},
	El:      []Corpus{Corpus_LOCAL, Corpus_NEWS, Corpus_VIDEO},
	Ed:      map[int32]Corpus{789: Corpus_WEB, 567: Corpus_IMAGES},
	Ll:      [][]float32{[]float32{12.3, 1.23}, []float32{1.234, 12.34, 123.4}},
	Vl:      []interface{}{123, "adsf", nil, v, 345.123, nil},
	Vd:      map[int32]interface{}{0: nil, 12: Corpus_IMAGES, 23: "adsf", 34: v2, 45: 345.123},
	Ld:      map[int32][]interface{}{12: []interface{}{Corpus_IMAGES, "adsf"}, 34: []interface{}{v2, 345.123}},
	Fld:     map[int32][]float32{123: []float32{222.111, 345.123}},
	Pyl:     []*GoType{&GoType{pp: 123, ap: "adsf"}, nil, &GoType{pp: 456, ap: "xxxx"}},
	Pyd:     map[int32]*GoType{321: &GoType{pp: 123, ap: "adsf"}, 654: &GoType{pp: 456, ap: "xxxx"}, 320: nil},
	Pyv1:    123,
	Pyv2:    &GoType{pp: 123, ap: "adsf"},
	Dd:      map[int32]map[int32]interface{}{12: map[int32]interface{}{111: Corpus_MAGES, 222: "adsf"}, 34: map[int32]interface{}{333: v2, 444: 345.123}},
	Fdd:     map[int32]map[int32]float32{123: map[int32]float32{12: 222.111, 23: 345.123}},
	Nv:      123456,
	Lv:      []interface{}{123, "adsf"},
	Flv:     []float32{222.111, 345.123},
	Dv:      map[int32]interface{}{333: 123, 444: "adsf"},
	Fdv:     map[int32]float32{333: 222.111, 444: 345.123},
}
