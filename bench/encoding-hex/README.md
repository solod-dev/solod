# encoding/hex benchmarks

Run the benchmark:

```text
make bench name=encoding-hex
```

### Encode

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/encoding-binary
cpu: Apple M1
Benchmark_Encode/256-8      5442978      193.2 ns/op    1325.30 MB/s
Benchmark_Encode/1024-8     1627232      734.5 ns/op    1394.06 MB/s
Benchmark_Encode/4096-8      413481     2940 ns/op      1393.37 MB/s
Benchmark_Encode/16384-8     101353    11617 ns/op      1410.39 MB/s
```

So 0.2:

```text
Benchmark_Encode_256        6328112      171.3 ns/op    1494.17 MB/s
Benchmark_Encode_1024       1842578      656.1 ns/op    1560.65 MB/s
Benchmark_Encode_4096        466254     2607 ns/op      1571.08 MB/s
Benchmark_Encode_16384       116136    10408 ns/op      1574.11 MB/s
```

### Decode

Go 1.26.1:

```text
goos: darwin
goarch: arm64
pkg: solod.dev/bench/encoding-binary
cpu: Apple M1
BenchmarkDecode/256-8       9465538      127.2 ns/op    2011.83 MB/s
BenchmarkDecode/1024-8      2437354      502.2 ns/op    2038.83 MB/s
BenchmarkDecode/4096-8       610184     1963 ns/op      2086.39 MB/s
BenchmarkDecode/16384-8      154268     7797 ns/op      2101.36 MB/s
```

So 0.2:

```text
Benchmark_Decode_256       11874603      96.19 ns/op    2661.28 MB/s
Benchmark_Decode_1024       3400048     358.2 ns/op     2859.02 MB/s
Benchmark_Decode_4096        780894    1422 ns/op       2879.45 MB/s
Benchmark_Decode_16384       216157    5600 ns/op       2925.72 MB/s
```
