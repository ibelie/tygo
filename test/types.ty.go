// Generated by tygo.  DO NOT EDIT!

package main

import "fmt"
import "github.com/ibelie/tygo"

type Corpus uint

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

func (i Corpus) ByteSize() (size int) {
	if i != 0 {
		size += tygo.SizeVarint(uint64(i))
	}
	return
}

func (i Corpus) Serialize(output *tygo.ProtoBuf) {
	output.WriteVarint(uint64(i))
}

func (i *Corpus) Deserialize(input *tygo.ProtoBuf) (err error) {
	x, err := input.ReadVarint()
	*i = Corpus(x)
	return
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

func (s *Vector2) MaxFieldNum() int {
	return 6
}

func (s *Vector2) ByteSize() (size int) {
	if s != nil {
		// property: s.X
		// type: float32
		if s.X != 0 {
			size += 1 + 4
		}

		// property: s.Y
		// type: fixedpoint<1, -10>
		if s.Y != -10 {
			size += 1 + tygo.SizeVarint(uint64((s.Y - -10) * 10))
		}

		// property: s.B
		// type: bytes
		if len(s.B) > 0 {
			l := len([]byte(s.B))
			size += 1 + tygo.SizeVarint(uint64(l)) + l
		}

		// property: s.S
		// type: string
		if len(s.S) > 0 {
			l := len([]byte(s.S))
			size += 1 + tygo.SizeVarint(uint64(l)) + l
		}

		// property: s.E
		// type: Corpus
		if s.E != 0 {
			size += 1 + tygo.SizeVarint(uint64(s.E))
		}

		// property: s.P
		// type: *GoType
		if s.P != nil {
			tSize := s.P.ByteSize()
			size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

	}
	return
}

func (s *Vector2) Serialize(output *tygo.ProtoBuf) {
}

func (s *Vector2) Deserialize(input *tygo.ProtoBuf) (err error) {
	return
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

func (s *Fighter_Part1) MaxFieldNum() int {
	return 11
}

func (s *Fighter_Part1) ByteSize() (size int) {
	if s != nil {
		// property: s.Pos
		// type: *Vector2
		if s.Pos != nil {
			tSize := s.Pos.ByteSize()
			size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.IsAwake
		// type: bool
		if s.IsAwake {
			size += 1 + 1
		}

		// property: s.Hp
		// type: float32
		if s.Hp != 0 {
			size += 1 + 4
		}

		// property: s.Poss
		// type: map[int32]*Vector2
		if len(s.Poss) > 0 {
			for k, v := range s.Poss {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: *Vector2
				if v != nil {
					tSizee := v.ByteSize()
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
				size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Posi
		// type: map[int32]float32
		if len(s.Posi) > 0 {
			for k, v := range s.Posi {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: float32
				if v != 0 {
					tSize += 1 + 4
				}
				size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Posl
		// type: []*Vector2
		if len(s.Posl) > 0 {
			for _, e := range s.Posl {
				// list element
				// type: *Vector2
				if e != nil {
					tSize := e.ByteSize()
					size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
				} else {
					size += 1 + 1
				}
			}
		}

		// property: s.Posll
		// type: [][]*Vector2
		if len(s.Posll) > 0 {
			for _, e := range s.Posll {
				tSize := 0
				// list element
				// type: []*Vector2
				if len(e) > 0 {
					for _, e := range e {
						// list element
						// type: *Vector2
						if e != nil {
							tSizee := e.ByteSize()
							tSize += tygo.SizeVarint(uint64(tSizee)) + tSizee
						} else {
							tSize += 1
						}
					}
				}
				size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Pyl
		// type: []*GoType
		if len(s.Pyl) > 0 {
			for _, e := range s.Pyl {
				// list element
				// type: *GoType
				if e != nil {
					tSize := e.ByteSize()
					size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
				} else {
					size += 1 + 1
				}
			}
		}

		// property: s.Pyd
		// type: map[int32]*GoType
		if len(s.Pyd) > 0 {
			for k, v := range s.Pyd {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: *GoType
				if v != nil {
					tSizee := v.ByteSize()
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
				size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Pyv1
		// type: variant<int32, *GoType>
		if s.Pyv1 != nil {
			tSize := 0
			switch v := s.Pyv1.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: *GoType
			case *GoType:
				// type: *GoType
				{
					tSizee := v.ByteSize()
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, *GoType>: %v", v))
			}
			size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.Pyv2
		// type: variant<int32, *GoType>
		if s.Pyv2 != nil {
			tSize := 0
			switch v := s.Pyv2.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: *GoType
			case *GoType:
				// type: *GoType
				{
					tSizee := v.ByteSize()
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, *GoType>: %v", v))
			}
			size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

	}
	return
}

func (s *Fighter_Part1) Serialize(output *tygo.ProtoBuf) {
}

func (s *Fighter_Part1) Deserialize(input *tygo.ProtoBuf) (err error) {
	return
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

func (s *Fighter_Part2) MaxFieldNum() int {
	return 19
}

func (s *Fighter_Part2) ByteSize() (size int) {
	if s != nil {
		size += s.Fighter_Part1.ByteSize()
		// property: s.Fl
		// type: []float32
		if len(s.Fl) > 0 {
			tSize := len(s.Fl) * 4
			size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.Bl
		// type: []bytes
		if len(s.Bl) > 0 {
			for _, e := range s.Bl {
				// list element
				// type: bytes
				if len(e) > 0 {
					l := len([]byte(e))
					size += 1 + tygo.SizeVarint(uint64(l)) + l
				} else {
					size += 1 + 1
				}
			}
		}

		// property: s.Sl
		// type: []string
		if len(s.Sl) > 0 {
			for _, e := range s.Sl {
				// list element
				// type: string
				if len(e) > 0 {
					l := len([]byte(e))
					size += 1 + tygo.SizeVarint(uint64(l)) + l
				} else {
					size += 1 + 1
				}
			}
		}

		// property: s.Bd
		// type: map[string]bytes
		if len(s.Bd) > 0 {
			for k, v := range s.Bd {
				tSize := 0
				// dict key
				// type: string
				if len(k) > 0 {
					l := len([]byte(k))
					tSize += 1 + tygo.SizeVarint(uint64(l)) + l
				}
				// dict value
				// type: bytes
				if len(v) > 0 {
					l := len([]byte(v))
					tSize += 1 + tygo.SizeVarint(uint64(l)) + l
				}
				size += 1 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Sd
		// type: map[int32]string
		if len(s.Sd) > 0 {
			for k, v := range s.Sd {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: string
				if len(v) > 0 {
					l := len([]byte(v))
					tSize += 1 + tygo.SizeVarint(uint64(l)) + l
				}
				size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.El
		// type: []Corpus
		if len(s.El) > 0 {
			tSize := 0
			for _, e := range s.El {
				// list element
				// type: Corpus
				tSize += tygo.SizeVarint(uint64(e))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.Ed
		// type: map[int32]Corpus
		if len(s.Ed) > 0 {
			for k, v := range s.Ed {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: Corpus
				if v != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(v))
				}
				size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Ll
		// type: [][]float32
		if len(s.Ll) > 0 {
			for _, e := range s.Ll {
				tSize := 0
				// list element
				// type: []float32
				if len(e) > 0 {
					tSizee := len(e) * 4
					tSize += tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
				size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

	}
	return
}

func (s *Fighter_Part2) Serialize(output *tygo.ProtoBuf) {
}

func (s *Fighter_Part2) Deserialize(input *tygo.ProtoBuf) (err error) {
	return
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

func (s *Fighter) MaxFieldNum() int {
	return 35
}

func (s *Fighter) ByteSize() (size int) {
	if s != nil {
		size += s.Fighter_Part2.ByteSize()
		// property: s.V0
		// type: variant<int32, float32, bytes, *Vector2>
		if s.V0 != nil {
			tSize := 0
			switch v := s.V0.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: float32
			case float32:
				// type: float32
				tSize += 1 + 4
			// variant type: bytes
			case []byte:
				// type: bytes
				{
					l := len([]byte(v))
					tSize += 1 + tygo.SizeVarint(uint64(l)) + l
				}
			// variant type: *Vector2
			case *Vector2:
				// type: *Vector2
				{
					tSizee := v.ByteSize()
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// addition type: float64 -> float32
			case float64:
				tSize += 5
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, float32, bytes, *Vector2>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.V1
		// type: variant<int32, float32, bytes, *Vector2>
		if s.V1 != nil {
			tSize := 0
			switch v := s.V1.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: float32
			case float32:
				// type: float32
				tSize += 1 + 4
			// variant type: bytes
			case []byte:
				// type: bytes
				{
					l := len([]byte(v))
					tSize += 1 + tygo.SizeVarint(uint64(l)) + l
				}
			// variant type: *Vector2
			case *Vector2:
				// type: *Vector2
				{
					tSizee := v.ByteSize()
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// addition type: float64 -> float32
			case float64:
				tSize += 5
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, float32, bytes, *Vector2>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.V2
		// type: variant<int32, float32, bytes, *Vector2>
		if s.V2 != nil {
			tSize := 0
			switch v := s.V2.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: float32
			case float32:
				// type: float32
				tSize += 1 + 4
			// variant type: bytes
			case []byte:
				// type: bytes
				{
					l := len([]byte(v))
					tSize += 1 + tygo.SizeVarint(uint64(l)) + l
				}
			// variant type: *Vector2
			case *Vector2:
				// type: *Vector2
				{
					tSizee := v.ByteSize()
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// addition type: float64 -> float32
			case float64:
				tSize += 5
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, float32, bytes, *Vector2>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.V3
		// type: variant<int32, float32, bytes, *Vector2>
		if s.V3 != nil {
			tSize := 0
			switch v := s.V3.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: float32
			case float32:
				// type: float32
				tSize += 1 + 4
			// variant type: bytes
			case []byte:
				// type: bytes
				{
					l := len([]byte(v))
					tSize += 1 + tygo.SizeVarint(uint64(l)) + l
				}
			// variant type: *Vector2
			case *Vector2:
				// type: *Vector2
				{
					tSizee := v.ByteSize()
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// addition type: float64 -> float32
			case float64:
				tSize += 5
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, float32, bytes, *Vector2>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.V4
		// type: variant<int32, float32, bytes, *Vector2>
		if s.V4 != nil {
			tSize := 0
			switch v := s.V4.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: float32
			case float32:
				// type: float32
				tSize += 1 + 4
			// variant type: bytes
			case []byte:
				// type: bytes
				{
					l := len([]byte(v))
					tSize += 1 + tygo.SizeVarint(uint64(l)) + l
				}
			// variant type: *Vector2
			case *Vector2:
				// type: *Vector2
				{
					tSizee := v.ByteSize()
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// addition type: float64 -> float32
			case float64:
				tSize += 5
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, float32, bytes, *Vector2>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.Vl
		// type: []variant<int32, fixedpoint<3, 0>, string, *Vector2>
		if len(s.Vl) > 0 {
			for _, e := range s.Vl {
				// list element
				// type: variant<int32, fixedpoint<3, 0>, string, *Vector2>
				if e != nil {
					tSize := 0
					switch v := e.(type) {
					// variant type: int32
					case int32:
						// type: int32
						tSize += 1 + tygo.SizeVarint(uint64(v))
					// variant type: fixedpoint<3, 0>
					case float64:
						// type: fixedpoint<3, 0>
						tSize += 1 + tygo.SizeVarint(uint64((v - 0) * 1000))
					// variant type: string
					case string:
						// type: string
						{
							l := len([]byte(v))
							tSize += 1 + tygo.SizeVarint(uint64(l)) + l
						}
					// variant type: *Vector2
					case *Vector2:
						// type: *Vector2
						{
							tSizee := v.ByteSize()
							tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
						}
					// addition type: int
					case int:
						tSize += 1 + tygo.SizeVarint(uint64(v))
					default:
						panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, fixedpoint<3, 0>, string, *Vector2>: %v", v))
					}
					size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
				} else {
					size += 2 + 1
				}
			}
		}

		// property: s.Vd
		// type: map[int32]variant<Corpus, float64, string, *Vector2>
		if len(s.Vd) > 0 {
			for k, v := range s.Vd {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: variant<Corpus, float64, string, *Vector2>
				if v != nil {
					tSizee := 0
					switch v := v.(type) {
					// variant type: Corpus
					case Corpus:
						// type: Corpus
						tSizee += 1 + tygo.SizeVarint(uint64(v))
					// variant type: float64
					case float64:
						// type: float64
						tSizee += 1 + 8
					// variant type: string
					case string:
						// type: string
						{
							l := len([]byte(v))
							tSizee += 1 + tygo.SizeVarint(uint64(l)) + l
						}
					// variant type: *Vector2
					case *Vector2:
						// type: *Vector2
						{
							tSizeee := v.ByteSize()
							tSizee += 1 + tygo.SizeVarint(uint64(tSizeee)) + tSizeee
						}
					// addition type: int -> float64
					case int:
						tSizee += 9
					default:
						panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<Corpus, float64, string, *Vector2>: %v", v))
					}
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
				size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Ld
		// type: map[int32][]variant<Corpus, float64, string, *Vector2>
		if len(s.Ld) > 0 {
			for k, v := range s.Ld {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: []variant<Corpus, float64, string, *Vector2>
				if len(v) > 0 {
					for _, e := range v {
						// list element
						// type: variant<Corpus, float64, string, *Vector2>
						if e != nil {
							tSizee := 0
							switch v := e.(type) {
							// variant type: Corpus
							case Corpus:
								// type: Corpus
								tSizee += 1 + tygo.SizeVarint(uint64(v))
							// variant type: float64
							case float64:
								// type: float64
								tSizee += 1 + 8
							// variant type: string
							case string:
								// type: string
								{
									l := len([]byte(v))
									tSizee += 1 + tygo.SizeVarint(uint64(l)) + l
								}
							// variant type: *Vector2
							case *Vector2:
								// type: *Vector2
								{
									tSizeee := v.ByteSize()
									tSizee += 1 + tygo.SizeVarint(uint64(tSizeee)) + tSizeee
								}
							// addition type: int -> float64
							case int:
								tSizee += 9
							default:
								panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<Corpus, float64, string, *Vector2>: %v", v))
							}
							tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
						} else {
							tSize += 1 + 1
						}
					}
				}
				size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Fld
		// type: map[int32][]float32
		if len(s.Fld) > 0 {
			for k, v := range s.Fld {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: []float32
				if len(v) > 0 {
					tSizee := len(v) * 4
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
				size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Dd
		// type: map[int32]map[int32]variant<int32, Corpus, float64, string, *Vector2>
		if len(s.Dd) > 0 {
			for k, v := range s.Dd {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: map[int32]variant<int32, Corpus, float64, string, *Vector2>
				if len(v) > 0 {
					for k, v := range v {
						tSizee := 0
						// dict key
						// type: int32
						if k != 0 {
							tSizee += 1 + tygo.SizeVarint(uint64(k))
						}
						// dict value
						// type: variant<int32, Corpus, float64, string, *Vector2>
						if v != nil {
							tSizeee := 0
							switch v := v.(type) {
							// variant type: int32
							case int32:
								// type: int32
								tSizeee += 1 + tygo.SizeVarint(uint64(v))
							// variant type: Corpus
							case Corpus:
								// type: Corpus
								tSizeee += 1 + tygo.SizeVarint(uint64(v))
							// variant type: float64
							case float64:
								// type: float64
								tSizeee += 1 + 8
							// variant type: string
							case string:
								// type: string
								{
									l := len([]byte(v))
									tSizeee += 1 + tygo.SizeVarint(uint64(l)) + l
								}
							// variant type: *Vector2
							case *Vector2:
								// type: *Vector2
								{
									tSizeeee := v.ByteSize()
									tSizeee += 1 + tygo.SizeVarint(uint64(tSizeeee)) + tSizeeee
								}
							// addition type: int
							case int:
								tSizeee += 1 + tygo.SizeVarint(uint64(v))
							default:
								panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, Corpus, float64, string, *Vector2>: %v", v))
							}
							tSizee += 1 + tygo.SizeVarint(uint64(tSizeee)) + tSizeee
						}
						tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
					}
				}
				size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Fdd
		// type: map[int32]map[int32]float32
		if len(s.Fdd) > 0 {
			for k, v := range s.Fdd {
				tSize := 0
				// dict key
				// type: int32
				if k != 0 {
					tSize += 1 + tygo.SizeVarint(uint64(k))
				}
				// dict value
				// type: map[int32]float32
				if len(v) > 0 {
					for k, v := range v {
						tSizee := 0
						// dict key
						// type: int32
						if k != 0 {
							tSizee += 1 + tygo.SizeVarint(uint64(k))
						}
						// dict value
						// type: float32
						if v != 0 {
							tSizee += 1 + 4
						}
						tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
					}
				}
				size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
			}
		}

		// property: s.Nv
		// type: variant<nil, int32>
		if s.Nv != nil {
			tSize := 0
			switch v := s.Nv.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<nil, int32>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.Lv
		// type: variant<int32, []variant<float32, string>>
		if s.Lv != nil {
			tSize := 0
			switch v := s.Lv.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: []variant<float32, string>
			case []interface{}:
				// type: []variant<float32, string>
				if len(v) > 0 {
					for _, e := range v {
						// list element
						// type: variant<float32, string>
						if e != nil {
							tSizee := 0
							switch v := e.(type) {
							// variant type: float32
							case float32:
								// type: float32
								tSizee += 1 + 4
							// variant type: string
							case string:
								// type: string
								{
									l := len([]byte(v))
									tSizee += 1 + tygo.SizeVarint(uint64(l)) + l
								}
							// addition type: int -> float32
							case int:
								tSizee += 5
							// addition type: float64 -> float32
							case float64:
								tSizee += 5
							default:
								panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<float32, string>: %v", v))
							}
							tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
						} else {
							tSize += 1 + 1
						}
					}
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, []variant<float32, string>>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.Flv
		// type: variant<int32, []float32>
		if s.Flv != nil {
			tSize := 0
			switch v := s.Flv.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: []float32
			case []float32:
				// type: []float32
				if len(v) > 0 {
					tSizee := len(v) * 4
					tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, []float32>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.Dv
		// type: variant<int32, map[int32]variant<float32, string>>
		if s.Dv != nil {
			tSize := 0
			switch v := s.Dv.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: map[int32]variant<float32, string>
			case map[int32]interface{}:
				// type: map[int32]variant<float32, string>
				if len(v) > 0 {
					for k, v := range v {
						tSizee := 0
						// dict key
						// type: int32
						if k != 0 {
							tSizee += 1 + tygo.SizeVarint(uint64(k))
						}
						// dict value
						// type: variant<float32, string>
						if v != nil {
							tSizeee := 0
							switch v := v.(type) {
							// variant type: float32
							case float32:
								// type: float32
								tSizeee += 1 + 4
							// variant type: string
							case string:
								// type: string
								{
									l := len([]byte(v))
									tSizeee += 1 + tygo.SizeVarint(uint64(l)) + l
								}
							// addition type: int -> float32
							case int:
								tSizeee += 5
							// addition type: float64 -> float32
							case float64:
								tSizeee += 5
							default:
								panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<float32, string>: %v", v))
							}
							tSizee += 1 + tygo.SizeVarint(uint64(tSizeee)) + tSizeee
						}
						tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
					}
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, map[int32]variant<float32, string>>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

		// property: s.Fdv
		// type: variant<int32, map[int32]float32>
		if s.Fdv != nil {
			tSize := 0
			switch v := s.Fdv.(type) {
			// variant type: int32
			case int32:
				// type: int32
				tSize += 1 + tygo.SizeVarint(uint64(v))
			// variant type: map[int32]float32
			case map[int32]float32:
				// type: map[int32]float32
				if len(v) > 0 {
					for k, v := range v {
						tSizee := 0
						// dict key
						// type: int32
						if k != 0 {
							tSizee += 1 + tygo.SizeVarint(uint64(k))
						}
						// dict value
						// type: float32
						if v != 0 {
							tSizee += 1 + 4
						}
						tSize += 1 + tygo.SizeVarint(uint64(tSizee)) + tSizee
					}
				}
			// addition type: int
			case int:
				tSize += 1 + tygo.SizeVarint(uint64(v))
			default:
				panic(fmt.Sprintf("[Tygo][Variant] Unexpect type for variant<int32, map[int32]float32>: %v", v))
			}
			size += 2 + tygo.SizeVarint(uint64(tSize)) + tSize
		}

	}
	return
}

func (s *Fighter) Serialize(output *tygo.ProtoBuf) {
}

func (s *Fighter) Deserialize(input *tygo.ProtoBuf) (err error) {
	return
}

// RPG Params(a0: *Fighter, a1: variant<nil, int32>, a2: fixedpoint<3, 0>)
func (s *Fighter) SerializeRPGParam(a0 *Fighter, a1 interface{}, a2 float64) (data []byte) {
	size := 0
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	return
}

// RPG Params(a0: *Fighter, a1: variant<nil, int32>, a2: fixedpoint<3, 0>)
func (s *Fighter) DeserializeRPGParam(data []byte) (a0 *Fighter, a1 interface{}, a2 float64, err error) {
	return
}

// RPG Results(a0: *Vector2)
func (s *Fighter) SerializeRPGResult(a0 *Vector2) (data []byte) {
	size := 0
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	return
}

// RPG Results(a0: *Vector2)
func (s *Fighter) DeserializeRPGResult(data []byte) (a0 *Vector2, err error) {
	return
}

// GPR Params(a0: map[int32]variant<Corpus, float64, string, *Vector2>)
func (s *Fighter) SerializeGPRParam(a0 map[int32]interface{}) (data []byte) {
	size := 0
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	return
}

// GPR Params(a0: map[int32]variant<Corpus, float64, string, *Vector2>)
func (s *Fighter) DeserializeGPRParam(data []byte) (a0 map[int32]interface{}, err error) {
	return
}

// GPR Results(a0: *Fighter, a1: int32)
func (s *Fighter) SerializeGPRResult(a0 *Fighter, a1 int32) (data []byte) {
	size := 0
	if size <= 0 {
		return
	}
	data = make([]byte, size)
	return
}

// GPR Results(a0: *Fighter, a1: int32)
func (s *Fighter) DeserializeGPRResult(data []byte) (a0 *Fighter, a1 int32, err error) {
	return
}
