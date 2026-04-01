# io benchmarks

Requires GCC/Clang and mimalloc (for heap allocations in So). If mimalloc isn't available, the benchmarks will use the default libc allocator, which is much slower.

Run the benchmark:

```text
make bench name=bytes
```

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/io
cpu: Apple M1
BenchmarkCopyNSmall-8    3263336      487.0 ns/op      1340 B/op    1 allocs/op
BenchmarkCopyNLarge-8      52023    21419 ns/op      115350 B/op    2 allocs/op
```

So (mimalloc):

```text
Benchmark_CopyNSmall     9174732      419.3 ns/op      1872 B/op    0 allocs/op
Benchmark_CopyNLarge       75915    16004 ns/op      113151 B/op    0 allocs/op
```
