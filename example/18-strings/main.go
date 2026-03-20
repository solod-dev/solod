// A So string is a read-only slice of bytes. The language
// and the standard library treat strings specially - as
// containers of text encoded in UTF-8.
// In other languages, strings are made of "characters".
// In So, the concept of a character is called a `rune` - it's
// an integer that represents a Unicode code point.
package main

import "solod.dev/so/c/stdio"

func main() {
	// `s` is a `string` assigned a literal value
	// representing the word "hello" in the Thai
	// language. So string literals are UTF-8
	// encoded text.
	const s = "สวัสดี"

	// Since strings are equivalent to `[]byte`, this
	// will produce the length of the raw bytes stored within.
	println("Len:", len(s))

	// Indexing into a string produces the raw byte values at
	// each index. This loop generates the hex values of all
	// the bytes that constitute the code points in `s`.
	for i := 0; i < len(s); i++ {
		stdio.Printf("%x ", s[i])
	}
	println()

	// To count how many _runes_ are in a string, we can use
	// convert the string to a rune slice. Note that this conversion
	// is an O(n) operation in both time and memory,
	// because it has to decode each UTF-8 rune sequentially.
	// Some Thai characters are represented by UTF-8 code points
	// that can span multiple bytes, so the result of this count
	// may be surprising.
	runes := []rune(s)
	println("Rune count:", len(runes))

	// A `range` loop handles strings specially and decodes
	// each `rune` along with its offset in the string.
	for idx, runeValue := range s {
		stdio.Printf("0x%x starts at %d\n", runeValue, int32(idx))
		// This demonstrates passing a `rune` value to a function.
		examineRune(runeValue)
	}
}

func examineRune(r rune) {
	// Values enclosed in single quotes are _rune literals_. We
	// can compare a `rune` value to a rune literal directly.
	if r == 't' {
		println("found tee")
	} else if r == 'ส' {
		println("found so sua")
	}
}
