package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
	"solod.dev/so/strings"
	"solod.dev/so/testing"
)

var arena *mem.Arena

//so:volatile
var sink string

//so:volatile
var sinkInt int

func Clone(b *testing.B) {
	a := b.Allocator()
	str := strings.Repeat(nil, "a", 1024)
	defer mem.FreeString(nil, str)
	for b.Loop() {
		s := strings.Clone(a, str)
		mem.FreeString(a, s)
		if arena != nil {
			arena.Reset()
		}
	}
}

func Compare(b *testing.B) {
	str1 := strings.Repeat(nil, "01234567890αβγδεζ", 64)
	str2 := strings.Repeat(nil, "01234567890αβγδεζ", 64)
	defer mem.FreeString(nil, str1)
	defer mem.FreeString(nil, str2)
	for b.Loop() {
		sinkInt = strings.Compare(str1, str2)
		if arena != nil {
			arena.Reset()
		}
	}
}

func Fields(b *testing.B) {
	a := b.Allocator()
	str := strings.Repeat(nil, "01234567890αβ γδεζ", 16)
	defer mem.FreeString(nil, str)
	for b.Loop() {
		fields := strings.Fields(a, str)
		mem.FreeSlice(a, fields)
		if arena != nil {
			arena.Reset()
		}
	}
}

func Index(b *testing.B) {
	sb := strings.NewBuilder(nil)
	for range 64 {
		sb.WriteString("01234567890αβγδεζ")
	}
	sb.WriteRune('ω')
	str := sb.String() // 1025 chars, search for last
	defer mem.FreeString(nil, str)
	for b.Loop() {
		sinkInt = strings.Index(str, "ω")
		if arena != nil {
			arena.Reset()
		}
	}
}

func IndexByte(b *testing.B) {
	sb := strings.NewBuilder(nil)
	for range 64 {
		sb.WriteString("01234567890αβγδεζ")
	}
	sb.WriteByte('X')
	str := sb.String() // 1025 chars, search for last
	defer mem.FreeString(nil, str)
	for b.Loop() {
		sinkInt = strings.Index(str, "X")
		if arena != nil {
			arena.Reset()
		}
	}
}

func Repeat(b *testing.B) {
	a := b.Allocator()
	str := "0123456789abcdef"
	for b.Loop() {
		s := strings.Repeat(a, str, 64)
		mem.FreeString(a, s)
		if arena != nil {
			arena.Reset()
		}
	}
}

func ReplaceAll(b *testing.B) {
	a := b.Allocator()
	str := strings.Repeat(nil, "0123456789abcdef", 16)
	defer mem.FreeString(nil, str)
	for b.Loop() {
		s := strings.ReplaceAll(a, str, "a", "A")
		mem.FreeString(a, s)
		if arena != nil {
			arena.Reset()
		}
	}
}

func Split(b *testing.B) {
	a := b.Allocator()
	str := strings.Repeat(nil, "01234567890αβγδεζ", 16)
	defer mem.FreeString(nil, str)
	for b.Loop() {
		fields := strings.Split(a, str, "γ")
		mem.FreeSlice(a, fields)
		if arena != nil {
			arena.Reset()
		}
	}
}

func Trim(b *testing.B) {
	sb := strings.NewBuilder(nil)
	sb.WriteString("ηθικλμνξοπρστυφχψω")
	for range 64 {
		sb.WriteString("01234567890αβγδεζ")
	}
	sb.WriteString("ηθικλμνξοπρστυφχψω")
	str := sb.String()
	defer mem.FreeString(nil, str)
	for b.Loop() {
		sink = strings.Trim(str, "ωψχφυτσρποξνμλκιθη")
	}
}

func TrimSuffix(b *testing.B) {
	str := strings.Repeat(nil, "01234567890αβγδεζ", 16)
	suffix := "01234567890αβγδεζ"
	defer mem.FreeString(nil, str)
	for b.Loop() {
		sink = strings.TrimSuffix(str, suffix)
	}
}

func ToUpper(b *testing.B) {
	a := b.Allocator()
	str := strings.Repeat(nil, "01234567890αβγδεζ", 16)
	defer mem.FreeString(nil, str)
	for b.Loop() {
		s := strings.ToUpper(a, str)
		mem.FreeString(a, s)
		if arena != nil {
			arena.Reset()
		}
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "Clone", F: Clone},
		{Name: "Compare", F: Compare},
		{Name: "Fields", F: Fields},
		{Name: "Index", F: Index},
		{Name: "IndexByte", F: IndexByte},
		{Name: "Repeat", F: Repeat},
		{Name: "ReplaceAll", F: ReplaceAll},
		{Name: "Split", F: Split},
		{Name: "ToUpper", F: ToUpper},
		{Name: "Trim", F: Trim},
		{Name: "TrimSuffix", F: TrimSuffix},
	}

	fmt.Println("Malloc-based allocator:")
	testing.RunBenchmarks(mem.System, benchs)

	fmt.Println("Arena allocator:")
	var buf [4096]byte
	a := mem.NewArena(buf[:])
	arena = &a
	testing.RunBenchmarks(arena, benchs)
}
