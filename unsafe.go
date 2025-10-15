package hexid

import "unsafe"

//go:inline
func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

//go:inline
func b2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
