package main

func main() {
	v := vm{}
	v.devices = append(v.devices, output{})
	v.devices = append(v.devices, input{})
	v.devices = append(v.devices, goroutine{})
	err := v.loadImage("ngaImage")
	if err != nil {
		panic(err)
	}
	v.run()
}
