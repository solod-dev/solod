# Soan: the "better C" is just C

**Soan** is a subset of Go that transpiles to regular C with zero runtime.

You write regular Go code — structs, methods, interfaces, slices, multiple returns, defer — and get plain C11 code as output. There's no garbage collector or reference counting. Everything is stack-allocated: slices have a fixed capacity, strings are immutable pointer-length pairs, and interfaces are inline vtable structs. Heap allocations will be handled by the standard library, not built into the language runtime.

There are no maps, channels, goroutines, closures, or generics. Instead, you get a language that feels like Go, uses standard Go tools for type-checking, and compiles to C code you could maintain by hand. C interop is first-class: if you declare a function without a body, it's treated as an extern; if you mark type as extern, it comes from your own headers. CGO is not used — Soan provides zero-cost interop with C.

Soan is for people who want Go's syntax and ergonomics for the kind of programs C is good at.

## Example

This Go code:

```go
package main

type Person struct {
    Name string
    Age  int
}

func (p *Person) Sleep() int {
    p.Age += 1
    return p.Age
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    p.Sleep()
    p.Sleep()
    p.Sleep()
    println(p.Name, "is now", p.Age, "years old.")
}
```

Translates to a header file `main.h`:

```c
#include "so.h"

typedef struct main_Person {
    so_String Name;
    so_int Age;
} main_Person;
so_int main_Person_Sleep(void* self);
```

Plus an implementation file `main.c`:

```c
#include "main.h"

so_int main_Person_Sleep(void* self) {
    main_Person* p = (main_Person*)self;
    p->Age += 1;
    return p->Age;
}

int main(void) {
    main_Person p = {.Name = so_strlit("Alice"), .Age = 30};
    main_Person_Sleep(&p);
    main_Person_Sleep(&p);
    main_Person_Sleep(&p);
    so_println("%s %s %lld %s", p.Name.ptr, "is now", p.Age, "years old.");
}
```

To learn more about the supported language features, see [language.md].

## Testing

Soan doesn't have its own testing framework. Since Soan code is valid Go code, you can just use `go test` like you normally would. Plus, your tests can use all Go features because they're never transpiled.

The transpilation logic is covered by the Soan compiler's own tests.
