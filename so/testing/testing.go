package testing

import (
	"solod.dev/so/flag"
	"solod.dev/so/fmt"
	"solod.dev/so/io"
	"solod.dev/so/mem"
	"solod.dev/so/os"
	"solod.dev/so/strings"
)

// The C backing for the variadic Errorf/Fatalf methods, which cannot be
// expressed in So (a So variadic packs its args into a slice; a real C variadic
// is needed to forward them to fmt).
//
//so:embed testing.h
var testing_h string

//so:embed testing.c
var testing_c string

// T is the context passed to a test function. It records failure and skip
// state for a single test.
//
// The plain message methods (Log, Error, Fatal, Skip) take a preformatted
// string. For formatted messages use the variadic [T.Errorf] and [T.Fatalf]:
//
//	t.Errorf("Index = %d, want 6", got)
//
// So also has no recover, so T cannot unwind a running test. Fatal only marks
// the test failed and prints the message; by convention the test function must
// return right after calling it:
//
//	if err != nil {
//		t.Fatal("open failed")
//		return
//	}
//
// By design, a hard crash (panic or segfault) in a test aborts
// the entire test run, not just the current test.
type T struct {
	name    string
	w       io.Writer
	failed  bool
	skipped bool

	alloc mem.Tracker
}

// Name returns the name of the running test.
func (t *T) Name() string { return t.name }

// Allocator returns the memory allocator for the test. Allocations made
// through it are tracked, and after the test function returns the runner fails
// the test if any of them were not freed. Use it in place of [mem.System] to
// enable leak checking:
//
//	alloc := t.Allocator()
//	p := mem.Alloc[int](alloc)
//	defer mem.Free(alloc, p)
//
// Allocations made through any other allocator are not tracked.
func (t *T) Allocator() mem.Allocator { return &t.alloc }

// Failed reports whether the test has failed.
func (t *T) Failed() bool { return t.failed }

// Fail marks the test as failed but continues execution.
func (t *T) Fail() { t.failed = true }

// Log records msg in the test log.
func (t *T) Log(msg string) {
	fmt.Fprintf(t.w, "    %s\n", msg)
}

// Error is equivalent to Log followed by Fail.
func (t *T) Error(msg string) {
	t.Log(msg)
	t.Fail()
}

// Errorf formats its arguments like [fmt.Sprintf], then
// behaves like [T.Error] (Log followed by Fail).
//
//so:extern
func (t *T) Errorf(format string, args ...any) {
	buf := fmt.NewBuffer(fmt.BufSize)
	t.Error(fmt.Sprintf(buf, format, args...))
}

// Fatal is equivalent to Log followed by Fail. The test function must return
// after calling it; see [T].
func (t *T) Fatal(msg string) {
	t.Log(msg)
	t.Fail()
}

// Fatalf formats its arguments like [fmt.Sprintf], then behaves like [T.Fatal].
// The test function must return after calling it; see [T].
//
//so:extern
func (t *T) Fatalf(format string, args ...any) {
	buf := fmt.NewBuffer(fmt.BufSize)
	t.Fatal(fmt.Sprintf(buf, format, args...))
}

// Skip marks the test as skipped. Like Fatal, the test must return afterwards.
func (t *T) Skip(msg string) {
	t.Log(msg)
	t.skipped = true
}

// Test represents a single test to be run by the test runner.
type Test struct {
	Name string
	F    func(t *T)
}

// RunTests runs the given tests for package pkg, prints per-test results
// to stdout, and exits with a non-zero status if any test failed.
// args is the runner's os.Args; RunTests parses flags from it.
func RunTests(pkg string, args []string, tests []Test) {
	var run string
	fs := flag.NewFlagSet("so test", flag.ContinueOnError)
	fs.StringVar(&run, "run", "", "run only tests whose names start with this prefix")
	if err := fs.Parse(args[1:]); err != nil {
		os.Exit(2)
	}

	failed := 0
	skipped := 0
	total := 0
	for _, tc := range tests {
		if !strings.HasPrefix(tc.Name, run) {
			continue
		}
		total++

		t := &T{name: tc.Name, w: os.Stdout}
		t.alloc.Allocator = mem.System
		fmt.Fprintf(t.w, "=== RUN   %s\n", t.name)
		tc.F(t)

		// Fail a passing test that leaked memory allocated through t.Allocator().
		if !t.failed && !t.skipped {
			stats := t.alloc.Stats()
			if stats.Mallocs != stats.Frees {
				fmt.Fprintf(t.w, "    memory leak: %d unfreed allocation(s), %d byte(s)\n",
					stats.Mallocs-stats.Frees, stats.Alloc)
				t.failed = true
			}
		}

		if t.skipped {
			fmt.Fprintf(t.w, "--- SKIP: %s\n", t.name)
			skipped++
			continue
		}
		if t.failed {
			fmt.Fprintf(t.w, "--- FAIL: %s\n", t.name)
			failed++
			continue
		}
		fmt.Fprintf(t.w, "--- PASS: %s\n", t.name)
	}

	if total == 0 {
		fmt.Fprintf(os.Stdout, "ok\t%s\t%d tests [no tests to run]\n", pkg, total)
		return
	}
	if failed > 0 {
		fmt.Fprintf(os.Stdout, "FAIL\t%s\t%d of %d failed\n", pkg, failed, total)
		os.Exit(1)
	}
	if skipped > 0 {
		fmt.Fprintf(os.Stdout, "ok\t%s\t%d tests (%d skipped)\n", pkg, total, skipped)
		return
	}
	fmt.Fprintf(os.Stdout, "ok\t%s\t%d tests\n", pkg, total)
}
