package hexid

import (
	"fmt"
	"testing"
)

func BenchmarkID_AppendText(b *testing.B) {
	var buf []byte
	g, _ := NewGenerator()
	id := g.ID()
	b.ResetTimer()

	for range b.N {
		buf, _ = id.AppendText(buf[:0])
	}
}

func BenchmarkIDFromString(b *testing.B) {
	id := Generate().String()
	b.ResetTimer()

	for range b.N {
		_, _ = IDFromString(id)
	}
}

func TestMultiplierInverse(t *testing.T) {
	var multiplier = multiplier
	var invMultiplier = invMultiplier

	// Multiplication on uint64 is performed modulo 2^64.
	product := multiplier * invMultiplier

	if product != 1 {
		t.Errorf("Expected multiplier * invMultiplier mod 2^64 to equal 1, got %x", product)
	}
}

func Example_coder() {
	id := ID(123)
	fmt.Println(id)

	// Output: 4be605be3466b3f5
}

func ExampleIDFromString() {
	_, err := IDFromString("3784432400289806371")

	if err != nil {
		fmt.Println(err)
	}

	// Output: invalid ID
}
