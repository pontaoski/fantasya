package main

import (
	"encoding/binary"
	"os"
)

type output struct{}

func (o output) Query() (id cell, version cell) {
	return 0, 0
}
func (o output) Invoke(v *vm) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(v.PopData()))
	_, err := os.Stdout.Write(bs)
	if err != nil {
		panic(err)
	}
}
