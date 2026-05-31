# Solod vs. Go benchmarks

Here are some benchmarks that show how So performs on common tasks compared to Go.

[Bufio](#buffered-io) •
[Bytes](#byte-functions) •
[Maps](#maps) •
[Strings](#string-conversion) •
[Time](#time)

See the benchmarks for individual packages in the nested folders:

[bufio](./bufio/README.md) •
[bytes](./bytes/README.md) •
[encoding/binary](./encoding-binary/README.md) •
[encoding/hex](./encoding-hex/README.md) •
[io](./io/README.md) •
[log/slog](./log-slog/README.md) •
[maps](./maps/README.md) •
[math/rand](./math-rand/README.md) •
[net/netip](./net-netip/README.md) •
[path](./path/README.md) •
[strconv](./strconv/README.md) •
[strings](./strings/README.md) •
[time](./time/README.md)

## Buffered I/O

So is ~3x faster than Go for reading and writing, and ~4x faster for scanning.

| Benchmark           |     Go |     So | Winner        |
| ------------------- | -----: | -----: | ------------- |
| Reader (buffered)   | 3089ns | 1073ns | **So** - 2.9x |
| Reader (unbuffered) | 1269ns |  412ns | **So** - 3.1x |
| Writer (buffered)   | 2988ns | 1038ns | **So** - 2.9x |
| Writer (unbuffered) | 4928ns | 1537ns | **So** - 3.2x |
| Scanner             |  443ns |  112ns | **So** - 4.0x |

Apple M1 • Go 1.26.1 • [details](./bufio/README.md)

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

Apple M1 • Go 1.26.1 • [details](./bytes/README.md#functions)

## Byte buffer

So reads 1.3x faster and writes 2-4x faster than Go.
Memory usage is the same for both.

| Benchmark  |      Go | So (mimalloc) | So (arena) | Winner        |
| ---------- | ------: | ------------: | ---------: | ------------- |
| ReadString |  2329ns |        1757ns |     1719ns | **So** - 1.3x |
| WriteByte  |  8858ns |        2608ns |     2643ns | **So** - 3.4x |
| WriteRune  | 15110ns |        3902ns |     3956ns | **So** - 3.8x |
| WriteBlock | 17238ns |        7830ns |     7510ns | **So** - 2.2x |

Apple M1 • Go 1.26.1 • [details](./bytes/README.md#buffer)

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

Apple M1 • Go 1.26.1 • [details](./maps/README.md)

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

Apple M1 • Go 1.26.1 • [details](./strconv/README.md)

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

Apple M1 • Go 1.26.1 • [details](./strings/README.md#functions)

## String builder

So is 2-4x faster than Go and uses 10%-20% less memory.

| Benchmark                |    Go | So (mimalloc) | So (arena) | Winner        |
| ------------------------ | ----: | ------------: | ---------: | ------------- |
| Write bytes (auto-grow)  | 245ns |         118ns |       59ns | **So** - 2.1x |
| Write bytes (pre-grow)   | 109ns |          29ns |       25ns | **So** - 3.8x |
| Write string (auto-grow) | 224ns |         116ns |       57ns | **So** - 1.9x |
| Write string (pre-grow)  | 113ns |          29ns |       26ns | **So** - 3.9x |

Apple M1 • Go 1.26.1 • [details](./strings/README.md#builder)

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

Apple M1 • Go 1.26.1 • [details](./time/README.md)

## Methodology

So is compiled with `-Ofast -march=native -flto -funroll-loops` and uses mimalloc as the system allocator. Go is run with default `go test -bench=.` settings.

The Winner column shows the worse result between mimalloc and arena for each So benchmark.
