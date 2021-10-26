package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type vm struct {
	sp cell
	rp cell
	ip cell

	data    [stackDepth]cell
	address [addresses]cell
	memory  [imageSize + 1]cell

	devices []device
}

func (v *vm) tos() *cell {
	if v.sp < 0 {
		println("stack underflow")
	}
	return &v.data[v.sp]
}

func (v *vm) nos() *cell {
	if v.sp < 0 {
		println("stack underflow")
	}
	return &v.data[v.sp-1]
}

func (v *vm) tors() *cell {
	if v.sp < 0 {
		println("return stack underflow")
	}
	return &v.address[v.rp]
}

func (v *vm) PushData(c cell) {
	v.sp++
	*v.tos() = c
}

func (v *vm) PopData() cell {
	a := *v.tos()
	v.sp--
	return a
}

func (v *vm) PushReturn(c cell) {
	v.rp++
	*v.tors() = c
}

func (v *vm) PopReturn() cell {
	a := *v.tors()
	v.rp--
	return a
}

func (v *vm) loadImage(s string) error {
	imageData, err := ioutil.ReadFile(s)
	if err != nil {
		return fmt.Errorf("failed to read image file: %+w", err)
	}

	buf := bytes.NewReader(imageData)

	nowVar := cell(0)
	i := 0
	for {
		err := binary.Read(buf, binary.LittleEndian, &nowVar)
		v.memory[i] = nowVar
		i++

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}

//go:inline
func flag(b bool) cell {
	if b {
		return -1
	} else {
		return 0
	}
}

//go:inline
func (v *vm) procInst(i inst) {
	switch i {
	case nop:
	case lit:
		v.sp++
		v.ip++
		*v.tos() = v.memory[v.ip]
	case dup:
		v.sp++
		*v.tos() = *v.nos()
	case drop:
		*v.tos() = 0
		v.sp--
		if v.sp < 0 {
			v.ip = imageSize
		}
	case swap:
		a := *v.tos()
		*v.tos() = *v.nos()
		*v.nos() = a
	case push:
		v.rp++
		*v.tors() = *v.tos()
		v.procInst(drop)
	case pop:
		v.sp++
		*v.tos() = *v.tors()
		v.rp--
	case jump:
		v.ip = *v.tos() - 1
		v.procInst(drop)
	case call:
		v.rp++
		*v.tors() = v.ip
		v.ip = *v.tos() - 1
		v.procInst(drop)
	case ccall:
		a := *v.tos()
		v.procInst(drop)
		b := *v.tos()
		v.procInst(drop)
		if b != 0 {
			v.rp++
			*v.tors() = v.ip
			v.ip = a - 1
		}
	case ret:
		v.ip = *v.tors()
		v.rp--
	case eq:
		*v.nos() = flag(*v.nos() == *v.tos())
		v.procInst(drop)
	case neq:
		*v.nos() = flag(*v.nos() != *v.tos())
		v.procInst(drop)
	case lt:
		*v.nos() = flag(*v.nos() < *v.tos())
		v.procInst(drop)
	case gt:
		*v.nos() = flag(*v.nos() > *v.tos())
		v.procInst(drop)
	case fetch:
		switch *v.tos() {
		case -1:
			*v.tos() = v.sp - 1
		case -2:
			*v.tos() = v.rp
		case -3:
			*v.tos() = imageSize
		case -4:
			*v.tos() = math.MinInt32
		case -5:
			*v.tos() = math.MaxInt32
		default:
			*v.tos() = v.memory[*v.tos()]
		}
	case store:
		v.memory[*v.tos()] = *v.nos()
		v.procInst(drop)
		v.procInst(drop)
	case add:
		*v.nos() += *v.tos()
		v.procInst(drop)
	case sub:
		*v.nos() -= *v.tos()
		v.procInst(drop)
	case mul:
		*v.nos() *= *v.tos()
		v.procInst(drop)
	case divmod:
		a := *v.tos()
		b := *v.nos()

		*v.tos() = b / a
		*v.nos() = b % a
	case and:
		*v.nos() = *v.tos() & *v.nos()
		v.procInst(drop)
	case or:
		*v.nos() = *v.tos() | *v.nos()
		v.procInst(drop)
	case xor:
		*v.nos() = *v.tos() ^ *v.nos()
		v.procInst(drop)
	case shift:
		y := *v.tos()
		x := *v.nos()

		if *v.tos() < 0 {
			*v.nos() = *v.nos() << (*v.tos() * -1)
		} else if x < 0 && y > 0 {
			*v.nos() = x>>y | ^(^0 >> y)
		} else {
			*v.nos() = x >> y
		}
		v.procInst(drop)
	case zret:
		if *v.tos() == 0 {
			v.procInst(drop)
			v.ip = *v.tors()
			v.rp--
		}
	case end:
		v.ip = imageSize
	case ioEnum:
		v.sp++
		*v.tos() = cell(len(v.devices))
	case ioQuery:
		dev := *v.tos()
		v.procInst(drop)
		id, ver := v.devices[dev].Query()
		v.PushData(ver)
		v.PushData(id)
	case ioInteract:
		dev := *v.tos()
		v.procInst(drop)
		for _, it := range v.devices {
			if id, _ := it.Query(); id == dev {
				it.Invoke(v)
				return
			}
		}
		panic("unhandled device " + fmt.Sprint(dev))
	default:
		panic("unhandled instruction")
	}
}

func validatePackedOpcode(opcode cell) int {
	raw := opcode
	valid := -1

	for i := 0; i < 4; i++ {
		current := raw & 0xff
		if !(current >= 0 && current <= 26) {
			valid = 0
		}

		raw = raw >> 8
	}

	return valid
}

func (v *vm) processPackedInst(opcode cell) {
	raw := opcode
	i := 0
	for i = 0; i < 4; i++ {
		v.procInst(inst(raw & 0xFF))
		raw = raw >> 8
	}
}

func (v *vm) resume() {
	for v.ip < imageSize {
		opcode := v.memory[v.ip]

		if validatePackedOpcode(opcode) != 0 {
			v.processPackedInst(opcode)
		} else if opcode >= 0 && opcode < cell(numOps) {
			v.procInst(inst(opcode))
		} else {
			fmt.Printf("Invalid instruction at %d, opcode %d", v.ip, opcode)
			return
		}

		v.ip++
	}
}

func (v *vm) run() {
	v.ip = 0

	v.resume()

	for i := cell(1); i <= v.sp; i++ {
		fmt.Printf("%d ", v.data[i])
	}
	fmt.Printf("\n")
}
