// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings_test

import (
	"fmt"
	"unsafe"

	"solod.dev/so/strings"
	"solod.dev/so/unicode"
)

func ExampleClone() {
	s := "abc"
	clone := strings.Clone(nil, s)
	fmt.Println(s == clone)
	fmt.Println(unsafe.StringData(s) == unsafe.StringData(clone))
	// Output:
	// true
	// false
}

func ExampleBuilder() {
	var b strings.Builder
	for i := 3; i >= 1; i-- {
		fmt.Fprintf(&b, "%d...", i)
	}
	b.WriteString("ignition")
	fmt.Println(b.String())

	// Output: 3...2...1...ignition
}

func ExampleCompare() {
	fmt.Println(strings.Compare("a", "b"))
	fmt.Println(strings.Compare("a", "a"))
	fmt.Println(strings.Compare("b", "a"))
	// Output:
	// -1
	// 0
	// 1
}

func ExampleContains() {
	fmt.Println(strings.Contains("seafood", "foo"))
	fmt.Println(strings.Contains("seafood", "bar"))
	fmt.Println(strings.Contains("seafood", ""))
	fmt.Println(strings.Contains("", ""))
	// Output:
	// true
	// false
	// true
	// true
}

func ExampleContainsAny() {
	fmt.Println(strings.ContainsAny("team", "i"))
	fmt.Println(strings.ContainsAny("fail", "ui"))
	fmt.Println(strings.ContainsAny("ure", "ui"))
	fmt.Println(strings.ContainsAny("failure", "ui"))
	fmt.Println(strings.ContainsAny("foo", ""))
	fmt.Println(strings.ContainsAny("", ""))
	// Output:
	// false
	// true
	// true
	// true
	// false
	// false
}

func ExampleContainsRune() {
	// Finds whether a string contains a particular Unicode code point.
	// The code point for the lowercase letter "a", for example, is 97.
	fmt.Println(strings.ContainsRune("aardvark", 97))
	fmt.Println(strings.ContainsRune("timeout", 97))
	// Output:
	// true
	// false
}

func ExampleContainsFunc() {
	f := func(r rune) bool {
		return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u'
	}
	fmt.Println(strings.ContainsFunc("hello", f))
	fmt.Println(strings.ContainsFunc("rhythms", f))
	// Output:
	// true
	// false
}

func ExampleCount() {
	fmt.Println(strings.Count("cheese", "e"))
	fmt.Println(strings.Count("five", "")) // before & after each rune
	// Output:
	// 3
	// 5
}

func ExampleCut() {
	show := func(s, sep string) {
		before, after := strings.Cut(s, sep)
		fmt.Printf("Cut(%q, %q) = %q, %q\n", s, sep, before, after)
	}
	show("Gopher", "Go")
	show("Gopher", "ph")
	show("Gopher", "er")
	show("Gopher", "Badger")
	// Output:
	// Cut("Gopher", "Go") = "", "pher"
	// Cut("Gopher", "ph") = "Go", "er"
	// Cut("Gopher", "er") = "Goph", ""
	// Cut("Gopher", "Badger") = "Gopher", ""
}

func ExampleCutPrefix() {
	show := func(s, prefix string) {
		after, found := strings.CutPrefix(s, prefix)
		fmt.Printf("CutPrefix(%q, %q) = %q, %v\n", s, prefix, after, found)
	}
	show("Gopher", "Go")
	show("Gopher", "ph")
	// Output:
	// CutPrefix("Gopher", "Go") = "pher", true
	// CutPrefix("Gopher", "ph") = "Gopher", false
}

func ExampleCutSuffix() {
	show := func(s, suffix string) {
		before, found := strings.CutSuffix(s, suffix)
		fmt.Printf("CutSuffix(%q, %q) = %q, %v\n", s, suffix, before, found)
	}
	show("Gopher", "Go")
	show("Gopher", "er")
	// Output:
	// CutSuffix("Gopher", "Go") = "Gopher", false
	// CutSuffix("Gopher", "er") = "Goph", true
}

func ExampleFields() {
	fmt.Printf("Fields are: %q", strings.Fields(nil, "  foo bar  baz   "))
	// Output: Fields are: ["foo" "bar" "baz"]
}

func ExampleFieldsFunc() {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsDigit(c)
	}
	fmt.Printf("Fields are: %q", strings.FieldsFunc(nil, "  foo1;bar2,baz3...", f))
	// Output: Fields are: ["foo1" "bar2" "baz3"]
}

func ExampleHasPrefix() {
	fmt.Println(strings.HasPrefix("Gopher", "Go"))
	fmt.Println(strings.HasPrefix("Gopher", "C"))
	fmt.Println(strings.HasPrefix("Gopher", ""))
	// Output:
	// true
	// false
	// true
}

func ExampleHasSuffix() {
	fmt.Println(strings.HasSuffix("Amigo", "go"))
	fmt.Println(strings.HasSuffix("Amigo", "O"))
	fmt.Println(strings.HasSuffix("Amigo", "Ami"))
	fmt.Println(strings.HasSuffix("Amigo", ""))
	// Output:
	// true
	// false
	// false
	// true
}

func ExampleIndex() {
	fmt.Println(strings.Index("chicken", "ken"))
	fmt.Println(strings.Index("chicken", "dmr"))
	// Output:
	// 4
	// -1
}

func ExampleIndexFunc() {
	f := func(c rune) bool {
		return unicode.Is(unicode.White_Space, c)
	}
	fmt.Println(strings.IndexFunc("Hello, 世界", f))
	fmt.Println(strings.IndexFunc("Helloworld", f))
	// Output:
	// 6
	// -1
}

