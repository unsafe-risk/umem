package arena

import (
	"testing"
)

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
	}
	r.Free()
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
	}
	r.Free()
}

//go:noinline
func StdNewPerson() *Person {
	p := new(Person)
	return p
}

func BenchmarkAllocateStdNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < nAlloc; j++ {
			p := StdNewPerson()
			p.Name = "John"
			p.Age = 42
			p.Address = "London"
			p.number = i
			p.uuid = "12345"
		}
	}
}
