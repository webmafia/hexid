package hexid

import (
	"github.com/webmafia/fast"
)

// Produces a deterministic ID from string(s). The timestamp part of the resulting ID
// will be be pointless, but guaranteed to be before 2004-01-10 13:37:04 UTC.
func HashedID(s ...string) ID {
	// Hash the input with FNV-1a.
	h := newFnv64a()

	for _, s := range s {
		h.Write(fast.StringToBytes(s))
	}

	id := h.Sum64()

	// Clamp the timestamp (32 first bits) to 2004-01-10 13:37:04 UTC.
	// We do this to ensure that it's impossible to get a collision with an ID of the future.
	return ID(id & 0x3FFFFFFFFFFFFFFF)
}

// Produces a deterministic ID from a byte slice. The timestamp part of the resulting ID
// will be be pointless, but guaranteed to be before 2004-01-10 13:37:04 UTC.
func HashedIDBytes(b []byte) ID {

	// Hash the input with FNV-1a.
	h := newFnv64a()
	h.Write(b)
	id := h.Sum64()

	// Clamp the timestamp (32 first bits) to 2004-01-10 13:37:04 UTC.
	// We do this to ensure that it's impossible to get a collision with an ID of the future.
	return ID(id & 0x3FFFFFFFFFFFFFFF)
}
