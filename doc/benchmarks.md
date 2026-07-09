# Solod vs. Go benchmarks

Here are some benchmarks that show how So performs on common tasks compared to Go.

[bufio](#buffered-io) •
[bytes](#byte-functions) •
[conc](#concurrency) •
[crypto/crand](#cryptographic-random) •
[encoding/binary](#binary-encoding) •
[encoding/hex](#hex-encoding) •
[io](#stream-copying) •
[log/slog](#structured-logging) •
[maps](#maps) •
[math/rand](#pseudorandom-numbers) •
[net/netip](#ip-addresses) •
[path](#path-manipulation) •
[strconv](#string-conversion) •
[strings](#string-functions) •
[sync](#synchronization) •
[time](#time) •
[uuid](#uuid)

## Buffered I/O

So is ~3x faster than Go for reading and writing, and ~4x faster for scanning.

| Benchmark           |     Go |     So | Winner        |
| ------------------- | -----: | -----: | ------------- |
| Reader (buffered)   | 3089ns | 1073ns | **So** - 2.9x |
| Reader (unbuffered) | 1269ns |  412ns | **So** - 3.1x |
| Writer (buffered)   | 2988ns | 1038ns | **So** - 2.9x |
| Writer (unbuffered) | 4928ns | 1537ns | **So** - 3.2x |
| Scanner             |  443ns |  112ns | **So** - 4.0x |

Apple M1 • Go 1.26.1

## Byte functions

So is generally ~1.5x faster than Go, except for Index operations.
Memory usage is the same for both.

| Benchmark  |    Go | So (mimalloc) | So (arena) | Winner        |
| ---------- | ----: | ------------: | ---------: | ------------- |
| Clone      | 102ns |          41ns |       32ns | **So** - 2.5x |
| Compare    |  34ns |          25ns |       25ns | **So** - 1.4x |
| Index      |  21ns |          32ns |       32ns | Go - 0.7x     |
| IndexByte  |  16ns |          25ns |       25ns | Go - 0.6x     |
| Repeat     | 106ns |          56ns |       48ns | **So** - 1.9x |
| ReplaceAll | 247ns |         258ns |      242ns | ~same         |
| Split      | 510ns |         422ns |      421ns | **So** - 1.2x |
| ToUpper    | 322ns |         176ns |      171ns | **So** - 1.8x |
| Trim       |  47ns |          44ns |       44ns | **So** - 1.1x |
| TrimSuffix |   4ns |           2ns |        2ns | **So** - 1.8x |

Apple M1 • Go 1.26.1

## Byte buffer

So reads 1.3x faster and writes 2-4x faster than Go.
Memory usage is the same for both.

| Benchmark  |      Go | So (mimalloc) | So (arena) | Winner        |
| ---------- | ------: | ------------: | ---------: | ------------- |
| ReadString |  2329ns |        1757ns |     1719ns | **So** - 1.3x |
| WriteByte  |  8858ns |        2608ns |     2643ns | **So** - 3.4x |
| WriteRune  | 15110ns |        3902ns |     3956ns | **So** - 3.8x |
| WriteBlock | 17238ns |        7830ns |     7510ns | **So** - 2.2x |

Apple M1 • Go 1.26.1

## Concurrency

`conc.Pool` is a fixed set of worker threads draining a shared task queue, built
on So's `Mutex` and `Cond`. Each dispatch crosses into the kernel to wake a
worker (see [Cond](#cond)), so the pool suits coarse-grained tasks: on realistic
workloads that per-task cost is amortized and So stays within ~1.1x of Go.

The benchmarks run 8 workers on both sides - So's `conc.Pool` against an
equivalent Go pool of persistent goroutines draining a buffered channel. Each
CPU-bound task runs computations of ~40µs; each IO-bound task blocks for 1ms,
standing in for a network or disk round-trip.

| Benchmark        |  Go |   So | Winner    |
| ---------------- | --: | ---: | --------- |
| Work (CPU-bound) | 7ms |  8ms | Go - 0.9x |
| IO (IO-bound)    | 9ms | 10ms | Go - 0.9x |

For CPU-bound work So's faster compute nearly offsets its heavier dispatch; for
IO-bound work the dispatch cost hides behind the blocking waits. Note the pool
is capped at `NumThreads` OS threads, so unlike Go's goroutines it cannot fan a
single batch out to thousands of concurrent IO waits.

Apple M1 • Go 1.26.1

## Cryptographic random

So is faster than Go for small reads and random text, and about the same for large reads.

| Benchmark |     Go |     So | Winner        |
| --------- | -----: | -----: | ------------- |
| Read 4B   |   69ns |   40ns | **So** - 1.7x |
| Read 32B  |  242ns |  211ns | **So** - 1.1x |
| Read 4KB  | 1215ns | 1184ns | ~same         |
| Text      |  264ns |  213ns | **So** - 1.2x |

Apple M1 • Go 1.26.1

## Binary encoding

So encodes fixed-size integers about 2x faster than Go.

| Benchmark       |     Go |     So | Winner        |
| --------------- | -----: | -----: | ------------- |
| BE PutUint64    | 0.63ns | 0.32ns | **So** - 2.0x |
| BE AppendUint64 | 1.77ns | 0.95ns | **So** - 1.9x |
| LE PutUint64    | 0.63ns | 0.31ns | **So** - 2.0x |
| LE AppendUint64 | 1.73ns | 0.95ns | **So** - 1.8x |

Apple M1 • Go 1.26.1

## Hex encoding

So encodes ~1.1x and decodes ~1.4x faster than Go. The ratios hold across buffer sizes from 256B to 16KB; representative figures:

| Benchmark   |     Go |     So | Winner        |
| ----------- | -----: | -----: | ------------- |
| Encode 256B |  193ns |  171ns | **So** - 1.1x |
| Encode 4KB  | 2940ns | 2607ns | **So** - 1.1x |
| Decode 256B |  127ns |   96ns | **So** - 1.3x |
| Decode 4KB  | 1963ns | 1422ns | **So** - 1.4x |

Apple M1 • Go 1.26.1

## Stream copying

So's `io.CopyN` is ~1.2-1.3x faster than Go and, routed through an allocator, reports no per-op allocations. So uses mimalloc.

| Benchmark   |      Go |      So | Winner        |
| ----------- | ------: | ------: | ------------- |
| CopyN small |   487ns |   419ns | **So** - 1.2x |
| CopyN large | 21419ns | 16004ns | **So** - 1.3x |

Apple M1 • Go 1.26.1

## Structured logging

So is 4-7x faster than Go, and logging with attributes allocates nothing in So versus three allocations in Go.

| Benchmark       |    Go |   So | Winner        |
| --------------- | ----: | ---: | ------------- |
| No attributes   | 166ns | 39ns | **So** - 4.3x |
| With attributes | 259ns | 38ns | **So** - 6.8x |

Apple M1 • Go 1.26.1

## Maps

### Int keys

For heap-allocated maps, So is ~1.4x faster than Go across all operations.

So's built-in map is even faster, but it's only useful in certain situations — it's fixed size and stack-allocated.

| Benchmark |      Go | So (mimalloc) | So (arena) | So (built-in) | Winner        |
| --------- | ------: | ------------: | ---------: | ------------: | ------------- |
| Set       | 35645ns |       26333ns |    25515ns |           n/a | **So** - 1.4x |
| Set (pre) |  9676ns |        8813ns |     8704ns |        3109ns | **So** - 1.1x |
| Get       |  5594ns |        1581ns |     1537ns |        2577ns | **So** - 3.5x |
| Delete    | 23968ns |       14889ns |    14859ns |           n/a | **So** - 1.6x |

### String keys

So modifications are ~1.4x faster than Go, while lookups are slightly slower.

| Benchmark |      Go | So (mimalloc) | So (arena) | So (built-in) | Winner        |
| --------- | ------: | ------------: | ---------: | ------------: | ------------- |
| Set       | 47805ns |       31055ns |    30749ns |           n/a | **So** - 1.5x |
| Set (pre) | 14699ns |       12101ns |    12233ns |        6585ns | **So** - 1.2x |
| Get       |  9216ns |       10170ns |     9907ns |       10531ns | Go - 0.9x     |
| Delete    | 33819ns |       24227ns |    24392ns |           n/a | **So** - 1.4x |

Apple M1 • Go 1.26.1

## Pseudorandom numbers

So's raw source generator is ~1.6x faster, but the package-level helpers (global source, bounded ints, floats) are about 2x slower than Go.

| Benchmark     |    Go |    So | Winner        |
| ------------- | ----: | ----: | ------------- |
| Source Uint64 | 4.7ns | 2.8ns | **So** - 1.6x |
| Global Uint64 | 4.8ns | 8.8ns | Go - 0.5x     |
| Uint64        | 4.5ns | 8.8ns | Go - 0.5x     |
| Int64N (1e9)  | 4.6ns | 9.1ns | Go - 0.5x     |
| Int64N (4e18) | 9.1ns |  12ns | Go - 0.8x     |
| Float64       | 4.4ns | 9.3ns | Go - 0.5x     |

Apple M1 • Go 1.26.1

## IP addresses

So parses IPv6 ~1.4-1.5x faster and formats addresses 2-4x faster than Go, allocating nothing. The exception is parsing a zoned IPv6 address, which makes an `if_nametoindex` syscall and is far slower.

Parsing:

| Benchmark     |   Go |      So | Winner        |
| ------------- | ---: | ------: | ------------- |
| Parse v4      | 18ns |    16ns | **So** - 1.1x |
| Parse v6      | 81ns |    55ns | **So** - 1.5x |
| Parse v6e     | 47ns |    33ns | **So** - 1.4x |
| Parse v6+v4   | 48ns |    40ns | **So** - 1.2x |
| Parse v6+zone | 64ns | 19087ns | Go - syscall  |

Formatting:

| Benchmark      |   Go |   So | Winner        |
| -------------- | ---: | ---: | ------------- |
| String v4      | 20ns |  9ns | **So** - 2.3x |
| String v6      | 53ns | 17ns | **So** - 3.1x |
| String v6+v4   | 23ns | 11ns | **So** - 2.0x |
| String v6+zone | 60ns | 14ns | **So** - 4.3x |

Apple M1 • Go 1.26.1

## Path manipulation

Slash paths are roughly on par with Go. Matching is marginally slower in So; Join is slower with mimalloc but faster with an arena.

| Benchmark   |    Go | So (mimalloc) | So (arena) | Winner    |
| ----------- | ----: | ------------: | ---------: | --------- |
| Join        |  61ns |          73ns |       58ns | Go - 0.8x |
| Match true  | 105ns |         113ns |        n/a | Go - 0.9x |
| Match false | 106ns |         114ns |        n/a | Go - 0.9x |

Apple M1 • Go 1.26.1

## String conversion

### Parsing

So parses floats ~1.5x faster and ints ~2x faster than Go.

| Benchmark       |   Go |   So | Winner        |
| --------------- | ---: | ---: | ------------- |
| Atof64 decimal  | 21ns | 12ns | **So** - 1.7x |
| Atof64 float    | 24ns | 15ns | **So** - 1.6x |
| Atof64 exp      | 25ns | 21ns | **So** - 1.2x |
| Atof64 big      | 38ns | 25ns | **So** - 1.5x |
| ParseInt 7-bit  | 10ns |  4ns | **So** - 2.5x |
| ParseInt 26-bit | 14ns |  7ns | **So** - 2.0x |
| ParseInt 31-bit | 16ns |  9ns | **So** - 1.9x |
| ParseInt 56-bit | 24ns | 15ns | **So** - 1.6x |
| ParseInt 62-bit | 26ns | 17ns | **So** - 1.6x |

### Formatting

So formats floats ~1.2x faster and ints ~2x faster than Go.

| Benchmark           |   Go |   So | Winner        |
| ------------------- | ---: | ---: | ------------- |
| FormatFloat decimal | 30ns | 27ns | **So** - 1.1x |
| FormatFloat float   | 43ns | 34ns | **So** - 1.3x |
| FormatFloat exp     | 35ns | 30ns | **So** - 1.2x |
| FormatFloat big     | 39ns | 33ns | **So** - 1.2x |
| FormatInt 7-bit     | 14ns |  5ns | **So** - 3.0x |
| FormatInt 26-bit    | 17ns |  7ns | **So** - 2.3x |
| FormatInt 31-bit    | 20ns |  8ns | **So** - 2.3x |
| FormatInt 56-bit    | 24ns | 12ns | **So** - 2.0x |
| FormatInt 62-bit    | 26ns | 13ns | **So** - 2.0x |

Apple M1 • Go 1.26.1

## String functions

So is generally ~1.3x faster than Go, except for Index operations.
Memory usage is the same for both.

| Benchmark  |     Go | So (mimalloc) | So (arena) | Winner        |
| ---------- | -----: | ------------: | ---------: | ------------- |
| Clone      |   99ns |          42ns |       34ns | **So** - 2.4x |
| Compare    |   47ns |          36ns |       36ns | **So** - 1.3x |
| Fields     | 1524ns |         908ns |      912ns | **So** - 1.7x |
| Index      |   25ns |          35ns |       34ns | Go - 0.7x     |
| IndexByte  |   22ns |          33ns |       33ns | Go - 0.7x     |
| Repeat     |  127ns |          64ns |       67ns | **So** - 1.9x |
| ReplaceAll |  243ns |         200ns |      203ns | **So** - 1.2x |
| Split      | 1899ns |        1399ns |     1423ns | **So** - 1.3x |
| ToUpper    | 2066ns |        1602ns |     1622ns | **So** - 1.3x |
| Trim       |  501ns |         373ns |      375ns | **So** - 1.3x |

Apple M1 • Go 1.26.1

## String builder

So is 2-4x faster than Go and uses 10%-20% less memory.

| Benchmark                |    Go | So (mimalloc) | So (arena) | Winner        |
| ------------------------ | ----: | ------------: | ---------: | ------------- |
| Write bytes (auto-grow)  | 245ns |         118ns |       59ns | **So** - 2.1x |
| Write bytes (pre-grow)   | 109ns |          29ns |       25ns | **So** - 3.8x |
| Write string (auto-grow) | 224ns |         116ns |       57ns | **So** - 1.9x |
| Write string (pre-grow)  | 113ns |          29ns |       26ns | **So** - 3.9x |

Apple M1 • Go 1.26.1

## Synchronization

So's synchronization primitives are built on POSIX threads: `Mutex` and `Cond`
wrap a pthread mutex and condition variable. The mutex beats Go's for short,
spin-friendly critical sections but loses once contention forces threads to park
in the kernel. `Cond` is slower because it always parks threads in the kernel
instead of a user-space scheduler. `Once` takes a lock-free atomic fast path,
so uncontended it is close to Go; under contention it inherits the same kernel
dispatch cost as `Cond`.

The contended benchmarks run 8 worker threads that share one primitive, using a
persistent thread pool on the So side and an equivalent persistent goroutine pool
on the Go side.

### Mutex

Uncontended lock/unlock is ~1.6x faster than Go. Under contention the result
depends on how long the lock is held. With an empty critical section (the _spin_
row) a waiting thread reacquires the lock while still spinning and almost never
parks, so So's thin pthread wrapper wins by ~2.8x. Give the critical section a
small (~1µs) amount of real work (the _work_ row), and waiters exhaust their spin
budget and park in the kernel; every handoff then costs a wakeup syscall, and So
drops to ~0.5x of Go. The _work_ critical section runs identically on both sides
single-threaded, so the gap is purely the parking cost, not the work.

| Benchmark           |    Go |    So | Winner        |
| ------------------- | ----: | ----: | ------------- |
| Uncontended         |  14ns |   9ns | **So** - 1.6x |
| TryLock             |  14ns |   9ns | **So** - 1.6x |
| Contended spin (8t) | 600µs | 215µs | **So** - 2.8x |
| Contended work (8t) |   9ms |  16ms | Go - 0.5x     |

### Cond

So's condition variable is ~7-10x slower than Go across waiter counts: each
wakeup crosses into the kernel, while Go wakes goroutines in user space. Figures
are per 1000 rendezvous rounds.

| Benchmark  |     Go |    So | Winner     |
| ---------- | -----: | ----: | ---------- |
| 1 waiter   | 0.15ms | 1.5ms | Go - 0.10x |
| 2 waiters  | 0.39ms | 2.9ms | Go - 0.13x |
| 4 waiters  | 0.87ms | 7.3ms | Go - 0.12x |
| 8 waiters  |  2.0ms |  14ms | Go - 0.15x |
| 16 waiters |  3.9ms |  28ms | Go - 0.14x |
| 32 waiters |  9.0ms |  60ms | Go - 0.15x |

### Once

So's `Do` takes a lock-free atomic fast path: once the initializer has run,
every call is just an atomic load. Uncontended, both sides do that single load
and land within ~1.2x of each other. Under contention the gap is because of
`conc.Pool` dispatch: waking the eight workers crosses into the kernel, the
same cost that makes `Cond` slow, rather than anything in `Once`.

| Benchmark      |    Go |    So | Winner    |
| -------------- | ----: | ----: | --------- |
| Uncontended    | 2.1ns | 2.6ns | Go - 0.8x |
| Contended (8t) | 6.0µs |  32µs | Go - 0.2x |

### Atomic

So's atomic types map directly to the compiler's `__atomic` builtins - the same
hardware instructions Go emits - so performance is on par with Go across the board.

Single-value ops use `Uint64`; the contended row runs 8 threads adding to one counter.

| Benchmark       |    Go |    So | Winner |
| --------------- | ----: | ----: | ------ |
| Load            |   2ns |   2ns | ~same  |
| Store           |   2ns |   2ns | ~same  |
| Add             |   7ns |   7ns | ~same  |
| Swap            |   7ns |   6ns | ~same  |
| CompareAndSwap  |  13ns |  13ns | ~same  |
| Add (8 threads) | 180µs | 180µs | ~same  |

Apple M1 • Go 1.26.1

## Time

Regular time functions and methods in So are slightly slower than in Go.
In parsing and formatting, So is 5x faster for predefined layouts (RFC3339, DateTime, etc.),
about the same for custom parsing, and 5x slower for custom formatting (due to strftime overhead).

| Benchmark    |   Go |    So | Winner        |
| ------------ | ---: | ----: | ------------- |
| Date         |  7ns |   2ns | **So** - 3.2x |
| ISOWeek      |  9ns |   2ns | **So** - 4.3x |
| Now          | 34ns |  39ns | Go - 0.9x     |
| Since        | 17ns |  25ns | Go - 0.7x     |
| UnixNano     | 34ns |  38ns | Go - 0.9x     |
| Until        | 17ns |  24ns | Go - 0.7x     |
| Format       | 39ns |   4ns | **So** - 8.8x |
| FormatCustom | 55ns | 250ns | Go - 0.2x     |
| Parse        | 27ns |   6ns | **So** - 4.9x |
| ParseCustom  | 55ns |  45ns | **So** - 1.2x |

Apple M1 • Go 1.26.1

## UUID

So generates v4 UUIDs a bit faster and formats them ~4x faster; v7 generation and parsing are on par with Go.

| Benchmark     |    Go |    So | Winner        |
| ------------- | ----: | ----: | ------------- |
| NewV4         | 251ns | 212ns | **So** - 1.2x |
| NewV7         |  72ns |  79ns | Go - 0.9x     |
| String        |  34ns |   9ns | **So** - 3.9x |
| Parse (ok)    |  29ns |  29ns | ~same         |
| Parse (error) |  26ns |  29ns | Go - 0.9x     |

Apple M1 • Go 1.27

## Methodology

So is compiled with `-Ofast -march=native -flto -funroll-loops` and uses mimalloc as the system allocator. Go is run with default `go test -bench=.` settings.

The Winner column shows the worse result between mimalloc and arena for each So benchmark.
