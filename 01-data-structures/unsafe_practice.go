package main

import (
	"fmt"
	"math"
	"unsafe"
)

func sizeOf() {
	i := int(5)
	j := int32(5)
	fmt.Println(unsafe.Sizeof(i)) // "8" -> 64 bytes
	fmt.Println(unsafe.Sizeof(j)) // "4" -> 32 bytes
	// Sizeof reports only the size of the fixed part of each data structure,
	// like the pointer and length of a string, but not indirect parts like the contents of the string
	fmt.Println(unsafe.Sizeof("abcdefgh")) // "16"
	fmt.Println(unsafe.Alignof(true))      // 1
}

// Given a float64, return a uint64 with the same binary representation
func float64Touint64(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

func sameLocationStrings(s1 string, s2 string) bool {
	// dereference string pointers to get the underlying data array addresses and compare them
	p1 := *(*uintptr)(unsafe.Pointer(&s1))
	p2 := *(*uintptr)(unsafe.Pointer(&s2))

	return p1 == p2
}

func sumIntSliceValues(slice []int) int {
	// determine int size, it's machine-dependent
	intSize := unsafe.Sizeof(int(1))
	fmt.Printf("Int size: %d\n", intSize)

	// determine pointer size, it's machine-dependent
	pointerSize := unsafe.Sizeof(&slice)
	fmt.Printf("Pointer size: %d\n", pointerSize)

	// determine slice size
	sliceSize := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&slice)) + pointerSize))
	fmt.Printf("Slice size: %d\n", sliceSize)

	sum := 0
	for i := 0; i < sliceSize; i++ {
		// determine uintptr address of struct start to perform integer arithmetic,
		// then dereference to get the [i]th element of the underlying array
		sum += *(*int)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&slice)) + intSize*uintptr(i)))
	}

	return sum
}

func main() {

	// 1. Given a float64, return a uint64 with the same binary representation
	f1 := 1.23
	u1 := float64Touint64(f1)
	fmt.Printf("%064b\n", u1)
	// verify
	fmt.Printf("%064b\n", math.Float64bits(f1))

	// 2. Given two strings, return a boolean that indicates whether the underlying string data is stored at the same memory location.
	s1 := "abcde"
	s2 := "abcdek"
	fmt.Printf("Same location? - %v\n", sameLocationStrings(s1, s2))

	s3 := "abbb"
	s4 := s3[:]

	fmt.Printf("Same location? - %v\n", sameLocationStrings(s3, s4))

	// 3. Given an []int slice, return the sum of values in the slice without using range or the [] operator.
	sl := []int{12, 3, 11, 3, 10}
	fmt.Printf("Sum of elements: %d", sumIntSliceValues(sl))
}
