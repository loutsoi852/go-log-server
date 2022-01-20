package main

import (
	"fmt"
	"math"
)

type MyFloat float64
type MyInt int

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func (i MyInt) Test() {
	if i > 2 {
		fmt.Println("over 2")
	} else {
		fmt.Println("not over 2")
	}
}

func main() {
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())

	i := MyInt(3)
	i.Test()
}
