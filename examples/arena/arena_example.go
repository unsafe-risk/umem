package main

import (
	"fmt"

	"github.com/unsafe-risk/umem/arena"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	a := arena.New()

	p := arena.NewOf[Person](a)
	p.Name = "John"
	p.Age = 30

	fmt.Println(p)

	a.Free()
	// p is invalid after Free()

	// fmt.Println(p)
	// panic: unexpected fault address
	// fatal error: fault
}
