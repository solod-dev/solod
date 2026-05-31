# v0.2 changelog

This document lists the major changes in the [main](https://github.com/solod-dev/solod/commits/main/) branch since the latest release ([v0.1](https://github.com/solod-dev/solod/releases/tag/v0.1.0)).

New directives: `so:volatile`, `so:thread_local`, `so:attr`.<br>
[600e881](https://github.com/solod-dev/solod/commit/600e881fe72cf5f9857745b489c6dedf9a864ea3)

32-bit target support.<br>
[deac815](https://github.com/solod-dev/solod/commit/deac815a5100f119765ffcf8b5961ef579c7a766)
[de30cde](https://github.com/solod-dev/solod/commit/de30cdec169be0f7f8835853ccde5f78e3e4c233)

WebAssembly support (WASI).<br>
[3d0791b](https://github.com/solod-dev/solod/commit/3d0791b69e8fd5053fd508dbbb8c9cebfb0b3ff7)

Freestanding mode (no libc dependency).<br>
[1cfc8c7](https://github.com/solod-dev/solod/commit/1cfc8c7cd602a379332e6c128ebd2bde007c9a63)

Implement `error` as a regular interface (it was special-cased before).<br>
[6c8f0bd](https://github.com/solod-dev/solod/commit/6c8f0bd68e4ba8693d22be59f763676889270070)

⚠️ Auto-generate result types for `(T, error)` multi-return values where `T` is a custom struct type. This is a breaking change: previously, you had to manually define a `T{Result}` type for any such `T`.<br>
[745b174](https://github.com/solod-dev/solod/commit/745b174e11c08bee91cdaeaf8a8b2aa083863b61)

Numeric C types in the `so/c` package for better interop.<br>
[5914a75](https://github.com/solod-dev/solod/commit/5914a7591bc44335b4556893a3f848e1a6c9cc8c)

Type aliases.<br>
[deeccb9](https://github.com/solod-dev/solod/commit/deeccb98d22f342e6bdceb9d7827e9d464af9603)

Support for third-party packages (`go install` or vendoring) and multi-module projects.<br>
[bba8265](https://github.com/solod-dev/solod/commit/bba8265883b10814803510518693b224b70d2d98)

`net/netip` package.<br>
[5f87292](https://github.com/solod-dev/solod/commit/5f87292ac5b334cd1080a20ca75cc5a1c2c3ea59)

`mem/Arena.Free` reclaims the last allocation if the pointer matches.<br>
[f9adba6](https://github.com/solod-dev/solod/commit/f9adba6baca67b2ed332e3aaaf1e59b44113d1db)

Report So source location (file and line number) instead of C source location when panicking (`track-source` flag, off by default).<br>
[fb78b7a](https://github.com/solod-dev/solod/commit/fb78b7af20525055e320e5f01cb5bb8198ab18ff)

Check for nil pointer dereference on struct pointer field access and interface method calls (`check-nil` flag, off by default).<br>
[426961e](https://github.com/solod-dev/solod/commit/426961e0ef463cc2390e6d1a930555f2db581f7e)

⚠️ Block-scoped `defer` is no longer supported. This is a breaking change.<br>
[fb49cca](https://github.com/solod-dev/solod/commit/fb49ccab2316815308f690f2690e1c3bf19ee59b)

`encoding/hex` package.<br>
[42a5cf0](https://github.com/solod-dev/solod/commit/42a5cf0d7f6f08ba9e862bd5a738cf74448b2711)

`uuid` package.<br>
[fc8f2fa](https://github.com/solod-dev/solod/commit/fc8f2fabdac147f576fcede45bd178b313a7e25a)
