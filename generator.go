package hexid

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

// Non-thread-safe ID generator. Can generate up to 2^15 (32,768) locally unique IDs per millisecond per node.
type Generator struct {
	seq  uint32
	node uint8
}

// Create an ID generator. The generator is NOT thread-safe.
func NewGenerator(node ...uint8) (g Generator, err error) {
	var n uint8 = 1

	if len(node) > 0 {
		n = node[0]

		if n < 1 || n > 63 {
			err = fmt.Errorf("node must be between 1 and 63")
			return
		}
	}

	return Generator{
		seq:  rand.Uint32(),
		node: n,
	}, nil
}

func (g *Generator) ID() (id ID) {
	id = newID(time.Now(), g.node, uint16(g.seq))
	g.seq++
	return
}

func (g *Generator) IDFromTime(ts time.Time) (id ID) {
	id = newID(ts, g.node, uint16(g.seq))
	g.seq++
	return
}

// Thread-safe ID generator. Can generate up to 2^15 (32,768) locally unique IDs per millisecond per node.
type AtomicGenerator struct {
	seq  uint32
	node uint8
}

// Create an atomic ID generator. The generator is thread-safe.
func NewAtomicGenerator(node ...uint8) (g AtomicGenerator, err error) {
	var n uint8 = 1

	if len(node) > 0 {
		n = node[0]

		if n < 1 || n > 63 {
			err = fmt.Errorf("node must be between 1 and 63")
			return
		}
	}

	return AtomicGenerator{
		seq:  rand.Uint32(),
		node: n,
	}, nil
}

func (g *AtomicGenerator) ID() ID {
	return newID(time.Now(), g.node, uint16(atomic.AddUint32(&g.seq, 1)))
}

func (g *AtomicGenerator) IDFromTime(ts time.Time) ID {
	return newID(ts, g.node, uint16(atomic.AddUint32(&g.seq, 1)))
}
