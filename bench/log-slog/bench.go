package main

import (
	"solod.dev/so/log/slog"
	"solod.dev/so/mem"
	"solod.dev/so/testing"
)

//so:volatile
var sink int

type soHandler struct{}

func (h *soHandler) Enabled(level slog.Level) bool {
	return true
}
func (h *soHandler) Handle(r slog.Record) error {
	sink += len(r.Message)
	return nil
}

func NoAttr(b *testing.B) {
	h := &soHandler{}
	l := slog.New(h)
	for b.Loop() {
		l.Log(slog.LevelInfo, "msg")
	}
}

func Attr(b *testing.B) {
	h := &soHandler{}
	l := slog.New(h)
	for b.Loop() {
		l.Log(slog.LevelInfo, "msg",
			slog.Int("a", 1), slog.String("b", "two"), slog.Bool("c", true))
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "NoAttr", F: NoAttr},
		{Name: "Attr", F: Attr},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
