// Package c provides convenience helpers for C interop.
// It bridges C's null-terminated strings and raw pointers
// with So's string and slice types.
package c

import "unsafe"

//so:embed c.h
var c_h string

// Char represents a C char type.
//
//so:extern char
type Char byte

// ConstChar represents a C char type with a const modifier.
//
//so:extern so_const_char
type ConstChar byte

// Alignof returns the alignment of type T in bytes.
//
//	alignof(T)
//
//so:extern
func Alignof[T any]() int {
	var v T
	return int(unsafe.Alignof(v))
}

// Alloca allocates an array of the given length
// on the stack and returns a pointer to it.
//
//so:extern
func Alloca[T any](n int) *T {
	v := make([]T, n)
	return &v[0]
}

// Assert aborts the program with the given message
// if the condition is not true.
// If assertions are disabled, does nothing.
//
//	assert((cond) && msg)
//
//so:extern
func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

// Bytes wraps a raw byte pointer and length into a []byte without copying.
// If ptr is nil, returns nil.
//
//	(so_Slice){ptr, n, n}
//
//so:extern
func Bytes(ptr *byte, n int) []byte {
	if ptr == nil {
		return nil
	}
	return unsafe.Slice(ptr, n)
}

// CString converts a So string to a null-terminated C string.
// Allocates memory on the stack using alloca until the calling function returns.
//
//so:extern so_cstr nodecay
func CString(s string) *Char {
	return (*Char)(unsafe.StringData(s))
}

// PtrAdd adds offset bytes to a pointer and returns the result.
//
//	ptr + offset
//
//so:extern
func PtrAdd[T any](ptr *T, offset int) *T {
	raw := ptrVal(ptr)
	p := unsafe.Add(raw, offset)
	return (*T)(p)
}

// PtrAs casts a raw pointer (void*) to *T.
//
//	(T*)(ptr)
//
//so:extern
func PtrAs[T any](ptr any) *T {
	raw := ptrVal(ptr)
	return (*T)(raw)
}

// PtrAt returns a pointer to the element at the given index in an array or slice.
//
//	&ptr[index]
//
//so:extern
func PtrAt[T any](ptr *T, index int) *T {
	return PtrAdd(ptr, index*Sizeof[T]())
}

// Raw emits a raw block of C code.
//
//so:extern
func Raw(code string) { _ = code }

// Sizeof returns the size of type T in bytes.
//
//	sizeof(T)
//
//so:extern
func Sizeof[T any]() int {
	var v T
	return int(unsafe.Sizeof(v))
}

// Slice wraps a raw pointer and length into a []T without copying.
// If ptr is nil, returns nil.
//
//	(so_Slice){ptr, len, cap}
//
//so:extern
func Slice[T any](ptr *T, len int, cap int) []T {
	if ptr == nil {
		return nil
	}
	s := unsafe.Slice(ptr, cap)
	return s[:len]
}

// String converts a null-terminated C string to a So string.
// If ptr is nil, returns "".
//
//	(so_String){s, strlen(s)}
//
//so:extern
func String[T Char|ConstChar](ptr *T) string { _ = ptr; return "" }

// Val emits a typed C expression.
//
//so:extern
func Val[T any](expr string) T {
	var v T
	_ = expr
	return v
}

// Zero returns the zero value of type T.
//
//	{0}
//
//so:extern
func Zero[T any]() T {
	var v T
	return v
}

// ptrVal extracts a raw pointer from an interface containing any pointer type.
// For testing only: in C any is void*, so unwrapping is unnecessary.
//
//so:extern
func ptrVal(v any) unsafe.Pointer {
	type iface struct {
		_    unsafe.Pointer
		data unsafe.Pointer
	}
	return (*iface)(unsafe.Pointer(&v)).data
}
