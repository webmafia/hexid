package hexid

import (
	"fmt"
	"testing"
	"time"
)

func ExampleGenerator() {
	g, err := NewGenerator()

	if err != nil {
		panic(err)
	}

	// Ensure a deterministic sequence in this example
	g.seq = 1

	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	for range 4 {
		id := g.IDFromTime(ts)
		fmt.Println(id)
	}

	// Output:
	//
	// 58c1fa4a4a00ca4f
	// c7af08e7eeda149e
	// 369c178593b35eed
	// a5892623388ca93c
}

func BenchmarkGenerator(b *testing.B) {
	b.Run("New", func(b *testing.B) {
		for range b.N {
			_, _ = NewGenerator()
		}
	})

	b.Run("ID", func(b *testing.B) {
		g, _ := NewGenerator()
		b.ResetTimer()

		for range b.N {
			_ = g.ID()
		}
	})

	b.Run("IDFromTime", func(b *testing.B) {
		g, _ := NewGenerator()
		ts := time.Now()
		b.ResetTimer()

		for range b.N {
			_ = g.IDFromTime(ts)
		}
	})
}

func BenchmarkAtomicGenerator(b *testing.B) {
	b.Run("New", func(b *testing.B) {
		for range b.N {
			_, _ = NewAtomicGenerator()
		}
	})

	b.Run("ID", func(b *testing.B) {
		g, _ := NewAtomicGenerator()
		b.ResetTimer()

		for range b.N {
			_ = g.ID()
		}
	})

	b.Run("IDFromTime", func(b *testing.B) {
		g, _ := NewAtomicGenerator()
		ts := time.Now()
		b.ResetTimer()

		for range b.N {
			_ = g.IDFromTime(ts)
		}
	})
}

func TestID_uniqueness(t *testing.T) {
	g := Generator{seq: 0}
	ts := time.Now()

	id1 := g.IDFromTime(ts)
	id2 := g.IDFromTime(ts)

	if id1 == id2 {
		t.Fatalf("identical IDs '%s' (%d)", id1, id2)
	}
}

func TestGeneratorWrapAroundDoesNotDuplicate(t *testing.T) {
	// Use a fixed time to ensure that the time portion of the ID remains constant.
	fixedTime := time.Unix(1600000000, 123456000)

	// Create a generator and force the sequence to the maximum effective value.
	g, _ := NewGenerator()
	// Set the sequence such that the effective sequence (seq & 0x3FFFFF)
	// is (1<<22)-1, i.e. 4194303.
	g.seq = (1 << 22) - 1

	// Generate an ID with the maximum effective sequence.
	id1 := g.IDFromTime(fixedTime)
	// The next call will use g.seq = 1<<22 which yields an effective sequence of 0.
	id2 := g.IDFromTime(fixedTime)

	if id1 == id2 {
		t.Fatalf("IDs are identical despite wrap-around: id1=%v, id2=%v", id1, id2)
	}
}
