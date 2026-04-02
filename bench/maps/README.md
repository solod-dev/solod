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
Benchmark_Set-8  31677  35580 ns/op  74264 B/op  20 allocs/op
Benchmark_Get-8  218179  5573 ns/op  0 B/op  0 allocs/op
Benchmark_Has-8  211342  5660 ns/op  0 B/op  0 allocs/op
Benchmark_Delete-8  50260  23892 ns/op  36944 B/op  5 allocs/op
Benchmark_SetDelete-8  63715  18956 ns/op  0 B/op  0 allocs/op
```

So (mimalloc):

```text
Benchmark_Set     17412       66479 ns/op     98112 B/op        27 allocs/op
Benchmark_Get    730503        1605 ns/op         0 B/op         0 allocs/op
Benchmark_Has    744277        1600 ns/op         0 B/op         0 allocs/op
Benchmark_Delete     33830       35621 ns/op     73728 B/op         6 allocs/op
Benchmark_SetDelete    277848        4366 ns/op       192 B/op         3 allocs/op
```

So (arena):

```text
Benchmark_Set     18577       64354 ns/op     98112 B/op        27 allocs/op
Benchmark_Get    743310        1668 ns/op         0 B/op         0 allocs/op
Benchmark_Has    753578        1606 ns/op         0 B/op         0 allocs/op
Benchmark_Delete     33637       36772 ns/op     73728 B/op         6 allocs/op
Benchmark_SetDelete    262053        4439 ns/op       192 B/op         3 allocs/op
```
