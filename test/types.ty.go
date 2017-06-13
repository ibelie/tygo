// Generated by tygo.  DO NOT EDIT!

package main

import "fmt"
import "github.com/ibelie/tygo"
import "io"

type Corpus int

const (
	Corpus_UNIVERSAL Corpus = 0
	Corpus_WEB       Corpus = 1
	Corpus_IMAGES    Corpus = 2
	Corpus_LOCAL     Corpus = 3
	Corpus_NEWS      Corpus = 4
	Corpus_PRODUCTS  Corpus = 5
	Corpus_VIDEO     Corpus = 6
)

func (i Corpus) String() string {
	switch i {
	case Corpus_UNIVERSAL:
		return "UNIVERSAL"
	case Corpus_WEB:
		return "WEB"
	case Corpus_IMAGES:
		return "IMAGES"
	case Corpus_LOCAL:
		return "LOCAL"
	case Corpus_NEWS:
		return "NEWS"
	case Corpus_PRODUCTS:
		return "PRODUCTS"
	case Corpus_VIDEO:
		return "VIDEO"
	default:
		panic(fmt.Sprintf("[Tygo][Corpus] Unexpect enum value: %d", i))
		return "UNKNOWN"
	}
}

type Vector2 struct {
	tygo.Tygo
	X float32 // float32
	Y float64 // fixedpoint<1, -10>
	B []byte  // bytes
	S string  // string
	E Corpus  // Corpus
	P *GoType // *GoType
}

func (s *Vector2) ByteSize() (int, error) {
	return 0, nil
}

func (s *Vector2) Serialize(w io.Writer) error {
	return nil
}

func (s *Vector2) Deserialize(r io.Reader) error {
	return nil
}

type Fighter_Part1 struct {
	tygo.Tygo
	Pos     *Vector2           // *Vector2
	IsAwake bool               // bool
	Hp      float32            // float32
	Poss    map[int32]*Vector2 // map[int32]*Vector2
	Posi    map[int32]float32  // map[int32]float32
	Posl    []*Vector2         // []*Vector2
	Posll   [][]*Vector2       // [][]*Vector2
	Pyl     []*GoType          // []*GoType
	Pyd     map[int32]*GoType  // map[int32]*GoType
	Pyv1    interface{}        // variant<int32, *GoType>
	Pyv2    interface{}        // variant<int32, *GoType>
}

func (s *Fighter_Part1) ByteSize() (int, error) {
	return 0, nil
}

func (s *Fighter_Part1) Serialize(w io.Writer) error {
	return nil
}

func (s *Fighter_Part1) Deserialize(r io.Reader) error {
	return nil
}

type Fighter_Part2 struct {
	Fighter_Part1
	Fl []float32         // []float32
	Bl [][]byte          // []bytes
	Sl []string          // []string
	Bd map[string][]byte // map[string]bytes
	Sd map[int32]string  // map[int32]string
	El []Corpus          // []Corpus
	Ed map[int32]Corpus  // map[int32]Corpus
	Ll [][]float32       // [][]float32
}

func (s *Fighter_Part2) ByteSize() (int, error) {
	return 0, nil
}

func (s *Fighter_Part2) Serialize(w io.Writer) error {
	return nil
}

func (s *Fighter_Part2) Deserialize(r io.Reader) error {
	return nil
}

type Fighter struct {
	Fighter_Part2
	V0  interface{}                     // variant<int32, float32, bytes, *Vector2>
	V1  interface{}                     // variant<int32, float32, bytes, *Vector2>
	V2  interface{}                     // variant<int32, float32, bytes, *Vector2>
	V3  interface{}                     // variant<int32, float32, bytes, *Vector2>
	V4  interface{}                     // variant<int32, float32, bytes, *Vector2>
	Vl  []interface{}                   // []variant<int32, fixedpoint<3, 0>, string, *Vector2>
	Vd  map[int32]interface{}           // map[int32]variant<Corpus, float64, string, *Vector2>
	Ld  map[int32][]interface{}         // map[int32][]variant<Corpus, float64, string, *Vector2>
	Fld map[int32][]float32             // map[int32][]float32
	Dd  map[int32]map[int32]interface{} // map[int32]map[int32]variant<int32, Corpus, float64, string, *Vector2>
	Fdd map[int32]map[int32]float32     // map[int32]map[int32]float32
	Nv  interface{}                     // variant<nil, int32>
	Lv  interface{}                     // variant<int32, []variant<float32, string>>
	Flv interface{}                     // variant<int32, []float32>
	Dv  interface{}                     // variant<int32, map[int32]variant<float32, string>>
	Fdv interface{}                     // variant<int32, map[int32]float32>
}

func (s *Fighter) ByteSize() (int, error) {
	return 0, nil
}

func (s *Fighter) Serialize(w io.Writer) error {
	return nil
}

func (s *Fighter) Deserialize(r io.Reader) error {
	return nil
}

func (s *Fighter) SerializeRPGParam(a0 *Fighter, a1 interface{}, a2 float64) (data string, err error) {
	return
}

func (s *Fighter) DeserializeRPGParam(data string) (a0 *Fighter, a1 interface{}, a2 float64, err error) {
	return
}

func (s *Fighter) SerializeRPGResult(a0 *Vector2) (data string, err error) {
	return
}

func (s *Fighter) DeserializeRPGResult(data string) (a0 *Vector2, err error) {
	return
}

func (s *Fighter) SerializeGPRParam(a0 map[int32]interface{}) (data string, err error) {
	return
}

func (s *Fighter) DeserializeGPRParam(data string) (a0 map[int32]interface{}, err error) {
	return
}

func (s *Fighter) SerializeGPRResult(a0 *Fighter, a1 int32) (data string, err error) {
	return
}

func (s *Fighter) DeserializeGPRResult(data string) (a0 *Fighter, a1 int32, err error) {
	return
}
