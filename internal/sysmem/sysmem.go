package umem

import (
	_ "runtime"
	"unsafe"
)

var memstat sysMemStat

func SYSAllocOS(n uintptr) unsafe.Pointer {
	return SYSAlloc_runtime(n, &memstat)
}

func SYSFreeOS(v unsafe.Pointer, n uintptr) {
	SYSFree_runtime(v, n, &memstat)
}

//go:linkname SYSAlloc_runtime runtime.sysAlloc
//go:linkname SYSFree_runtime runtime.sysFree

type sysMemStat uint64

func SYSAlloc_runtime(n uintptr, sysStat *sysMemStat) unsafe.Pointer
func SYSFree_runtime(v unsafe.Pointer, n uintptr, sysStat *sysMemStat)
