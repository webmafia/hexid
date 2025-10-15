package hexid

import (
	"testing"
)

func FuzzHashedIDNoCollisions(f *testing.F) {
	f.Fuzz(func(t *testing.T, a string, b string) {
		if a == b {
			return // Skip identical inputs
		}
		ida := HashedID(a)
		idb := HashedID(b)

		if ida == idb {
			t.Errorf("Collision detected: HashedID(%q) == HashedID(%q) == %d", a, b, ida)
		}
	})
}
