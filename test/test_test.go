package test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
)

const (
	multiplier    uint64 = 0x6eed0e9da4d94a4f
	invMultiplier uint64 = 0x2f72b4215a3d8caf
)

// Test that the constants are true modular inverses mod 2^64.
func TestMultiplierInverse(t *testing.T) {
	// force runtime evaluation (avoid compile-time overflow)
	m := multiplier
	inv := invMultiplier
	product := m * inv // wraps mod 2^64 automatically

	if product != 1 {
		t.Fatalf("multiplier*invMultiplier ≡ 0x%016x (want 1)", product)
	}

	if m%2 == 0 {
		t.Fatalf("multiplier must be odd to be invertible mod 2^64")
	}

	// round-trip check for sample values
	for _, v := range []uint64{0, 1, 2, 1234567890, 0xffffffffffffffff} {
		scr := v * m
		unscr := scr * inv
		if unscr != v {
			t.Fatalf("round-trip failed for %d: got %d", v, unscr)
		}
	}

	t.Log("✓ multiplier and inverse verified modulo 2^64")
}

func modInverse64(a uint64) uint64 {
	var x, y uint64 = a, 1
	for i := 0; i < 6; i++ {
		y *= 2 - x*y
	}
	return y
}

func Example_reverse() {
	const multiplier = 0x6eed0e9da4d94a4f
	inv := modInverse64(multiplier)
	fmt.Printf("multiplier:  0x%016x\n", multiplier)
	fmt.Printf("inverse:     0x%016x\n", inv)
	fmt.Printf("check:       0x%016x\n", multiplier*inv)
	// Output:
}

func scrambleToHex(id uint64) string {
	scrambled := id * multiplier
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], scrambled)
	return hex.EncodeToString(buf[:])
}

func unscrambleFromHex(hexStr string) uint64 {
	b, _ := hex.DecodeString(hexStr)
	v := binary.BigEndian.Uint64(b)
	original := v * invMultiplier
	return original
}

func Example() {
	ints := []uint64{1234567890, 9223372036854775807}

	for _, i := range ints {
		h := scrambleToHex(i)
		fmt.Println(i, h, unscrambleFromHex(h))
	}

	// Output:
}
