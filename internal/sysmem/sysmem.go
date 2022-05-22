package umem

import "unsafe"

//go:linkname SYSAllocOS runtime.sysAllocOS
func SYSAllocOS(n uintptr) unsafe.Pointer

//go:linkname SYSFreeOS runtime.sysFreeOS
func SYSFreeOS(v unsafe.Pointer, n uintptr)
