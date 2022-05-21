package umem

import "unsafe"

//go:linkname sysAllocOS runtime.sysAllocOS
func sysAllocOS(n uintptr) unsafe.Pointer

//go:linkname sysFreeOS runtime.sysFreeOS
func sysFreeOS(v unsafe.Pointer, n uintptr)
