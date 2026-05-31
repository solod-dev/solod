// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hex

import (
	"testing"

	"solod.dev/so/bytes"
	"solod.dev/so/io"
	"solod.dev/so/strings"
)

type encDecTest struct {
	enc string
	dec []byte
}

var encDecTests = []encDecTest{
	{"", []byte{}},
	{"0001020304050607", []byte{0, 1, 2, 3, 4, 5, 6, 7}},
	{"08090a0b0c0d0e0f", []byte{8, 9, 10, 11, 12, 13, 14, 15}},
	{"f0f1f2f3f4f5f6f7", []byte{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7}},
	{"f8f9fafbfcfdfeff", []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}},
	{"67", []byte{'g'}},
	{"e3a1", []byte{0xe3, 0xa1}},
}

func TestEncode(t *testing.T) {
	for i, test := range encDecTests {
		dst := make([]byte, EncodedLen(len(test.dec)))
		n := Encode(dst, test.dec)
		if n != len(dst) {
			t.Errorf("#%d: bad return value: got: %d want: %d", i, n, len(dst))
		}
		if string(dst) != test.enc {
			t.Errorf("#%d: got: %#v want: %#v", i, dst, test.enc)
		}
		dst = []byte("lead")
		dst = AppendEncode(nil, dst, test.dec)
		if string(dst[4:]) != test.enc {
			t.Errorf("#%d: got: %#v want: %#v", i, dst, test.enc)
		}
	}
}

func TestDecode(t *testing.T) {
	// Case for decoding uppercase hex characters, since
	// Encode always uses lowercase.
	decTests := append(encDecTests, encDecTest{"F8F9FAFBFCFDFEFF", []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}})
	for i, test := range decTests {
		dst := make([]byte, DecodedLen(len(test.enc)))
		n, err := Decode(dst, []byte(test.enc))
		if err != nil {
			t.Errorf("#%d: bad return value: got:%d want:%d", i, n, len(dst))
		} else if !bytes.Equal(dst, test.dec) {
			t.Errorf("#%d: got: %#v want: %#v", i, dst, test.dec)
		}
		dst = []byte("lead")
		dst, err = AppendDecode(nil, dst, []byte(test.enc))
		if err != nil {
			t.Errorf("#%d: AppendDecode error: %v", i, err)
		} else if string(dst[4:]) != string(test.dec) {
			t.Errorf("#%d: got: %#v want: %#v", i, dst, test.dec)
		}
	}
}

func TestEncodeToString(t *testing.T) {
	for i, test := range encDecTests {
		s := EncodeToString(nil, test.dec)
		if s != test.enc {
			t.Errorf("#%d got:%s want:%s", i, s, test.enc)
		}
	}
}

func TestDecodeString(t *testing.T) {
	for i, test := range encDecTests {
		dst, err := DecodeString(nil, test.enc)
		if err != nil {
			t.Errorf("#%d: unexpected err value: %s", i, err)
			continue
		}
		if !bytes.Equal(dst, test.dec) {
			t.Errorf("#%d: got: %#v want: #%v", i, dst, test.dec)
		}
	}
}

var errTests = []struct {
	in  string
	out string
	err error
}{
	{"", "", nil},
	{"0", "", ErrLength},
	{"zd4aa", "", ErrInvalidByte},
	{"d4aaz", "\xd4\xaa", ErrInvalidByte},
	{"30313", "01", ErrLength},
	{"0g", "", ErrInvalidByte},
	{"00gg", "\x00", ErrInvalidByte},
	{"0\x01", "", ErrInvalidByte},
	{"ffeed", "\xff\xee", ErrLength},
}

func TestDecodeErr(t *testing.T) {
	for _, tt := range errTests {
		out := make([]byte, len(tt.in)+10)
		n, err := Decode(out, []byte(tt.in))
		if string(out[:n]) != tt.out || err != tt.err {
			t.Errorf("Decode(%q) = %q, %v, want %q, %v", tt.in, string(out[:n]), err, tt.out, tt.err)
		}
	}
}

func TestDecodeStringErr(t *testing.T) {
	for _, tt := range errTests {
		out, err := DecodeString(nil, tt.in)
		if string(out) != tt.out || err != tt.err {
			t.Errorf("DecodeString(nil, %q) = %q, %v, want %q, %v", tt.in, out, err, tt.out, tt.err)
		}
	}
}

