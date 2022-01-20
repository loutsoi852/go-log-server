package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

type Thing struct {
	A, B int
}

func (t *Thing) addThis(i int) {
	fmt.Println(t.A+i, t.B+i)
}

func test() {
	t := Thing{4, 9}
	t.addThis(3)
}

func main() {
	v := Vertex{3, 4}
	v.Scale(10)
	fmt.Println(v.Abs())

	test()
}
