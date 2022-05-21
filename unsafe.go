package umem

import "unsafe"

//go:linkname sysAlloc runtime.sysAlloc
func sysAlloc(n uintptr, sysStat *uint64) unsafe.Pointer

//go:linkname sysFree runtime.sysFree
func sysFree(v unsafe.Pointer, n uintptr, sysStat *uint64)
