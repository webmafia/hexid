package hexid

import (
	"fmt"
	"testing"
)

func ExampleHashedID() {
	a := HashedID("foobar")
	b := HashedID("foobaz")

	fmt.Println(a.Hashed(), a)
	fmt.Println(b.Hashed(), b)

	// Output:
	//
	// true 45ecc9eb54b12098
	// true 951ba2f26ae6feb0
}

func BenchmarkHashedID(b *testing.B) {
	for b.Loop() {
		_ = HashedID("foobar")
	}
}
