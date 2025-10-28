package hexid

import (
	"fmt"
	"testing"
	"time"
)

func ExampleIDFromEntropy() {
	id1, _ := IDFromString("24a3b3372e1a0a50")
	id2 := IDFromEntropy(id1.Unix(), id1.Entropy())

	fmt.Println(id1)
	fmt.Println(id2)

	// Output:
	// 24a3b3372e1a0a50
	// 24a3b3372e1a0a50
}

func TestIDFromEntropy(t *testing.T) {
	ts := time.Unix(1730000000, 123_000_000) // 123 ms
	id := IDFromTime(ts)

	unix := id.Unix()
	entropy := id.Entropy()

	// Reassemble
	id2 := IDFromEntropy(unix, entropy)

	if id2 != id {
		t.Fatalf("ID mismatch:\n got  %064b\n want %064b", id2, id)
	}

	// Double-check correctness of the reassembled parts
	if id2.Unix() != unix {
		t.Fatalf("unix mismatch: got %d, want %d", id2.Unix(), unix)
	}
	if id2.Entropy() != entropy {
		t.Fatalf("entropy mismatch: got %d, want %d", id2.Entropy(), entropy)
	}
}
