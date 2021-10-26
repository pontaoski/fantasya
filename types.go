package main

type inst int

type cell int32

const imageSize = 524288 * 4
const addresses = 2048
const stackDepth = 512

const (
	nop inst = iota
	lit
	dup
	drop
	swap
	push
	pop
	jump
	call
	ccall
	ret
	eq
	neq
	lt
	gt
	fetch
	store
	add
	sub
	mul
	divmod
	and
	or
	xor
	shift
	zret
	end
	ioEnum
	ioQuery
	ioInteract
	numOps
)

func (i inst) String() string {
	switch i {
	case nop:
		return "nop"
	case lit:
		return "lit"
	case dup:
		return "dup"
	case drop:
		return "drop"
	case swap:
		return "swap"
	case push:
		return "push"
	case pop:
		return "pop"
	case jump:
		return "jump"
	case call:
		return "call"
	case ccall:
		return "ccall"
	case ret:
		return "ret"
	case eq:
		return "eq"
	case neq:
		return "neq"
	case lt:
		return "lt"
	case gt:
		return "gt"
	case fetch:
		return "fetch"
	case store:
		return "store"
	case add:
		return "add"
	case sub:
		return "sub"
	case mul:
		return "mul"
	case divmod:
		return "divmod"
	case and:
		return "and"
	case or:
		return "or"
	case xor:
		return "xor"
	case shift:
		return "shift"
	case zret:
		return "zret"
	case end:
		return "end"
	case ioEnum:
		return "ioEnum"
	case ioQuery:
		return "ioQuery"
	case ioInteract:
		return "ioInteract"
	}
	return ""
}
