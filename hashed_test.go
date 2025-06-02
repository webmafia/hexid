package hexid

import (
	"fmt"
	"testing"
)

func ExampleHashedID() {
	fmt.Println(HashedID("foobar"))
	fmt.Println(HashedID("foobaz"))

	// Output:
}

func BenchmarkHashedID(b *testing.B) {
	for b.Loop() {
		_ = HashedID("foobar")
	}
}
