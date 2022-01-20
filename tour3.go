package main

import "fmt"

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func main() {
	q := []int{2, 3, 5, 7, 11, 13}
	fmt.Println(q)

	r := []bool{true, false, true, true, false, true}
	fmt.Println(r)

	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Println(s)

	s1 := []int{2, 3, 5, 7, 11, 13}
	s1 = s1[1:4]
	fmt.Println(s1)
	s1 = s1[:2]
	fmt.Println(s1)
	s1 = s1[1:]
	fmt.Println(s1)

	s2 := []int{2, 3, 5, 7, 11, 13}
	printSlice(s2)
	// Slice the slice to give it zero length.
	s2 = s2[0:]
	printSlice(s2)
	// Extend its length.
	s2 = s2[4:]
	printSlice(s2)
	// // Drop its first two values.
	s2 = s2[2:]
	printSlice(s2)

	var s3 []int
	fmt.Println(s3, len(s3), cap(s3))
	if s3 == nil {
		fmt.Println("nil!")
	}

}
