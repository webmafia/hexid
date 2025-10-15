package hexid

// nodeMask clears bits 20–15 (the 6-bit node field) while keeping all other bits intact.
const nodeMask uint64 = ^(uint64(0x3F) << 15) & 0x7FFFFFFFFFFFFFFF

// HashedID produces a deterministic 63-bit ID from one or more strings.
// The resulting ID is based on an FNV-1a hash and always has node ID = 0.
// The timestamp portion is meaningless but guaranteed not to collide
// with time-based IDs, since those always have node ≥ 1.
func HashedID(s ...string) ID {
	h := newFnv64a()

	for _, str := range s {
		h.Write(s2b(str))
	}

	id := h.Sum64()
	return ID(id & nodeMask)
}

// HashedIDBytes produces a deterministic 63-bit ID from a byte slice.
// The resulting ID always has node ID = 0.
func HashedIDBytes(b []byte) ID {
	h := newFnv64a()
	h.Write(b)
	id := h.Sum64()
	return ID(id & nodeMask)
}
