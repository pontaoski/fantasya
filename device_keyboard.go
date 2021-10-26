package main

import (
	"os"
)

type input struct{}

func (o input) Query() (id cell, version cell) {
	return 1, 0
}
func (o input) Invoke(v *vm) {
	b := make([]byte, 1)
	_, err := os.Stdin.Read(b)
	if err != nil {
		panic(err)
	}

	v.PushData(cell(b[0]))
}
