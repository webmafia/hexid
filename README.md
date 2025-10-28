# hexid

**Zero-dependency, zero-allocation, time-sortable, random-looking 63-bit IDs for Go. ~40ns per ID.**

`hexid` is a compact, deterministic ID system that produces 63-bit identifiers safe for PostgreSQL `BIGINT`.  
IDs are chronologically sortable, distributed-safe, and allocation-free (except when converting to hex strings).  
All encoding and decoding are fully deterministic between Go and PostgreSQL.

---

## 🚀 Key Features

- **🪶 Zero dependencies:** Pure Go — no third-party packages.  
- **⚡ Zero allocations:** Except when encoding to a new hex string.  
- **🐘 Compact & efficient:** 63-bit IDs fit safely in Postgres `BIGINT`.  
- **⏱️ Time-sortable:** Encodes seconds + milliseconds for chronological order.  
- **🌍 Distribution-safe:** 6-bit node field (up to 63 nodes).  
- **💥 High throughput:** ~25 million IDs/s per node (~40 ns per ID).  
- **🧠 Deterministic:** Identical encoding and decoding in Go and PostgreSQL.  
- **🔒 Hash mode:** Deterministic `HashedID()` for stable, non-time-based IDs.  

---

## 📦 Installation

```bash
go get github.com/webmafia/hexid
```

Import and use:

```go
import "github.com/webmafia/hexid"
```

---

## 🧩 ID Layout

```
 63 62                            31 30      21 20  15 14            0
┌──┬────────────────────────────────┬──────────┬──────┬───────────────┐
│X │ 32 bits unix seconds           │ 10 bits  │ 6 b  │ 15 bits       │
│  │                                │ ms       │ node │ sequence      │
└──┴────────────────────────────────┴──────────┴──────┴───────────────┘
X = unused (sign bit of int64)
````

| Field        | Bits        | Range             | Purpose                                                                                           |
| ------------ | ----------- | ----------------- | ------------------------------------------------------------------------------------------------- |
| Seconds      | 32          | 0 – 4,294,967,295 | Valid until year 2106                                                                             |
| Milliseconds | 10          | 0 – 999           | Sub-second precision                                                                              |
| Node         | 6           | 1 – 63            | Up to 63 generator nodes (`0` is reserved for [hashed IDs](#4-deterministic-non-time-hashed-ids)) |
| Sequence     | 15          | 0 – 32 767        | Per-ms per-node counter                                                                           |
| **Total**    | **63 bits** | < 2⁶³             | Safe in signed `BIGINT`                                                                           |

---

## 🧰 Usage

### 1. Global generator (thread-safe)

```go
id := hexid.Generate()
fmt.Println(id.String()) // scrambled 16-char hex string

id2 := hexid.IDFromTime(time.Now())
````

### 2. Local generator (faster, not thread-safe)

```go
g, _ := hexid.NewGenerator(5) // node ID 5
id := g.ID()
```

### 3. Thread-safe generator

```go
g, _ := hexid.NewAtomicGenerator(12)
id := g.ID()
```

### 4. Deterministic (non-time) hashed IDs

```go
h1 := hexid.HashedID("user", "42")
h2 := hexid.HashedIDBytes([]byte("my-unique-key"))
```

Hashed IDs always have `Node() == 0` and a zero timestamp.

---

## 🧩 ID Accessors

| Method              | Description                             |
| ------------------- | --------------------------------------- |
| `id.Unix()`         | Extract unix seconds.                   |
| `id.Millis()`       | Milliseconds within the second (0–999). |
| `id.Node()`         | Node ID (0–63).                         |
| `id.Seq()`          | Sequence number (0–32 767).             |
| `id.Time()`         | Reconstruct creation time.              |
| `id.String()`       | Scrambled 16-character hex encoding.    |
| `IDFromString(str)` | Decode from hex string.                 |
| `id.Bytes()`        | 8-byte big-endian binary form.          |

---

## 🐘 Encoding/decoding from PostgreSQL

Matching SQL functions for direct database use:

```sql
CREATE OR REPLACE FUNCTION hexid_encode(id bigint)
RETURNS text AS $$
  SELECT lpad(
    to_hex(
      ((id::numeric * 7993060983890856527)
       % 9223372036854775808)::bigint
    ), 16, '0');
$$ LANGUAGE sql IMMUTABLE STRICT;

CREATE OR REPLACE FUNCTION hexid_decode(hexid text)
RETURNS bigint AS $$
  SELECT (
    (('x' || hexid)::bit(64)::bigint::numeric *
     3418993122468531375) % 9223372036854775808
  )::bigint;
$$ LANGUAGE sql IMMUTABLE STRICT;
```

These produce and decode exactly the same hex values as Go’s `String()` / `IDFromString()`.

---

## 🧬 Collisions and ID Uniqueness
Generating an ID takes ~40 ns on a modern CPU thread (~25 000 IDs/ms). Each ID includes a millisecond timestamp and a 15-bit sequence counter (max = 32 767). IDs are guaranteed unique as long as:
- the generator’s node ID is unique, and
- the generation rate does not exceed ~32 767 IDs/ms (~30 ns per ID), preventing sequence overflow within a single millisecond.

---

## Benchmark
```
goos: darwin
goarch: arm64
pkg: github.com/webmafia/hexid
cpu: Apple M1 Pro
BenchmarkGenerator/New-10                   176147698             6.669 ns/op           0 B/op           0 allocs/op
BenchmarkGenerator/ID-10                     28672035            42.550 ns/op           0 B/op           0 allocs/op
BenchmarkGenerator/IDFromTime-10            557195217             2.162 ns/op           0 B/op           0 allocs/op
BenchmarkAtomicGenerator/New-10             180172480             6.674 ns/op           0 B/op           0 allocs/op
BenchmarkAtomicGenerator/ID-10               29088128            44.700 ns/op           0 B/op           0 allocs/op
BenchmarkAtomicGenerator/IDFromTime-10      173841322             6.885 ns/op           0 B/op           0 allocs/op
```

## ⚖️ License

MIT © 2025 The Web Mafia, Ltd.
