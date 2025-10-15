package hexid

import (
	"testing"
	"time"
)

// TestIDEncodingDecoding verifies that encoding an ID with newID()
// and then extracting its components gives back the original values.
//
// Layout under test (63 bits total):
// [32 bits seconds][10 bits ms][6 bits node][15 bits seq]
func TestIDEncodingDecoding(t *testing.T) {
	testCases := []struct {
		timestamp time.Time
		node      uint8
		seq       uint16
	}{
		{time.Unix(1708400000, 123_000_000), 1, 42},     // normal case
		{time.Unix(1600000000, 999_000_000), 7, 32767},  // max ms + max seq
		{time.Unix(1500000000, 0), 0, 0},                // all zero
		{time.Unix(1800000000, 500_000_000), 63, 12345}, // max node
	}

	for _, tc := range testCases {
		id := newID(tc.timestamp, tc.node, tc.seq)

		// Check timestamp seconds
		if got := id.Unix(); got != uint32(tc.timestamp.Unix()) {
			t.Errorf("Unix() mismatch: got %d, want %d", got, tc.timestamp.Unix())
		}

		// Check milliseconds
		wantMS := uint16(tc.timestamp.Nanosecond() / 1_000_000)
		if got := id.Millis(); got != wantMS {
			t.Errorf("Millis() mismatch: got %d, want %d", got, wantMS)
		}

		// Check node
		if got := id.Node(); got != tc.node {
			t.Errorf("Node() mismatch: got %d, want %d", got, tc.node)
		}

		// Check sequence
		if got := id.Seq(); got != tc.seq&0x7FFF { // 15 bits
			t.Errorf("Seq() mismatch: got %d, want %d", got, tc.seq&0x7FFF)
		}

		// Reconstruct time
		reconstructed := id.Time()
		if reconstructed.Unix() != tc.timestamp.Unix() {
			t.Errorf("Time() seconds mismatch: got %d, want %d", reconstructed.Unix(), tc.timestamp.Unix())
		}
		if ms := reconstructed.Nanosecond() / 1_000_000; uint16(ms) != wantMS {
			t.Errorf("Time() milliseconds mismatch: got %d, want %d", ms, wantMS)
		}
	}
}
