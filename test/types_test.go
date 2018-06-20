// Copyright 2017 - 2018 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package test

import (
	"github.com/ibelie/tygo"
	"testing"
)

func TestVector2(t *testing.T) {
	vd := &tygo.ProtoBuf{Buffer: make([]byte, v.ByteSize())}
	v.Serialize(vd)
	vd.Reset()
	v3 := &Vector2{}
	if err := v3.Deserialize(vd); err == nil {
		CompareVector2(t.Errorf, v, v3, "")
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
		CompareFighter(t.Errorf, fighter, fighter2)
	} else {
		t.Errorf("TestFighter Deserialize error: %v", err)
	}
}
