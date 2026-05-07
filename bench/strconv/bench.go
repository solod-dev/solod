package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
	"solod.dev/so/strconv"
	"solod.dev/so/testing"
)

//so:volatile
var sinkInt int

//so:volatile
var sinkFloat float64

func Atof64_Decimal(b *testing.B) {
	for b.Loop() {
		f, _ := strconv.ParseFloat("33909", 64)
		sinkFloat += f
	}
}

func Atof64_Float(b *testing.B) {
	for b.Loop() {
		f, _ := strconv.ParseFloat("339.7784", 64)
		sinkFloat += f
	}
}

func Atof64_Exp(b *testing.B) {
	for b.Loop() {
		f, _ := strconv.ParseFloat("-5.09e75", 64)
		sinkFloat += f
	}
}

func Atof64_Big(b *testing.B) {
	for b.Loop() {
		f, _ := strconv.ParseFloat("1844674407370955", 64)
		sinkFloat += f
	}
}

func ParseInt_7bit(b *testing.B) {
	buf := fmt.NewBuffer(64)
	s := fmt.Sprintf(buf, "%d", 1<<7-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func ParseInt_26bit(b *testing.B) {
	buf := fmt.NewBuffer(64)
	s := fmt.Sprintf(buf, "%d", 1<<26-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func ParseInt_31bit(b *testing.B) {
	buf := fmt.NewBuffer(64)
	s := fmt.Sprintf(buf, "%d", 1<<31-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func ParseInt_56bit(b *testing.B) {
	buf := fmt.NewBuffer(64)
	format := "%lld"
	s := fmt.Sprintf(buf, format, 1<<56-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func ParseInt_62bit(b *testing.B) {
	buf := fmt.NewBuffer(64)
	format := "%lld"
	s := fmt.Sprintf(buf, format, 1<<62-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func FormatFloat_Decimal(b *testing.B) {
	buf := make([]byte, 64)
	for b.Loop() {
		s := strconv.FormatFloat(buf, 33909, 'f', -1, 64)
		sinkInt += len(s)
	}
}

func FormatFloat_Float(b *testing.B) {
	buf := make([]byte, 64)
	for b.Loop() {
		s := strconv.FormatFloat(buf, 339.7784, 'f', -1, 64)
		sinkInt += len(s)
	}
}

func FormatFloat_Exp(b *testing.B) {
	buf := make([]byte, 64)
	for b.Loop() {
		s := strconv.FormatFloat(buf, -5.09e75, 'e', -1, 64)
		sinkInt += len(s)
	}
}

func FormatFloat_Big(b *testing.B) {
	buf := make([]byte, 64)
	for b.Loop() {
		s := strconv.FormatFloat(buf, 1844674407370955, 'f', -1, 64)
		sinkInt += len(s)
	}
}

func FormatInt_7bit(b *testing.B) {
	buf := make([]byte, 64)
	n := 1<<7 - 1
	for b.Loop() {
		s := strconv.FormatInt(buf, int64(n), 10)
		sinkInt += len(s)
	}
}

func FormatInt_26bit(b *testing.B) {
	buf := make([]byte, 64)
	n := 1<<26 - 1
	for b.Loop() {
		s := strconv.FormatInt(buf, int64(n), 10)
		sinkInt += len(s)
	}
}

func FormatInt_31bit(b *testing.B) {
	buf := make([]byte, 64)
	n := 1<<31 - 1
	for b.Loop() {
		s := strconv.FormatInt(buf, int64(n), 10)
		sinkInt += len(s)
	}
}

func FormatInt_56bit(b *testing.B) {
	buf := make([]byte, 64)
	n := 1<<56 - 1
	for b.Loop() {
		s := strconv.FormatInt(buf, int64(n), 10)
		sinkInt += len(s)
	}
}

func FormatInt_62bit(b *testing.B) {
	buf := make([]byte, 64)
	n := 1<<62 - 1
	for b.Loop() {
		s := strconv.FormatInt(buf, int64(n), 10)
		sinkInt += len(s)
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "Atof64_Decimal", F: Atof64_Decimal},
		{Name: "Atof64_Float", F: Atof64_Float},
		{Name: "Atof64_Exp", F: Atof64_Exp},
		{Name: "Atof64_Big", F: Atof64_Big},
		{Name: "ParseInt_7bit", F: ParseInt_7bit},
		{Name: "ParseInt_26bit", F: ParseInt_26bit},
		{Name: "ParseInt_31bit", F: ParseInt_31bit},
		{Name: "ParseInt_56bit", F: ParseInt_56bit},
		{Name: "ParseInt_62bit", F: ParseInt_62bit},
		{Name: "FormatFloat_Decimal", F: FormatFloat_Decimal},
		{Name: "FormatFloat_Float", F: FormatFloat_Float},
		{Name: "FormatFloat_Exp", F: FormatFloat_Exp},
		{Name: "FormatFloat_Big", F: FormatFloat_Big},
		{Name: "FormatInt_7bit", F: FormatInt_7bit},
		{Name: "FormatInt_26bit", F: FormatInt_26bit},
		{Name: "FormatInt_31bit", F: FormatInt_31bit},
		{Name: "FormatInt_56bit", F: FormatInt_56bit},
		{Name: "FormatInt_62bit", F: FormatInt_62bit},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