func TestEncoderDecoder(t *testing.T) {
	copyBuf := make([]byte, 7)
	for _, multiplier := range []int{1, 128, 192} {
		for _, test := range encDecTests {
			input := bytes.Repeat(nil, test.dec, multiplier)
			output := strings.Repeat(nil, test.enc, multiplier)

			var buf bytes.Buffer
			enc := NewEncoder(&buf)
			r := bytes.NewReader(input)
			if n, err := io.CopyBuffer(&enc, &r, copyBuf); n != int64(len(input)) || err != nil {
				t.Errorf("encoder.Write(%q*%d) = (%d, %v), want (%d, nil)", test.dec, multiplier, n, err, len(input))
				continue
			}

			if encDst := buf.String(); encDst != output {
				t.Errorf("buf(%q*%d) = %v, want %v", test.dec, multiplier, encDst, output)
				continue
			}

			dec := NewDecoder(&buf)
			var decBuf bytes.Buffer
			w := &decBuf
			if _, err := io.CopyBuffer(w, &dec, copyBuf); err != nil || decBuf.Len() != len(input) {
				t.Errorf("decoder.Read(%q*%d) = (%d, %v), want (%d, nil)", test.enc, multiplier, decBuf.Len(), err, len(input))
			}

			if !bytes.Equal(decBuf.Bytes(), input) {
				t.Errorf("decBuf(%q*%d) = %v, want %v", test.dec, multiplier, decBuf.Bytes(), input)
				continue
			}
		}
	}
}

func TestDecoderErr(t *testing.T) {
	for _, tt := range errTests {
		r := strings.NewReader(tt.in)
		dec := NewDecoder(&r)
		out, err := io.ReadAll(nil, &dec)
		wantErr := tt.err
		// Decoder is reading from stream, so it reports io.ErrUnexpectedEOF instead of ErrLength.
		if wantErr == ErrLength {
			wantErr = io.ErrUnexpectedEOF
		}
		if string(out) != tt.out || err != wantErr {
			t.Errorf("NewDecoder(%q) = %q, %v, want %q, %v", tt.in, out, err, tt.out, wantErr)
		}
	}
}

var expectedHexDump = []byte(`00000000  1e 1f 20 21 22 23 24 25  26 27 28 29 2a 2b 2c 2d  |.. !"#$%&'()*+,-|
00000010  2e 2f 30 31 32 33 34 35  36 37 38 39 3a 3b 3c 3d  |./0123456789:;<=|
00000020  3e 3f 40 41 42 43 44 45                           |>?@ABCDE|
`)

func TestDumper(t *testing.T) {
	var in [40]byte
	for i := range in {
		in[i] = byte(i + 30)
	}

	for stride := 1; stride < len(in); stride++ {
		var out bytes.Buffer
		dumper := NewDumper(&out)
		done := 0
		for done < len(in) {
			todo := done + stride
			if todo > len(in) {
				todo = len(in)
			}
			dumper.Write(in[done:todo])
			done = todo
		}

		dumper.Close()
		if !bytes.Equal(out.Bytes(), expectedHexDump) {
			t.Errorf("stride: %d failed. got:\n%s\nwant:\n%s", stride, out.Bytes(), expectedHexDump)
		}
	}
}

func TestDumper_doubleclose(t *testing.T) {
	var out strings.Builder
	dumper := NewDumper(&out)

	dumper.Write([]byte(`gopher`))
	dumper.Close()
	dumper.Close()
	dumper.Write([]byte(`gopher`))
	dumper.Close()

	expected := "00000000  67 6f 70 68 65 72                                 |gopher|\n"
	if out.String() != expected {
		t.Fatalf("got:\n%#v\nwant:\n%#v", out.String(), expected)
	}
}

func TestDumper_earlyclose(t *testing.T) {
	var out strings.Builder
	dumper := NewDumper(&out)

	dumper.Close()
	dumper.Write([]byte(`gopher`))

	expected := ""
	if out.String() != expected {
		t.Fatalf("got:\n%#v\nwant:\n%#v", out.String(), expected)
	}
}

func TestDump(t *testing.T) {
	var in [40]byte
	for i := range in {
		in[i] = byte(i + 30)
	}

	out := []byte(Dump(nil, in[:]))
	if !bytes.Equal(out, expectedHexDump) {
		t.Errorf("got:\n%s\nwant:\n%s", out, expectedHexDump)
	}
}
