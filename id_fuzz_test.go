package hexid

import (
	"testing"
	"time"
)

func FuzzID(f *testing.F) {
	f.Add(int64(1600000000), int64(123456789), uint8(1), uint16(42))
	f.Add(int64(0), int64(0), uint8(0), uint16(0))
	f.Add(time.Now().Unix(), int64(987654321), uint8(63), uint16(32767))

	f.Fuzz(func(t *testing.T, sec int64, nsec int64, node uint8, seq uint16) {
		if sec < 0 || nsec < 0 {
			t.Skip()
		}

		nsec = nsec % 1_000_000_000
		node = node & 0x3F // 6 bits (0–63)
		seq = seq & 0x7FFF // 15 bits (0–32767)

		ts := time.Unix(sec, nsec)
		id := newID(ts, node, seq)

		// Verify basic fields
		if got, want := id.Unix(), uint32(ts.Unix()); got != want {
			t.Errorf("Unix mismatch: got %d, want %d", got, want)
		}
		if got, want := id.Node(), node; got != want {
			t.Errorf("Node mismatch: got %d, want %d", got, want)
		}
		if got, want := id.Seq(), seq&0x7FFF; got != want {
			t.Errorf("Seq mismatch: got %d, want %d", got, want)
		}

		// Skip time reconstruction for hashed IDs (node == 0)
		if node == 0 {
			return
		}

		// Verify that Time() reconstructs milliseconds accurately
		ms := ts.Nanosecond() / 1_000_000
		expectedTime := time.Unix(ts.Unix(), int64(ms)*1_000_000)
		if got := id.Time(); !got.Equal(expectedTime) {
			t.Errorf("Time mismatch: got %v, want %v", got, expectedTime)
		}

		// Verify round-trip encoding/decoding
		decoded, err := IDFromString(id.String())
		if err != nil {
			t.Fatalf("DecodeText error: %v", err)
		}
		if decoded != id {
			t.Errorf("Decoded ID mismatch: got %v, want %v", decoded, id)
		}
	})
}
