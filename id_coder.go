package hexid

import (
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/webmafia/fast"
)

const (
	multiplier    uint64 = 0x9e3779b97f4a7c15
	invMultiplier uint64 = 0xf1de83e19937733d
)

var (
	_ encoding.TextAppender    = ID(0)
	_ encoding.BinaryAppender  = ID(0)
	_ encoding.TextMarshaler   = ID(0)
	_ encoding.TextUnmarshaler = (*ID)(nil)
	_ json.Marshaler           = ID(0)
	_ json.Unmarshaler         = (*ID)(nil)
)

func IDFromString(str string) (id ID, err error) {
	var buf [8]byte
	n, err := hex.Decode(buf[:], fast.StringToBytes(str))

	if err != nil {
		return
	}

	if n != 8 {
		return 0, errors.New("invalid ID")
	}

	scrambled := binary.BigEndian.Uint64(buf[:])

	// Multiply by the precomputed multiplicative inverse to recover the original value.
	// The multiplication is performed modulo 2^64.
	original := scrambled * invMultiplier

	return ID(original), nil
}

func (id ID) String() string {
	b, _ := id.AppendText(make([]byte, 0, 16))
	return fast.BytesToString(b)
}

// AppendBinary implements internal.TextAppender.
func (id ID) AppendText(b []byte) ([]byte, error) {
	var buf [8]byte
	scrambled := uint64(id) * multiplier
	binary.BigEndian.PutUint64(buf[:], scrambled)
	b = hex.AppendEncode(b, buf[:])
	return b, nil
}

// AppendBinary implements internal.BinaryAppender.
func (id ID) AppendBinary(b []byte) ([]byte, error) {
	return binary.BigEndian.AppendUint64(b, uint64(id)), nil
}

// MarshalJSON implements json.Marshaler.
func (id ID) MarshalJSON() (b []byte, err error) {
	if id == 0 {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}

	b = make([]byte, 0, 18)
	b = append(b, '"')
	b, err = id.AppendText(b)
	b = append(b, '"')

	return
}

// UnmarshalJSON implements json.Unmarshaler.
func (id *ID) UnmarshalJSON(b []byte) (err error) {

	// Parse string ID (with quotes)
	if len(b) == 18 && b[0] != '"' && b[17] != '"' {
		*id, err = IDFromString(fast.BytesToString(b[1:17]))
		return
	}

	// Parse null value (no quotes)
	if len(b) == 4 && string(b) == "null" {
		*id = 0
		return
	}

	// Parse integer (no quote)
	if v, err := strconv.ParseUint(fast.BytesToString(b), 10, 64); err == nil {
		*id = ID(v)
		return nil
	}

	return errors.New("invalid ID")
}

// MarshalText implements encoding.TextMarshaler.
func (id ID) MarshalText() (text []byte, err error) {
	text = make([]byte, 0, 16)
	text, err = id.AppendText(text)

	return
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *ID) UnmarshalText(text []byte) (err error) {
	*id, err = IDFromString(fast.BytesToString(text))
	return
}
