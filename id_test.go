package hexid

import (
	"testing"
	"time"
)

/*
TestIDEncodingDecoding verifies that encoding and decoding of IDs
round-trip correctly for the 63-bit layout:

[32 bits seconds][10 bits ms][6 bits node][15 bits seq]
*/

func TestIDEncodingDecoding(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp time.Time
		node      uint8
		seq       uint16
	}{
		{"normal", time.Unix(1708400000, 123_000_000), 1, 42},
		{"max-ms-seq", time.Unix(1600000000, 999_000_000), 7, 32767},
		{"hashed-node", time.Unix(1500000000, 0), 0, 0}, // Hashed → should return zero Time()
		{"max-node", time.Unix(1800000000, 500_000_000), 63, 12345},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := newID(tc.timestamp, tc.node, tc.seq)
			t.Logf("raw id: %064b", id)

			// Verify timestamp bits
			if got := id.Unix(); got != uint32(tc.timestamp.Unix()) {
				t.Errorf("Unix() mismatch: got %d, want %d", got, tc.timestamp.Unix())
			}

			// Verify milliseconds
			wantMS := uint16(tc.timestamp.Nanosecond() / 1_000_000)
			if got := id.Millis(); got != wantMS {
				t.Errorf("Millis() mismatch: got %d, want %d", got, wantMS)
			}

			// Verify node
			if got := id.Node(); got != tc.node {
				t.Errorf("Node() mismatch: got %d, want %d", got, tc.node)
			}

			// Verify sequence
			if got := id.Seq(); got != tc.seq&0x7FFF {
				t.Errorf("Seq() mismatch: got %d, want %d", got, tc.seq&0x7FFF)
			}

			// Verify reconstructed time
			reconstructed := id.Time()
			if tc.node == 0 {
				// Hashed ID → should produce zero time
				if !reconstructed.IsZero() {
					t.Errorf("Hashed ID should yield zero time, got %v", reconstructed)
				}
			} else {
				// Normal ID → should reconstruct correctly
				if reconstructed.Unix() != tc.timestamp.Unix() {
					t.Errorf("Time() seconds mismatch: got %d, want %d",
						reconstructed.Unix(), tc.timestamp.Unix())
				}
				if ms := reconstructed.Nanosecond() / 1_000_000; uint16(ms) != wantMS {
					t.Errorf("Time() milliseconds mismatch: got %d, want %d", ms, wantMS)
				}
			}
		})
	}
}

func TestSequenceOverflow(t *testing.T) {
	ts := time.Unix(1700000000, 500_000_000) // arbitrary
	node := uint8(5)

	// Force sequence to just before overflow
	id1 := newID(ts, node, 0x7FFF) // 32767
	id2 := newID(ts, node, 0x8000) // 32768 (should wrap)

	if id1.Seq() != 0x7FFF {
		t.Fatalf("expected 0x7FFF before overflow, got %d", id1.Seq())
	}
	if id2.Seq() != 0 {
		t.Fatalf("expected sequence wrap to 0, got %d", id2.Seq())
	}

	// Ensure timestamp and node bits are identical
	if id1.Unix() != id2.Unix() {
		t.Fatalf("seconds changed on overflow: %d vs %d", id1.Unix(), id2.Unix())
	}
	if id1.Millis() != id2.Millis() {
		t.Fatalf("millis changed on overflow: %d vs %d", id1.Millis(), id2.Millis())
	}
	if id1.Node() != id2.Node() {
		t.Fatalf("node changed on overflow: %d vs %d", id1.Node(), id2.Node())
	}
}
