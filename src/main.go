package main

import "fmt"

func main() {
	f := Sequence()

	fmt.Println(f())
	fmt.Println(f())

	f = Sequence()
	fmt.Println(f())
	fmt.Println(f())
}

func Sequence() func() int {
	x := 0
	return func() int {
		x++
		return x*x
	}
}

