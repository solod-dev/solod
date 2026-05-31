package main

import (
	"solod.dev/so/encoding/hex"
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
	"solod.dev/so/slices"
)

func main() {
	{
		// Encode.
		src := []byte("Hello Gopher!")
		dst := slices.Make[byte](mem.System, hex.EncodedLen(len(src)))
		hex.Encode(dst, src)
		if string(dst) != "48656c6c6f20476f7068657221" {
			panic("unexpected Encode result")
		}
		mem.FreeSlice(mem.System, dst)
	}
	{
		// EncodeToString.
		src := []byte("Hello Gopher!")
		encoded := hex.EncodeToString(mem.System, src)
		if encoded != "48656c6c6f20476f7068657221" {
			panic("unexpected EncodeToString result")
		}
		mem.FreeString(mem.System, encoded)
	}
	{
		// Decode.
		src := []byte("48656c6c6f20476f7068657221")
		dst := slices.Make[byte](mem.System, hex.DecodedLen(len(src)))
		n, err := hex.Decode(dst, src)
		if err != nil {
			panic(err)
		}
		if string(dst[:n]) != "Hello Gopher!" {
			panic("unexpected Decode result")
		}
		mem.FreeSlice(mem.System, dst)
	}
	{
		// DecodeString.
		const s = "48656c6c6f20476f7068657221"
		decoded, err := hex.DecodeString(mem.System, s)
		if err != nil {
			panic(err)
		}
		if string(decoded) != "Hello Gopher!" {
			panic("unexpected DecodeString result")
		}
		mem.FreeSlice(mem.System, decoded)
	}
	{
		// Dump.
		content := []byte("Go is an open source programming language.")
		dmp := hex.Dump(mem.System, content)
		fmt.Printf("%s", dmp)
		mem.FreeString(mem.System, dmp)
	}
}
