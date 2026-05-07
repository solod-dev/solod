# C interop

So provides several tools for easy C interop.

[Includes](#includes) •
[Extern declarations](#extern-declarations) •
[Extern options](#extern-options) •
[Inlining](#inlining) •
[Qualifiers](#qualifiers-v02) •
[Embeds](#embeds) •
[Raw C](#raw-c) •
[Helpers](#helpers)

## Includes

Include a C header file. By default, `so:include` emits in the `.h` file, making the header visible to consumers:

```go
//so:include <stdint.h>
```

Use `so:include.c` when the include is purely an implementation detail that should only appear in the `.c` file:

```go
//so:include.c "internal_helper.h"
```

## Extern declarations

Declare an external C type (excluded from emission) with `so:extern`:

```go
//so:extern
type Account struct {
    name    string
    balance int64
    flags   []uint8
}
```

Declare an external C function:

```go
//so:extern
func dec_balance(acc *Account, amount int64) int64 {
    return 42 // for testing
}
```

When calling extern functions, `string` and `[]T` arguments are automatically decayed to their C equivalents: string literals become raw C strings (`"hello"`), string values become `char*` (`.ptr`), and slices become raw pointers (`.ptr`). This means C macros don't need to extract `.ptr` themselves:

```go
//so:extern
func fopen(path string, mode string) *File { return nil }

// Go call:
f := fopen("/tmp/test.txt", "w")

// Generated C:
// fopen("/tmp/test.txt", "w")
// not fopen(so_str("/tmp/test.txt"), so_str("w"))
```

The `so:extern` directive supports two optional parameters: a C name override and the `nodecay` flag.

## Extern options

_Name override_ specifies the C name to use instead of the default package-prefixed name. Useful for extern types that must match a C header:

```go
//so:extern Account
type Account struct {
    name    string
    balance int64
}
// Uses "Account" in C instead of "main_Account"
```

_Nodecay_ skips the automatic decay of So types (`so_String`, `so_Slice`) to raw C types. Use this for C functions that are "So-aware" and accept So types directly:

```go
//so:extern nodecay
func set_name(acc *Account, name string)

// Generated C passes so_String directly:
// set_name(&acc, name)
// not set_name(&acc, so_cstr(name))
```

Both options can be combined:

```go
//so:extern MyFunc nodecay
func MyFunc(s string)
```

## Inlining

Force a function to be emitted as `static inline` in the header file using `//so:inline`. This is useful for small, frequently used functions when the compiler won't inline them automatically:

```go
//so:inline
func add(a, b int) int {
    return a + b
}
```

The function body is emitted directly in the `.h` file and skipped from the `.c` file. Works with both functions and methods.

## Qualifiers (v0.2)

### Volatile

Mark a package-level variable as `volatile` using `//so:volatile`:

```go
//so:volatile
var counter int
```

Only allowed on `var` declarations.

### Thread-local storage

Mark a package-level variable as thread-local using `//so:thread_local`. Uses C11 `_Thread_local`:

```go
//so:thread_local
var perThread int
```

Can be combined with `//so:volatile`:

```go
//so:volatile
//so:thread_local
var flags int
```

Only allowed on `var` declarations.

### Attributes

Add GCC/Clang `__attribute__` annotations using `//so:attr`. The text after `so:attr` is used as the attribute value:

```go
//so:attr packed
type header struct {
    version byte
    length  int
}
```

Multiple `//so:attr` lines on the same declaration are combined:

```go
//so:attr packed
//so:attr aligned(16)
type aligned struct {
    x int
}
```

Allowed on `var`, `const`, `type`, and `func` declarations.

## Embeds

Embed C files directly into the generated output using `//so:embed`:

```go
//so:embed main.h
var main_h string

//so:embed main.c
var main_c string
```

`.h` files are embedded into the generated header, `.c` files into the generated implementation. The embed variable declarations are not emitted as C variables - they serve as markers only.

## Raw C (experimental)

For ad-hoc C interop, the `so/c` package provides two compiler intrinsics that emit their string argument as raw C code. The argument must be a string literal.

`c.Val[T](expr)` emits a typed C expression. Use it to access C constants, macros, or call C functions inline:

```go
nan := c.Val[float64]("NAN")
x := c.Val[float64]("sqrt(49)")
```

`c.Raw(code)` emits a raw block of C code as a statement:

```go
var b int
c.Raw(`
int a = 7;
b = a * a;
`)
```

Be careful when using `c.Val` and `c.Raw`. C code written as string literals bypasses the type system and is hard to maintain, so it's usually better to use `so:extern` and `so:embed` instead.

## Helpers

The `so/c` package also provides low-level interop helpers for pointers, strings, and type information.

Functions:

- `Alignof` and `Sizeof` return the alignment and size of type T.
- `Alloca` allocates an array on the stack.
- `Assert` aborts with a message if a condition is false.
- `Bytes`, `Slice` and `String` wrap C pointers to So types.
- `CString` converts a So string to a null-terminated C string.
- `PtrAdd`, `PtrAs` and `PtrAt` manipulate pointers.
- `Zero` returns the zero value of type T.

Types:

- `Char` and `ConstChar` represent a C `char` type.
