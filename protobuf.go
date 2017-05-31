// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

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
	return (fieldNum << WireTypeBits) | wireType
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

func ReadTag(reader io.Reader, cutoff uint32) (uint32, error) {

}

func ReadVarint32(reader io.Reader) (uint32, error) {

}

func ReadVarint64(reader io.Reader) (uint64, error) {

}

func Read64(reader io.Reader) (uint64, error) {

}

func Read32(reader io.Reader) (uint32, error) {

}

func ReadByte(reader io.Reader) (byte, error) {

}

func SkipField(reader io.Reader, fieldNum uint32) error {

}

func WriteByte(writer io.Writer, data byte) error {

}

func WriteTag(writer io.Writer, data uint32) error {

}

func Write32(writer io.Writer, data uint32) error {

}

func Write64(writer io.Writer, data uint64) error {

}
