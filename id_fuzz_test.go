package hexid

import (
	"testing"
	"time"
)

func FuzzID(f *testing.F) {
	// Seed with some initial values.
	f.Add(int64(1600000000), int64(123456789), uint8(1), uint16(42))
	f.Add(int64(0), int64(0), uint8(0), uint16(0))
	f.Add(time.Now().Unix(), int64(987654321), uint8(63), uint16(32767)) // max seq value for 15 bits

	f.Fuzz(func(t *testing.T, sec int64, nsec int64, node uint8, seq uint16) {
		if sec < 0 || nsec < 0 {
			t.Skip()
		}

		// Clamp nanoseconds to [0, 1e9)
		nsec = nsec % 1_000_000_000
		node = node & 0x3F // 6 bits (0–63)
		seq = seq & 0x7FFF // 15 bits (0–32767)

		// Create a timestamp and generate an ID.
		ts := time.Unix(sec, nsec)
		id := newID(ts, node, seq)

		// Verify that the Unix seconds are preserved.
		if got, want := id.Unix(), uint32(ts.Unix()); got != want {
			t.Errorf("Unix mismatch: got %d, want %d", got, want)
		}

		// Verify that the Time method returns a timestamp truncated to milliseconds.
		ms := ts.Nanosecond() / 1_000_000
		expectedTime := time.Unix(ts.Unix(), int64(ms)*1_000_000)
		if got := id.Time(); !got.Equal(expectedTime) {
			t.Errorf("Time mismatch: got %v, want %v", got, expectedTime)
		}

		// Verify that node ID is preserved.
		if got, want := id.Node(), node; got != want {
			t.Errorf("Node mismatch: got %d, want %d", got, want)
		}

		// Verify that the sequence number is limited to 15 bits.
		if got, want := id.Seq(), seq&0x7FFF; got != want {
			t.Errorf("Seq mismatch: got %d, want %d", got, want)
		}

		// Verify that AppendText and DecodeText form a round-trip.
		decoded, err := IDFromString(id.String())
		if err != nil {
			t.Fatalf("DecodeText error: %v", err)
		}
		if decoded != id {
			t.Errorf("Decoded ID mismatch: got %v, want %v", decoded, id)
		}
	})
}
