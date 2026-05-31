// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uuid provides support for generating and manipulating UUIDs.
//
// See [RFC 9562] for details.
//
// Random components of new UUIDs are generated with a
// cryptographically secure random number generator.
//
// UUIDs may be generated using various algorithms.
// The [New] function returns a new UUID generated using
// an algorithm suitable for most purposes.
//
// [RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562.html
package uuid

import (
	"solod.dev/so/bytes"
	"solod.dev/so/crypto/crand"
	"solod.dev/so/encoding/binary"
	"solod.dev/so/encoding/hex"
	"solod.dev/so/errors"
	"solod.dev/so/time"
)

// Length of a canonical UUID string.
const UUIDLen = 36 // xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

var ErrInvalid = errors.New("uuid: invalid uuid")

// A UUID is a Universally Unique Identifier as specified in RFC 9562.
//
// UUIDs are comparable, such as with the == operator.
type UUID struct {
	Value [16]byte
}

// Parse returns the UUID represented by s.
//
// It accepts strings in the following forms:
//
//	f81d4fae-7dec-11d0-a765-00a0c91e6bf6
//	{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}
//	urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6
//	f81d4fae7dec11d0a76500a0c91e6bf6
//
// Alphabetic characters in the hexadecimal digits may be any case.
func Parse(s string) (UUID, error) {
	var u UUID
	err := u.UnmarshalText([]byte(s))
	return u, err
}

// MustParse returns the UUID represented by s.
//
// It panics if s is not a valid string representation of a UUID as defined by [Parse].
func MustParse(s string) UUID {
	u, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

// New returns a new UUID.
//
// Programs which do not have a need for a specific UUID generation algorithm should use New.
// At this time, New is equivalent to [NewV4].
func New() UUID {
	return NewV4()
}

// Nil returns the Nil UUID 00000000-0000-0000-0000-000000000000.
//
// The Nil UUID is defined in [Section 5.9 of RFC 9562].
// Note that this is not the same as the Go value nil.
//
// [Section 5.9 of RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562#section-5.9
func Nil() UUID {
	return UUID{}
}

// Max returns the Max UUID ffffffff-ffff-ffff-ffff-ffffffffffff.
//
// The Max UUID is defined in [Section 5.10 of RFC 9562].
//
// [Section 5.10 of RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562#section-5.10
func Max() UUID {
	return UUID{Value: [16]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}}
}

// Version returns the version number of u as defined in [Section 4.2 of RFC 9562].
//
// [Section 4.2 of RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562#section-4.2
func (u UUID) Version() int {
	return int(u.Value[6] >> 4)
}

// String returns the string representation of u.
// It uses the lowercase hex-and-dash representation defined in RFC 9562.
// buf length must be at least [UUIDLen] bytes.
func (u UUID) String(buf []byte) string {
	b, _ := u.MarshalText(buf)
	return string(b)
}

// MarshalText implements the [encoding.TextMarshaler] interface.
// The encoding is the same as returned by [UUID.String]
// buf length must be at least [UUIDLen] bytes.
func (u UUID) MarshalText(buf []byte) ([]byte, error) {
	return u.AppendText(buf[:0])
}

