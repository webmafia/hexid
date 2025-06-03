package hexid

import (
	"fmt"
	"testing"
)

func ExampleHashedID() {
	fmt.Println(HashedID("foobar"))
	fmt.Println(HashedID("foobaz"))

	// Output:
	//
	// 06e441e53719e608
	// f0e4e84ada932890
}

func BenchmarkHashedID(b *testing.B) {
	for b.Loop() {
		_ = HashedID("foobar")
	}
}
