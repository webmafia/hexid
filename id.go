package hexid

import "time"

/*
ID Layout (63 bits total, safe for Postgres BIGINT)

┌──────────────────────────────────────┬────────────┬──────────┬────────────┐
│ 32 bits unix seconds                │ 10 bits ms │ 6 bits node │ 15 bits seq │
└──────────────────────────────────────┴────────────┴──────────┴────────────┘
  [62...........................31]     [30......21] [20....15] [14........0]

- 32 bits seconds: valid until year 2106
- 10 bits milliseconds: 0–999 (precision within each second)
- 6 bits node: up to 64 generator nodes
- 15 bits sequence: 0–32767 IDs per millisecond per node
- Total = 63 bits (top bit always 0)

This structure guarantees:
- Chronological ordering
- Collision-free ID generation up to 40 ns/ID
- Safe storage in Postgres BIGINT
*/

type ID uint64

// newID generates a new 63-bit ID based on the given timestamp, node ID, and sequence counter.
func newID(ts time.Time, nodeID uint8, seq uint16) ID {
	const (
		msBits   = 10
		nodeBits = 6
		seqBits  = 15

		nodeShift = seqBits
		msShift   = nodeShift + nodeBits
		secShift  = msShift + msBits

		mask63 = 0x7FFFFFFFFFFFFFFF // ensure top bit = 0
	)

	secs := uint64(ts.Unix())
	msecs := uint64(ts.Nanosecond() / 1_000_000)

	id := (secs << secShift) |
		(msecs << msShift) |
		(uint64(nodeID) << nodeShift) |
		(uint64(seq) & ((1 << seqBits) - 1))

	return ID(id & mask63)
}

// Unix returns the Unix timestamp in seconds.
func (id ID) Unix() uint32 {
	return uint32(id >> 31) // shift away ms+node+seq bits
}

// Millis returns the millisecond part within the second.
func (id ID) Millis() uint16 {
	return uint16((id >> 21) & 0x3FF) // 10 bits
}

// Node returns the 6-bit node ID.
func (id ID) Node() uint8 {
	return uint8((id >> 15) & 0x3F)
}

// Seq returns the 15-bit sequence number.
func (id ID) Seq() uint16 {
	return uint16(id & 0x7FFF)
}

// Entropy returns everything after the Unix timestamp seconds (milliseconds + node + sequence)
func (id ID) Entropy() uint32 {
	// Extract the lower 31 bits: milliseconds (10) + node (6) + sequence (15)
	return uint32(id & 0x7FFFFFFF)
}

// Time reconstructs the approximate creation time of the ID.
func (id ID) Time() time.Time {
	if id.Hashed() {
		return time.Time{}
	}

	secs := int64(id >> 31)
	ms := int64((id >> 21) & 0x3FF)
	return time.Unix(secs, ms*1_000_000)
}

// Uint64 returns the raw numeric value of the ID.
func (id ID) Uint64() uint64 {
	return uint64(id)
}

// Int64 returns the raw numeric value of the ID.
func (id ID) Int64() int64 {
	return int64(id)
}

func (id ID) Hashed() bool {
	return id.Node() == 0
}

// IsZero reports whether the ID is zero.
func (id ID) IsZero() bool { return id == 0 }

// IsNil is equivalent to IsZero.
func (id ID) IsNil() bool { return id == 0 }
