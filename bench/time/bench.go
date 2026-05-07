package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/testing"
	"solod.dev/so/time"
)

//so:volatile
var sinkInt int64

//so:volatile
var sinkStr string

//so:volatile
var sinkTime time.Time

func Date(b *testing.B) {
	t := time.Now()
	for b.Loop() {
		date := t.Date(time.UTC)
		sinkInt = int64(date.Year)
	}
}

func Format(b *testing.B) {
	buf := make([]byte, 64)
	t := time.Unix(1265346057, 0)
	for b.Loop() {
		sinkStr = t.Format(buf, time.RFC3339, time.UTC)
	}
}

func FormatCustom(b *testing.B) {
	buf := make([]byte, 64)
	t := time.Unix(1265346057, 0)
	for b.Loop() {
		sinkStr = t.Format(buf, "%d.%m.%Y", time.UTC)
	}
}

func ISOWeek(b *testing.B) {
	t := time.Now()
	for b.Loop() {
		_, week := t.ISOWeek()
		sinkInt = int64(week)
	}
}

func Now(b *testing.B) {
	for b.Loop() {
		sinkTime = time.Now()
	}
}

func Parse(b *testing.B) {
	str := "2011-11-18T15:56:35Z"
	_, err := time.Parse(time.RFC3339, str, time.UTC)
	if err != nil {
		panic(err)
	}
	for b.Loop() {
		sinkTime, _ = time.Parse(time.RFC3339, str, time.UTC)
	}
}

func ParseCustom(b *testing.B) {
	str := "15.03.2024"
	_, err := time.Parse("%d.%m.%Y", str, time.UTC)
	if err != nil {
		panic(err)
	}
	for b.Loop() {
		sinkTime, _ = time.Parse("%d.%m.%Y", str, time.UTC)
	}
}

func Since(b *testing.B) {
	start := time.Now()
	for b.Loop() {
		sinkInt = int64(time.Since(start))
	}
}

func UnixNano(b *testing.B) {
	for b.Loop() {
		sinkInt = time.Now().UnixNano()
	}
}

func Until(b *testing.B) {
	end := time.Now().Add(1 * time.Hour)
	for b.Loop() {
		sinkInt = int64(time.Until(end))
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "Date", F: Date},
		{Name: "Format", F: Format},
		{Name: "FormatCustom", F: FormatCustom},
		{Name: "ISOWeek", F: ISOWeek},
		{Name: "Now", F: Now},
		{Name: "Parse", F: Parse},
		{Name: "ParseCustom", F: ParseCustom},
		{Name: "Since", F: Since},
		{Name: "UnixNano", F: UnixNano},
		{Name: "Until", F: Until},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
