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
