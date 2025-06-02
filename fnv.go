package hexid

import (
	"encoding/binary"
	"hash"
)

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

var _ hash.Hash64 = (*fnv64a)(nil)

// An allocation-free implementation of FNV-1a.
type fnv64a uint64

func newFnv64a() fnv64a          { return fnv64a(offset64) }
func (s *fnv64a) BlockSize() int { return 1 }
func (s *fnv64a) Reset()         { *s = offset64 }
func (s *fnv64a) Size() int      { return 8 }
func (s *fnv64a) Sum64() uint64  { return uint64(*s) }

func (s *fnv64a) Sum(b []byte) []byte {
	return binary.BigEndian.AppendUint64(b, s.Sum64())
}

func (s *fnv64a) Write(data []byte) (int, error) {
	hash := *s
	for _, c := range data {
		hash ^= fnv64a(c)
		hash *= prime64
	}
	*s = hash
	return len(data), nil
}
