package arena

import (
	"os"
	"reflect"
	"runtime"
	"unsafe"

	sysmem "github.com/unsafe-risk/umem/internal/sysmem"
)

// This Implementation is based on the proposal in the following url: https://github.com/golang/go/issues/51317

// Thread-unsafe Allocation Arena.
type Arena struct {
	// The start address of the region.
	head uintptr
	// Tail of the region.
	tail uintptr
}

func NewFinalizer() *Arena {
	a := &Arena{}
	runtime.SetFinalizer(a, arenaFinalizer)
	return a
}

func New() *Arena {
	a := &Arena{}
	return a
}

func arenaFinalizer(a *Arena) {
	a.Free()
}

// Page Structure
/*
	|  0  |  1  |  2  |  3  |  4  |  5  |  6  |  7  |
	|-----|-----|-----|-----|-----|-----|-----|-----|
	|  page size            |  page head            |
	|-----|-----|-----|-----|-----|-----|-----|-----|
	|  Next Page Ptr                                |
	|-----|-----|-----|-----|-----|-----|-----|-----|
	|                                               |
	|                                               |
	|                                               |
	|                      Data                     |
	|                                               |
	|                                               |
	|                                               |
	|-----|-----|-----|-----|-----|-----|-----|-----|
*/

var defaultPageSize = uintptr(os.Getpagesize()*16 - 16)

func (r *Arena) newPage(size uintptr) {
	// println("Allocating new page", size)
	sptr := sysmem.SYSAllocOS(size + 16)
	pagesize := (*uint32)(unsafe.Pointer(sptr))
	pagehead := (*uint32)(unsafe.Pointer(uintptr(sptr) + 4))
	nextpage := (*uint64)(unsafe.Pointer(uintptr(sptr) + 8))

	*pagesize = uint32(size)
	*pagehead = 0
	*nextpage = 0

	if r.tail != 0 {
		// Add to the tail of the region.
		tailNextPage := (*uint64)(unsafe.Pointer(r.tail + 8))
		if *tailNextPage != 0 {
			*nextpage = *tailNextPage
		}
		*tailNextPage = uint64(uintptr(sptr))
	}
	r.tail = uintptr(sptr)
	if r.head == 0 {
		r.head = uintptr(sptr)
	}
	// println("New page allocated", size, sptr)
}

func (r *Arena) Reset() {
	r.tail = r.head
	for r.head != 0 {
		pagehead := (*uint32)(unsafe.Pointer(r.head + 4))
		nextpage := (*uint64)(unsafe.Pointer(r.head + 8))
		*pagehead = 0
		r.head = uintptr(*nextpage)
	}
	r.head = r.tail
}

func (r *Arena) Free() {
	for r.head != 0 {
		pagesize := (*uint32)(unsafe.Pointer(r.head))
		nextpage := (*uint64)(unsafe.Pointer(r.head + 8))
		nexthead := uintptr(*nextpage)
		sysmem.SYSFreeOS(unsafe.Pointer(r.head), uintptr(*pagesize+16))
		r.head = nexthead
	}
	r.tail = 0
}

const align = unsafe.Alignof(uintptr(0))

func (r *Arena) allocate(size uintptr) uintptr {
	size = ((size + align - 1) / align) * align
retry:
	if r.tail == 0 {
		// println("tail is 0, allocating new page")
		if size > defaultPageSize {
			r.newPage(size)
		} else {
			r.newPage(defaultPageSize)
		}
	}

	pagesize := (*uint32)(unsafe.Pointer(r.tail))
	pagehead := (*uint32)(unsafe.Pointer(r.tail + 4))
	nextpage := (*uint64)(unsafe.Pointer(r.tail + 8))
	if uintptr(*pagesize-*pagehead) < size {
		if *nextpage != 0 {
			r.tail = uintptr(*nextpage)
			goto retry
		}
		if size > defaultPageSize {
			r.newPage(size)
		} else {
			r.newPage(defaultPageSize)
		}
		pagesize = (*uint32)(unsafe.Pointer(r.tail))
		pagehead = (*uint32)(unsafe.Pointer(r.tail + 4))
		nextpage = (*uint64)(unsafe.Pointer(r.tail + 8))
	}

	data := r.tail + 16 + uintptr(*pagehead)
	*pagehead += uint32(size)
	return data
}

func (r *Arena) NewBytes(size uintptr) []byte {
	sh := reflect.SliceHeader{
		Data: r.allocate(size),
		Len:  int(size),
		Cap:  int(size),
	}
	return *(*[]byte)(unsafe.Pointer(&sh))
}

func (r *Arena) NewString(b []byte) string {
	s := r.NewBytes(uintptr(len(b)))
	copy(s, b)
	return *(*string)(unsafe.Pointer(&s))
}

func (r *Arena) HeapString(b string) string {
	s := r.NewBytes(uintptr(len(b)))
	copy(s, b)
	return *(*string)(unsafe.Pointer(&s))
}

func (r *Arena) Allocate(size uintptr) unsafe.Pointer {
	return unsafe.Pointer(r.allocate(size))
}

func NewOf[T any](r *Arena) *T {
	var zero T
	p := (*T)(r.Allocate(unsafe.Sizeof(zero)))
	*p = zero
	return p
}

func NewOfUninitialized[T any](r *Arena) *T {
	var zero T
	return (*T)(r.Allocate(unsafe.Sizeof(zero)))
}

func NewSliceOfUninitialized[T any](r *Arena, len int) []T {
	var zero T
	p := r.Allocate(unsafe.Sizeof(zero) * uintptr(len))
	sh := reflect.SliceHeader{
		Data: uintptr(p),
		Len:  len,
		Cap:  len,
	}
	v := *(*[]T)(unsafe.Pointer(&sh))
	return v
}

func NewSliceOf[T any](r *Arena, len int) []T {
	var zero T
	v := NewSliceOfUninitialized[T](r, len)
	for i := 0; i < len; i++ {
		v[i] = zero
	}
	return v
}