func ExampleIndexAny() {
	fmt.Println(strings.IndexAny("chicken", "aeiouy"))
	fmt.Println(strings.IndexAny("crwth", "aeiouy"))
	// Output:
	// 2
	// -1
}

func ExampleIndexByte() {
	fmt.Println(strings.IndexByte("golang", 'g'))
	fmt.Println(strings.IndexByte("gophers", 'h'))
	fmt.Println(strings.IndexByte("golang", 'x'))
	// Output:
	// 0
	// 3
	// -1
}
func ExampleIndexRune() {
	fmt.Println(strings.IndexRune("chicken", 'k'))
	fmt.Println(strings.IndexRune("chicken", 'd'))
	// Output:
	// 4
	// -1
}

func ExampleLastIndex() {
	fmt.Println(strings.Index("go gopher", "go"))
	fmt.Println(strings.LastIndex("go gopher", "go"))
	fmt.Println(strings.LastIndex("go gopher", "rodent"))
	// Output:
	// 0
	// 3
	// -1
}

func ExampleLastIndexByte() {
	fmt.Println(strings.LastIndexByte("Hello, world", 'l'))
	fmt.Println(strings.LastIndexByte("Hello, world", 'o'))
	fmt.Println(strings.LastIndexByte("Hello, world", 'x'))
	// Output:
	// 10
	// 8
	// -1
}

func ExampleJoin() {
	s := []string{"foo", "bar", "baz"}
	fmt.Println(strings.Join(nil, s, ", "))
	// Output: foo, bar, baz
}

func ExampleRepeat() {
	fmt.Println("ba" + strings.Repeat(nil, "na", 2))
	// Output: banana
}

func ExampleReplace() {
	fmt.Println(strings.Replace(nil, "oink oink oink", "k", "ky", 2))
	fmt.Println(strings.Replace(nil, "oink oink oink", "oink", "moo", -1))
	// Output:
	// oinky oinky oink
	// moo moo moo
}

func ExampleReplaceAll() {
	fmt.Println(strings.ReplaceAll(nil, "oink oink oink", "oink", "moo"))
	// Output:
	// moo moo moo
}

func ExampleSplit() {
	fmt.Printf("%q\n", strings.Split(nil, "a,b,c", ","))
	fmt.Printf("%q\n", strings.Split(nil, "a man a plan a canal panama", "a "))
	fmt.Printf("%q\n", strings.Split(nil, " xyz ", ""))
	fmt.Printf("%q\n", strings.Split(nil, "", "Bernardo O'Higgins"))
	// Output:
	// ["a" "b" "c"]
	// ["" "man " "plan " "canal panama"]
	// [" " "x" "y" "z" " "]
	// [""]
}

func ExampleSplitN() {
	fmt.Printf("%q\n", strings.SplitN(nil, "a,b,c", ",", 2))
	z := strings.SplitN(nil, "a,b,c", ",", 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
	// Output:
	// ["a" "b,c"]
	// [] (nil = true)
}

func ExampleSplitAfter() {
	fmt.Printf("%q\n", strings.SplitAfter(nil, "a,b,c", ","))
	// Output: ["a," "b," "c"]
}

func ExampleMap() {
	rot13 := func(r rune) rune {
		if r >= 'A' && r <= 'Z' {
			return 'A' + (r-'A'+13)%26
		} else if r >= 'a' && r <= 'z' {
			return 'a' + (r-'a'+13)%26
		}
		return r
	}
	fmt.Println(strings.Map(nil, rot13, "'Twas brillig and the slithy gopher..."))
	// Output: 'Gjnf oevyyvt naq gur fyvgul tbcure...
}

func ExampleToUpper() {
	fmt.Println(strings.ToUpper(nil, "Gopher"))
	// Output: GOPHER
}

func ExampleToLower() {
	fmt.Println(strings.ToLower(nil, "Gopher"))
	// Output: gopher
}

func ExampleTrim() {
	fmt.Print(strings.Trim("¡¡¡Hello, Gophers!!!", "!¡"))
	// Output: Hello, Gophers
}

func ExampleTrimSpace() {
	fmt.Println(strings.TrimSpace(" \t\n Hello, Gophers \n\t\r\n"))
	// Output: Hello, Gophers
}

func ExampleTrimPrefix() {
	var s = "¡¡¡Hello, Gophers!!!"
	s = strings.TrimPrefix(s, "¡¡¡Hello, ")
	s = strings.TrimPrefix(s, "¡¡¡Howdy, ")
	fmt.Print(s)
	// Output: Gophers!!!
}

func ExampleTrimSuffix() {
	var s = "¡¡¡Hello, Gophers!!!"
	s = strings.TrimSuffix(s, ", Gophers!!!")
	s = strings.TrimSuffix(s, ", Marmots!!!")
	fmt.Print(s)
	// Output: ¡¡¡Hello
}

func ExampleTrimFunc() {
	fmt.Print(strings.TrimFunc("¡¡¡Hello, Gophers!!!", func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}))
	// Output: Hello, Gophers
}

func ExampleTrimLeft() {
	fmt.Print(strings.TrimLeft("¡¡¡Hello, Gophers!!!", "!¡"))
	// Output: Hello, Gophers!!!
}

func ExampleTrimRight() {
	fmt.Print(strings.TrimRight("¡¡¡Hello, Gophers!!!", "!¡"))
	// Output: ¡¡¡Hello, Gophers
}
