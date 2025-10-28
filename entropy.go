package hexid

// Reassambles an ID from `(ID).Unix()` and `(ID).Entropy()`.
func IDFromEntropy(unix, entropy uint32) ID {
	// Layout:
	// [63] unused
	// [62..31] = 32-bit unix seconds
	// [30..0]  = 31-bit entropy (ms + node + seq)
	return ID((uint64(unix) << 31) | uint64(entropy&0x7FFFFFFF))
}
