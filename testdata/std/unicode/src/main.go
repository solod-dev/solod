package main

import "solod.dev/so/unicode"

func main() {
	{
		// Is.
		if !unicode.IsDigit('0') {
			panic("IsDigit failed")
		}
		if !unicode.IsLetter('a') {
			panic("IsLetter failed")
		}
		if !unicode.IsLower('a') {
			panic("IsLower failed")
		}
		if !unicode.IsSpace(' ') {
			panic("IsSpace failed")
		}
		if !unicode.IsTitle('ᾭ') {
			panic("IsTitle failed")
		}
		if !unicode.IsUpper('A') {
			panic("IsUpper failed")
		}
	}
	{
		// To.
		if unicode.ToLower('A') != 'a' {
			panic("ToLower failed")
		}
		if unicode.ToTitle('a') != 'A' {
			panic("ToTitle failed")
		}
		if unicode.ToUpper('a') != 'A' {
			panic("ToUpper failed")
		}
		if unicode.To(unicode.UpperCase, 'a') != 'A' {
			panic("To failed")
		}
	}
}
