package arena

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestUmemAtomicTest(t *testing.T) {
	r := New()
	p := NewOfUninitialized[Person](r)
	p.Name = "John"
	p.Age = 42
	p.Address = "London"
	p.number = 123
	p.uuid = "12345"

	_ = NewOfUninitialized[byte](r)
	u64p := NewOfUninitialized[uint64](r)
	*u64p = 1234567890
	fmt.Printf("u64p: %p\n", u64p)
	atomic.AddUint64(u64p, 1)
	if *u64p != 1234567891 {
		t.Errorf("Expected 1234567891, got %d", *u64p)
		t.FailNow()
	}

	if !atomic.CompareAndSwapUint64(u64p, 1234567891, 1234567892) {
		t.Errorf("Expected 1234567891 to be swapped to 1234567892, but it wasn't")
		t.FailNow()
	}
	r.Free()
}

type Person struct {
	Name    string
	Age     int
	Address string
	number  int
	uuid    string
}

const nAlloc = 1000000

func BenchmarkAllocateUmemUninitializedPerson(b *testing.B) {
	r := New()
	for i := 0; i < b.N; i++ {
		for j := 0; j < nAlloc; j++ {
			p := NewOfUninitialized[Person](r)
			p.Name = "John"
			p.Age = 42
			p.Address = "London"
			p.number = i
			p.uuid = "12345"
		}
		r.Free()
	}
}

func BenchmarkAllocateUmemPerson(b *testing.B) {
	r := New()
	for i := 0; i < b.N; i++ {
		for j := 0; j < nAlloc; j++ {
			p := NewOf[Person](r)
			p.Name = "John"
			p.Age = 42
			p.Address = "London"
			p.number = i
			p.uuid = "12345"
		}
		r.Free()
	}
}

//go:noinline
func StdNewPerson() *Person {
	p := new(Person)
	return p
}

func BenchmarkAllocateStdNew(b *testing.B) {
	var ps []*Person = make([]*Person, nAlloc)
	for i := 0; i < b.N; i++ {
		for j := 0; j < nAlloc; j++ {
			p := StdNewPerson()
			p.Name = "John"
			p.Age = 42
			p.Address = "London"
			p.number = i
			p.uuid = "12345"
			ps[j] = p
		}
	}
}
