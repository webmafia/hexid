package hexid

import (
	"hash/fnv"
	"testing"
)

func BenchmarkStandardFnv(b *testing.B) {
	b.Run("New", func(b *testing.B) {
		for b.Loop() {
			_ = fnv.New64a()
		}
	})

	b.Run("Write", func(b *testing.B) {
		h := fnv.New64a()

		for b.Loop() {
			_, _ = h.Write([]byte{1})
		}
	})

	b.Run("Sum", func(b *testing.B) {
		h := fnv.New64a()
		var buf []byte

		for b.Loop() {
			buf = h.Sum(buf[:0])
		}
	})

	b.Run("Sum64", func(b *testing.B) {
		h := fnv.New64a()

		for b.Loop() {
			_ = h.Sum64()
		}
	})
}

func BenchmarkCustomFnv(b *testing.B) {
	b.Run("New", func(b *testing.B) {
		for b.Loop() {
			_ = newFnv64a()
		}
	})

	b.Run("Write", func(b *testing.B) {
		h := newFnv64a()

		for b.Loop() {
			_, _ = h.Write([]byte{1})
		}
	})

	b.Run("Sum", func(b *testing.B) {
		h := newFnv64a()
		var buf []byte

		for b.Loop() {
			buf = h.Sum(buf[:0])
		}
	})

	b.Run("Sum64", func(b *testing.B) {
		h := newFnv64a()

		for b.Loop() {
			_ = h.Sum64()
		}
	})
}
