# maps benchmarks

Requires GCC/Clang and mimalloc (for heap allocations in So). If mimalloc isn't available, the benchmarks will use the default libc allocator, which is much slower.

Run the benchmark:

```text
make bench name=maps
```

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/maps
cpu: Apple M1

Benchmark_IntSet-8       31677    35580 ns/op     74264 B/op    20 allocs/op
Benchmark_IntGet-8      218179     5573 ns/op         0 B/op     0 allocs/op
Benchmark_IntHas-8      211342     5660 ns/op         0 B/op     0 allocs/op
Benchmark_IntDelete-8    50260    23892 ns/op     36944 B/op     5 allocs/op
Benchmark_IntSetDel-8    63715    18956 ns/op         0 B/op     0 allocs/op

Benchmark_StrSet-8       24082    48677 ns/op    108760 B/op    20 allocs/op
Benchmark_StrGet-8      134481     8990 ns/op         0 B/op     0 allocs/op
Benchmark_StrHas-8      139606    10174 ns/op         0 B/op     0 allocs/op
Benchmark_StrDelete-8    34094    33878 ns/op     54608 B/op     5 allocs/op
Benchmark_StrSetDel-8    45928    26323 ns/op         0 B/op     0 allocs/op
```

So (mimalloc):

```text
Benchmark_IntSet         31851    35406 ns/op     98112 B/op    27 allocs/op
Benchmark_IntGet        744554     1594 ns/op         0 B/op     0 allocs/op
Benchmark_IntHas        753531     1583 ns/op         0 B/op     0 allocs/op
Benchmark_IntDelete      59000    20291 ns/op     73728 B/op     6 allocs/op
Benchmark_IntSetDel     276937     4325 ns/op       192 B/op     3 allocs/op

Benchmark_StrSet         28688    41381 ns/op    130816 B/op    27 allocs/op
Benchmark_StrGet        114166    10568 ns/op         0 B/op     0 allocs/op
Benchmark_StrHas        115996    10250 ns/op         0 B/op     0 allocs/op
Benchmark_StrDelete      37478    32123 ns/op     98304 B/op     6 allocs/op
Benchmark_StrSetDel     130804     9179 ns/op       256 B/op     3 allocs/op
```

So (arena):

```text
Benchmark_IntSet         35107    34483 ns/op     98112 B/op   27 allocs/op
Benchmark_IntGet        710773     1617 ns/op         0 B/op    0 allocs/op
Benchmark_IntHas        755857     1587 ns/op         0 B/op    0 allocs/op
Benchmark_IntDelete      58615    20559 ns/op     73728 B/op    6 allocs/op
Benchmark_IntSetDel     275040     4327 ns/op       192 B/op    3 allocs/op

Benchmark_StrSet         28239    42481 ns/op    130816 B/op   27 allocs/op
Benchmark_StrGet        112918    10438 ns/op         0 B/op    0 allocs/op
Benchmark_StrHas        115912    10230 ns/op         0 B/op    0 allocs/op
Benchmark_StrDelete      37699    35994 ns/op     98304 B/op    6 allocs/op
Benchmark_StrSetDel     129916     9195 ns/op       256 B/op    3 allocs/op
```