// AppendText implements the [encoding.TextAppender] interface.
// The encoding is the same as returned by [UUID.String].
// Requires at least [UUIDLen] bytes of spare capacity in buf.
// Always returns a nil error.
func (u UUID) AppendText(b []byte) ([]byte, error) {
	off := len(b)
	b = append(b, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"...)
	dst := b[off:]
	hex.Encode(dst[0:8], u.Value[0:4])
	hex.Encode(dst[9:13], u.Value[4:6])
	hex.Encode(dst[14:18], u.Value[6:8])
	hex.Encode(dst[19:23], u.Value[8:10])
	hex.Encode(dst[24:36], u.Value[10:16])
	return b, nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
// The UUID is expected in a form accepted by [Parse].
func (u *UUID) UnmarshalText(b []byte) error {
	var dst UUID
	switch len(b) {
	case len("urn:uuid:") + 36:
		// urn:uuid:00000000-0000-0000-0000-000000000000
		b = bytes.TrimPrefix(b, []byte("urn:uuid:"))
	case 2 + 36:
		// {00000000-0000-0000-0000-000000000000}
		b = bytes.TrimPrefix(b, []byte("{"))
		b = bytes.TrimSuffix(b, []byte("}"))
	case 32:
		// 00000000000000000000000000000000
		if _, err := hex.Decode(dst.Value[:], b); err != nil {
			return ErrInvalid
		}
		*u = dst
		return nil
	}
	// 00000000-0000-0000-0000-000000000000
	if len(b) != 36 {
		return ErrInvalid
	}
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return ErrInvalid
	}
	if _, err := hex.Decode(dst.Value[0:4], b[0:8]); err != nil {
		return ErrInvalid
	}
	if _, err := hex.Decode(dst.Value[4:6], b[9:13]); err != nil {
		return ErrInvalid
	}
	if _, err := hex.Decode(dst.Value[6:8], b[14:18]); err != nil {
		return ErrInvalid
	}
	if _, err := hex.Decode(dst.Value[8:10], b[19:23]); err != nil {
		return ErrInvalid
	}
	if _, err := hex.Decode(dst.Value[10:16], b[24:36]); err != nil {
		return ErrInvalid
	}
	*u = dst
	return nil
}

// Compare compares the UUID u with v.
// If u is before v, it returns -1.
// If u is after v, it returns +1.
// If they are the same, it returns 0.
//
// Compare uses the big-endian byte order defined in
// [Section 6.11 of RFC 9562] for sorting.
//
// [Section 6.11 of RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562#section-6.11
func (u UUID) Compare(v UUID) int {
	for i := range u.Value {
		ui := u.Value[i]
		vi := v.Value[i]
		if ui < vi {
			return -1
		}
		if ui > vi {
			return 1
		}
	}
	return 0
}

// NewV4 returns a new version 4 UUID.
//
// Version 4 UUIDs contain 122 bits of random data.
func NewV4() UUID {
	var u UUID
	crand.Read(u.Value[:])
	u.setVersion(4)
	u.setVariant(0b10)
	return u
}

// NewV7 returns a new version 7 UUID.
//
// Version 7 UUIDs contain a timestamp in the most significant 48 bits,
// and at least 62 bits of random data.
func NewV7() UUID {
	// UUIDv7 is defined in RFC 9562 section 5.7 as:
	//
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                           unix_ts_ms                          |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |          unix_ts_ms           |  ver  |       rand_a          |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |var|                        rand_b                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                            rand_b                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	// We store a 12 bit sub-millisecond timestamp fraction in the rand_a section,
	// as optionally permitted by the RFC.

	// Generate our 60-bit timestamp: 48 bits of millisecond-resolution,
	// followed by 12 bits of 1/4096-millisecond resolution.
	now := time.Now()
	secs := uint64(now.Unix())
	nanos := uint64(now.Nanosecond())
	msecs := nanos / 1000000
	frac := nanos - (1000000 * msecs)
	timestamp := (1000*secs + msecs) << 12 // ms shifted into position
	timestamp += (frac * 4096) / 1000000   // ns converted to 1/4096-ms units

	// Insert a gap for the 4 bits of the ver field into the timestamp.
	hibits := ((timestamp << 4) & 0xffff_ffff_ffff_0000) | (timestamp & 0x0ffff)

	var u UUID
	binary.BigEndian.PutUint64(u.Value[0:8], hibits)
	crand.Read(u.Value[8:])
	u.setVersion(7)
	u.setVariant(0b10)

	return u
}

func (u *UUID) setVersion(version byte) {
	u.Value[6] = (u.Value[6] & 0b0000_1111) | (version << 4)
}

func (u *UUID) setVariant(variant byte) {
	u.Value[8] = (u.Value[8] & 0b0011_1111) | (variant << 6)
}
