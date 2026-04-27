# path benchmarks

Run the benchmark:

```text
make bench name=path
```

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/path
cpu: Apple M1
Benchmark_Join-8          18895888     61.28 ns/op    24 B/op    1 allocs/op
Benchmark_MatchTrue-8     11316279    105.2 ns/op      0 B/op    0 allocs/op
Benchmark_MatchFalse-8    11308365    106.1 ns/op      0 B/op    0 allocs/op
```

So (mimalloc):

```text
Benchmark_Join            16236070     73.01 ns/op    36 B/op    2 allocs/op
Benchmark_MatchTrue       10236286    112.9 ns/op      0 B/op    0 allocs/op
Benchmark_MatchFalse      10443045    113.9 ns/op      0 B/op    0 allocs/op
```

So (arena):

```text
Benchmark_Join            20919043     57.85 ns/op    36 B/op    2 allocs/op
```
