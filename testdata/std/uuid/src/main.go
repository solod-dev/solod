package main

import "solod.dev/so/uuid"

func main() {
	const ustr = "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	{
		// NewV4 and NewV7.
		u4 := uuid.NewV4()
		if u4.Version() != 4 {
			panic("NewV4() version != 4")
		}
		u7 := uuid.NewV7()
		if u7.Version() != 7 {
			panic("NewV7() version != 7")
		}
	}
	{
		// String and Parse.
		u1 := uuid.MustParse(ustr)
		buf := make([]byte, uuid.UUIDLen)
		s := u1.String(buf)
		if s != ustr {
			panic("String() mismatch")
		}
		u2, err := uuid.Parse(s)
		if err != nil {
			panic(err)
		}
		if u1 != u2 {
			panic("Parse/String mismatch")
		}
	}
	{
		// Compare.
		unil := uuid.Nil()
		uid := uuid.MustParse(ustr)
		umax := uuid.Max()
		if uid.Compare(unil) <= 0 {
			panic("Compare: uid <= unil")
		}
		if uid.Compare(umax) >= 0 {
			panic("Compare: uid >= umax")
		}
		if uid.Compare(uid) != 0 {
			panic("Compare: uid != uid")
		}
	}
}
