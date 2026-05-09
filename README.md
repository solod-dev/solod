# Solod: Go can be a better C

Solod (**So**) is a strict subset of Go that translates to regular C.

Highlights:

- Go in, C out. You write regular Go code and get readable C11 as output.
- Zero runtime. No garbage collection, no reference counting, no hidden allocations.
- Rich standard library. Use familiar types and functions ported from Go's stdlib.
- Native C interop. Call C from So and So from C — no CGO, no overhead.
- Go tooling works out of the box. Syntax highlighting, LSP, linting and "go test".

So supports structs, methods, interfaces, slices, maps, multiple returns, and defer. Everything is stack-allocated by default; heap is opt-in through the standard library. To keep things simple, there are no channels, goroutines, closures, or generics.

So is for Go developers who want systems-level control without learning a new language. And for C programmers who like Go's safety, structure, and tooling.

[Example](#example) •
[Installation](#installation-and-usage) •
[Documentation](#documentation-and-examples) •
[Benchmarks](#testing-and-benchmarks) •
[Compatibility](#compatibility) •
[Design](#design-principles-and-faq) •
[Roadmap](#roadmap) •
[Contributing](#contributing)

## Example

This Go code in a file `main.go`:

```go
package main

import "solod.dev/so/time"

type Person struct {
    Name string
    Age  int
    Nums [3]int
}

func (p *Person) Sleep() int {
    p.Age += 1
    return p.Age
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    p.Sleep()
    println(p.Name, "is now", p.Age, "years old.")

    p.Nums[0] = 42
    println("1st lucky number is", p.Nums[0])

    year := time.Now().Year()
    println("The year is", year)
}
```

Translates to a header file `main.h`:

```c
#pragma once
#include "so/builtin/builtin.h"
#include "so/time/time.h"

typedef struct main_Person {
    so_String Name;
    so_int Age;
    so_int Nums[3];
} main_Person;

so_int main_Person_Sleep(void* self);
```

Plus an implementation file `main.c`:

```c
#include "main.h"

so_int main_Person_Sleep(void* self) {
    main_Person* p = self;
    p->Age += 1;
    return p->Age;
}

int main(void) {
    main_Person p = (main_Person){.Name = so_str("Alice"), .Age = 30};
    main_Person_Sleep(&p);
    so_println("%.*s %s %" PRId64 " %s", p.Name.len, p.Name.ptr, "is now", p.Age, "years old.");

    p.Nums[0] = 42;
    so_println("%s %" PRId64, "1st lucky number is", p.Nums[0]);

    so_int year = time_Time_Year(time_Now());
    so_println("%s %" PRId64, "The year is", year);
}
```

Check out more examples in [So by example](https://github.com/solod-dev/example) and learn about the supported language features in the [language tour](doc/spec.md).

## Installation and usage

Install the So command line tool:

```
go install solod.dev/cmd/so@latest
```

Create a new Go project and add the Solod dependency to use the So standard library:

```
go mod init example
go get solod.dev@latest
```

Write regular Go code, but use So packages instead of the standard Go packages:

```go
package main

import "solod.dev/so/math"

func main() {
    ans := math.Sqrt(1764)
    println("Hello, world! The answer is", int(ans))
}
```

Transpile to C:

```
so translate -o generated .
```

The translated C code will be saved in the `generated` directory.

You can also transpile to C and compile the code to a binary in one step. This uses the C compiler set by the `CC` environment variable (default `cc`):

```
so build -o main .
```

Or you can transpile, compile, and run without saving the binary:

```
so run .
```

All commands work with Go modules, not individual files (`so run .`, not `so run main.go`).

You can pass additional compiler and linker flags via `CFLAGS` and `LDFLAGS`:

```
CFLAGS="-O2" LDFLAGS="-lm" so build -o main .
```

On Linux (and some BSDs), the math library is not linked by default. If your program imports `so/math` — directly or through other packages like `so/io` — you'll need to add `-lm`:

```
LDFLAGS="-lm" so build -o main .
```

Keep in mind that So is new, so it's still a bit rough around the edges.

## Documentation and examples

**[Language tour](./doc/spec.md)**. To learn about So's features and limitations, check out the brief overview of the language.

**[Standard library](./doc/stdlib.md)**. So provides a growing set of packages similar to Go's stdlib.

**[Playground](https://codapi.org/so/)**. Try So online without installing anything. You can run the code or view the translated C output.

**[So by example](https://github.com/solod-dev/example)**. If you like learning by doing, try a hands-on introduction to So with annotated example programs.

**[AI skill](https://github.com/solod-dev/ai)**. You can have a clanker write So code for you. But where's the fun in that?

## Testing and benchmarks

**Testing**. So doesn't have its own testing framework. Since So code is valid Go code, you can just use `go test` like you normally would. Plus, your tests can use all Go features because they're never transpiled.

**[Benchmarks](bench/README.md)**. So truly shines when it comes to C interop, but it's also quite fast on regular Go code — typically on par with or faster than Go.

## Compatibility

So generates C11 code that relies on several GCC/Clang extensions:

- Binary literals (`0b1010`) in generated code.
- Statement expressions (`({...})`) in macros.
- `__attribute__((constructor))` for package-level initialization.
- `__auto_type` for local type inference in generated code.
- `__typeof__` for type inference in generic macros.
- `alloca` and VLAs for `make()` and other dynamic stack allocations.

Supported compilers: GCC, Clang, Emscripten, and `zig cc`. MSVC is not supported.

Supported operating systems: Linux, macOS, and Windows (core language only).

Supported platforms: amd64, arm64, riscv64, i386, and wasm32.

## Design principles and FAQ

**[Principles](doc/design.md)**. So is highly opinionated. Simplicity is key. Heap allocations are explicit. Strictly Go syntax.

**[Frequently asked questions](doc/faq.md)**. I have heard these several times, so it's worth answering.

## Roadmap

✅ Core language features.

✅ Core stdlib packages (v0.1):

```text
✓ bufio    ✓ io      ✓ path      ✓ strings
✓ bytes    ✓ maps    ✓ rand      ✓ strconv
✓ flag     ✓ math    ✓ slices    ✓ time
✓ fmt      ✓ os      ✓ slog      ✓ unicode
```

⏳ I'm currently gathering feedback and defining the scope for the [v0.2 release](doc/changelog.md).

⬜ Networking (v0.2).

⬜ Concurrency (v0.3).

🤔 32-bit targets.

🤔 Full Windows support.

## Contributing

Bug fixes are welcome. For anything other than bug fixes, please open an issue first to discuss your proposed changes. To prevent feature bloat, it's important to discuss any new features before adding them.

AI-assisted submissions are fine on one condition: you, the human, have read all the code and fully understand what it does. Code reviewed only by another AI will not suffice.

Make sure to add or update tests as needed.

## License

Go stdlib code by the [Go Authors](https://github.com/golang/go).

Transpiler and So stdlib code by [Anton Zhiyanov](https://antonz.org/).

Released under the BSD 3-Clause License.
