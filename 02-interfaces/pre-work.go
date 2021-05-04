package main

import (
	"fmt"
	"strconv"
	"unsafe"
)

type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype
type tflag uint8

type iface struct { // `iface`
	tab *struct { // `itab`
		inter *struct { // `interfacetype`
			typ struct { // `_type`
				size       uintptr
				ptrdata    uintptr
				hash       uint32
				tflag      tflag
				align      uint8
				fieldalign uint8
				kind       uint8
				equal      func(unsafe.Pointer, unsafe.Pointer) bool
				// gcdata stores the GC type data for the garbage collector.
				// If the KindGCProg bit is set in kind, gcdata is a GC program.
				// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
				gcdata    *byte
				str       nameOff
				ptrToThis typeOff
			}
			pkgpath struct {
				bytes *byte
			}
			mhdr []struct { // `imethod`
				name nameOff
				ityp typeOff
			}
		}
		_type *struct { // `_type`
			size       uintptr
			ptrdata    uintptr
			hash       uint32
			tflag      tflag
			align      uint8
			fieldalign uint8
			kind       uint8
			equal      func(unsafe.Pointer, unsafe.Pointer) bool
			// gcdata stores the GC type data for the garbage collector.
			// If the KindGCProg bit is set in kind, gcdata is a GC program.
			// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
			gcdata    *byte
			str       nameOff
			ptrToThis typeOff
		}
		hash uint32
		_    [4]byte
		fun  [1]uintptr
	}
	data unsafe.Pointer
}

func getWrappedInt(i interface{}) int {
	return *(*int)((*iface)(unsafe.Pointer(&i)).data)
}

type Binary uint64

func (i Binary) String() string {
	return strconv.FormatUint(i.Get(), 2)
}

func (i Binary) String2() string {
	return fmt.Sprintf("Here's string representation %v", strconv.FormatUint(i.Get(), 2))
}

func (i Binary) Get() uint64 {
	return uint64(i)
}

func main() {
	// Given an interface{} variable that holds an int value, write a function that extracts the int value without using a type assertion or type switch.
	i := interface{}(7)
	val := getWrappedInt(i)
	fmt.Printf("Int value is : %d\n", val)

	// Given an arbitrary interface value, write a function that iterates through the corresponding itable and prints out information about methods
	// You can start by printing just the number of methods, but the eventual goal is for you to explore the underlying representations in an open-ended way.
	type Stringer interface {
		String() string
		String2() string
	}

	funcs := (*iface)(unsafe.Pointer(&i)).tab.fun
	unsafe.Pointer(funcs)
	fmt.Printf("Number of funcs: %d", len(funcs))
}
