// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package tygo

import (
	"io"
)

const (
	WireTypeBits = 3
	WireTypeMask = (1 << WireTypeBits) - 1
)

type WireType uint8

const (
	WireVarint WireType = iota
	WireFixed64
	WireBytes
	WireStartGroup
	WireEndGroup
	WireFixed32
)

func MAKE_TAG(fieldNum uint32, wireType WireType) uint32 {
	return (fieldNum << WireTypeBits) | uint32(wireType)
}

func MAX_TAG(tag uint32) uint32 {
	return MAKE_TAG(tag, WireTypeMask)
}

func TAG_FIELD(tag uint32) uint32 {
	return tag >> WireTypeBits
}

func TAG_WIRE(tag uint32) uint32 {
	return tag & WireTypeMask
}

type ProtoBuf struct {
	offset int
	Buffer []byte
}

func ReadTag(reader io.Reader, cutoff uint32) (uint32, error) {
	return 0, nil
}

func SkipField(reader io.Reader, fieldNum uint32) error {
	return nil
}
