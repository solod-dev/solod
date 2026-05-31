// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"testing"

	"solod.dev/so/encoding/binary"
	"solod.dev/so/time"
)

func TestNew(t *testing.T) {
	for _, test := range []struct {
		name    string
		newf    func() UUID
		version byte
		variant byte
	}{{
		name:    "New",
		newf:    New,
		version: 4,
		variant: 0b10,
	}, {
		name:    "NewV4",
		newf:    NewV4,
		version: 4,
		variant: 0b10,
	}, {
		name:    "NewV7",
		newf:    NewV7,
		version: 7,
		variant: 0b10,
	}} {
		u := test.newf()
		if got, want := version(u), test.version; got != want {
			t.Errorf("%v: got version %v, want %v", test.name, got, want)
		}
		if got, want := variant(u), test.variant; got != want {
			t.Errorf("%v: got variant %v, want %v", test.name, got, want)
		}
	}
}

func TestNewV7Millis(t *testing.T) {
	u := NewV7()
	got := binary.BigEndian.Uint64(u.Value[:8]) >> 16
	want := uint64(time.Now().UnixMilli())
	if got != want {
		t.Errorf("at %v, NewV7() = %v; millis = %x, want %x", time.Now(), u, got, want)
	}
}

func TestEncode(t *testing.T) {
	val := [16]byte{0xf8, 0x1d, 0x4f, 0xae, 0x7d, 0xec, 0x11, 0xd0, 0xa7, 0x65, 0x00, 0xa0, 0xc9, 0x1e, 0x6b, 0xf6}
	u := UUID{val}
	want := "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	buf := make([]byte, UUIDLen)
	if got := u.String(buf); got != want {
		t.Errorf("u.String() = %q, want %q", got, want)
	}
	if got, err := u.MarshalText(buf); string(got) != want || err != nil {
		t.Errorf("u.MarshalText() = %q, %v; want %q, nil", got, err, want)
	}
	prefix := []byte("urn:uuid:")
	if got, err := u.AppendText(prefix); string(got) != string(prefix)+want || err != nil {
		t.Errorf("u.MarshalAppend(%q) = %q, %v; want %q, nil", prefix, got, err, string(prefix)+want)
	}
}

func TestUnmarshalText(t *testing.T) {
	var got UUID
	err := got.UnmarshalText([]byte("f81d4fae-7dec-11d0-a765-00a0c91e6bf6"))
	if err != nil {
		t.Errorf("UnmarshalText: %v", err)
	}
	val := [16]byte{0xf8, 0x1d, 0x4f, 0xae, 0x7d, 0xec, 0x11, 0xd0, 0xa7, 0x65, 0x00, 0xa0, 0xc9, 0x1e, 0x6b, 0xf6}
	want := UUID{val}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseSuccess(t *testing.T) {
	val := [16]byte{
		0xf8, 0x1d, 0x4f, 0xae,
		0x7d, 0xec,
		0x11, 0xd0,
		0xa7, 0x65,
		0x00, 0xa0, 0xc9, 0x1e, 0x6b, 0xf6,
	}
	u1 := UUID{val}
	for _, test := range []struct {
		s string
		u UUID
	}{{
		s: "00000000-0000-0000-0000-000000000000",
		u: Nil(),
	}, {
		s: "ffffffff-ffff-ffff-ffff-ffffffffffff",
		u: Max(),
	}, {
		s: "f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
		u: u1,
	}, {
		s: "F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6",
		u: u1,
	}, {
		s: "f81d4fae7dec11d0a76500a0c91e6bf6",
		u: u1,
	}, {
		s: "{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}",
		u: u1,
	}, {
		s: "urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
		u: u1,
	}} {
		u, err := Parse(test.s)
		if err != nil {
			t.Errorf("Parse(%q) = _, %v; want success", test.s, err)
		} else if u != test.u {
			t.Errorf("Parse(%q) = %v, nil; want %v", test.s, u, test.u)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for _, test := range []string{
		"",
		"0000000000000-0000-0000-000000000000",
		"00000000-000000000-0000-000000000000",
		"00000000-0000-000000000-000000000000",
		"00000000-0000-0000-00000000000000000",
		"00000000-0000-0000-0000-00000000000",
		"x0000000-0000-0000-0000-000000000000",
		"00000000-x000-0000-0000-000000000000",
		"00000000-0000-x000-0000-000000000000",
		"00000000-0000-0000-x000-000000000000",
		"00000000-0000-0000-0000-x00000000000",
		"{x0000000-0000-0000-0000-000000000000}",
		"urn:uuid:x000000-0000-0000-0000-000000000000",
		"x0000000000000000000000000000000",
		// Some parsers permit hyphens in non-standard locations,
		// but we currently do not.
		"0000-0000-0000-0000-0000-0000-0000-0000",
		// Combinations of variant encodings that we could choose to parse,
		// but currently do not.
		"{00000000000000000000000000000000}",
		"{urn:uuid:00000000-0000-0000-0000-000000000000}",
		"urn:uuid:00000000000000000000000000000000",
	} {
		got, err := Parse(test)
		if err == nil {
			t.Errorf("Parse(%q) = %v, nil; want error", test, got)
		}
	}
}

func TestCompare(t *testing.T) {
	uuids := []UUID{
		Nil(),
		MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6"),
		Max(),
	}
	for i, u := range uuids {
		if got, want := u.Compare(u), 0; got != want {
			t.Errorf("%v.Compare(itself) = %v, want %v", u, got, want)
		}
		if i == 0 {
			continue
		}
		prev := uuids[i-1]
		if got, want := u.Compare(prev), 1; got != want {
			t.Errorf("%v.Compare(%v) = %v, want %v", u, prev, got, want)
		}
		if got, want := prev.Compare(u), -1; got != want {
			t.Errorf("%v.Compare(%v) = %v, want %v", prev, u, got, want)
		}
	}
}

func version(u UUID) byte {
	return u.Value[6] >> 4
}

func variant(u UUID) byte {
	return u.Value[8] >> 6
}
