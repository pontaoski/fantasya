package main

type goroutine struct{}

func (o goroutine) Query() (id cell, version cell) {
	return 727, 0
}
func (o goroutine) Invoke(v *vm) {
	cmd := v.PopData()

	switch cmd {
	// fork
	case 0:
		new := *v
		go new.resume()
	}
}
