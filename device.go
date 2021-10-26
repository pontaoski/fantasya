package main

type device interface {
	Query() (id cell, version cell)
	Invoke(v *vm)
}
