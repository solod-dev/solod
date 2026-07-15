package compiler

import (
	"fmt"
	"strings"
)

// benchKind discovers BenchmarkXxx(b *testing.B) functions
// and runs them via testing.RunBenchmarks.
var benchKind = kind{
	subdir:  "bench",
	command: "so bench",
	noun:    "benchmark",
	prefix:  "Benchmark",
	param:   "B",
	emit:    emitBenchmarks,
}

// Bench discovers BenchmarkXxx functions in the "bench" subdirectory of srcDir,
// generates a deterministic main.go runner there, and runs it via Run.
func Bench(srcDir string, opts Options) error {
	return benchKind.run(srcDir, nil, opts)
}

// emitBenchmarks writes the runner body dispatching the benchmarks via
// testing.RunBenchmarks. Benchmarks always use the system allocator; a package
// that needs a different one can maintain its own main.go and use `so run`.
func emitBenchmarks(b *strings.Builder, pkg string, names []string) {
	b.WriteString("import (\n")
	b.WriteString("\t\"solod.dev/so/mem\"\n")
	b.WriteString("\t\"solod.dev/so/testing\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func main() {\n")
	fmt.Fprintf(b, "\ttesting.RunBenchmarks(mem.System, %q, []testing.Benchmark{\n", pkg)
	for _, name := range names {
		fmt.Fprintf(b, "\t\t{Name: %q, F: %s},\n", name, name)
	}
	b.WriteString("\t})\n")
	b.WriteString("}\n")
}
