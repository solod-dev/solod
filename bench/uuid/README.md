# uuid benchmarks

Run the benchmark:

```text
make bench name=uuid
```

Go 1.27:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/uuid
cpu: Apple M1
BenchmarkNewV4-8            4347639    250.8 ns/op
BenchmarkNewV7-8           16745724     72.05 ns/op
BenchmarkString-8          34668553     33.66 ns/op
BenchmarkParseSuccess-8    41154498     29.16 ns/op
BenchmarkParseError-8      47801780     25.61 ns/op
```

So 0.2:

```text
Benchmark_NewV4             5078676    211.8 ns/op
Benchmark_NewV7            15316867     78.77 ns/op
Benchmark_String          138107875      8.715 ns/op
Benchmark_ParseSuccess     41671006     29.41 ns/op
Benchmark_ParseError       41529675     28.88 ns/op
```
