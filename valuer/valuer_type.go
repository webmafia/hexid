package valuer

type Type uint32

const (
	Int64Valuer  Type = iota // Encodes as int64 (default)
	Uint64Valuer             // Encodes as uint64
	StringValuer             // Encodes as HEX encoded string (16 bytes)
	BinaryValuer             // Encodes as 8 raw bytes
)
