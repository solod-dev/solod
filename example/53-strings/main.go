// The `so/strings` package provides many useful
// string-related functions. Here are some examples
// to give you a sense of the package.
package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
	"solod.dev/so/strings"
)

func main() {
	contains := strings.Contains("test", "es")
	fmt.Printf("Contains:  %d\n", contains)

	count := strings.Count("test", "t")
	fmt.Printf("Count:     %d\n", int32(count))

	hasPrefix := strings.HasPrefix("test", "te")
	fmt.Printf("HasPrefix: %d\n", hasPrefix)

	hasSuffix := strings.HasSuffix("test", "st")
	fmt.Printf("HasSuffix: %d\n", hasSuffix)

	index := strings.Index("test", "e")
	fmt.Printf("Index:     %d\n", int32(index))

	// Some functions allocate new strings to return the result.
	// We can use `mem.FreeString` to free them and avoid memory leaks.
	joined := strings.Join(nil, []string{"a", "b"}, "-")
	defer mem.FreeString(nil, joined)
	fmt.Printf("Join:      %s\n", joined)

	repeated := strings.Repeat(nil, "a", 5)
	defer mem.FreeString(nil, repeated)
	fmt.Printf("Repeat:    %s\n", repeated)

	replacedAll := strings.Replace(nil, "foo", "o", "0", -1)
	defer mem.FreeString(nil, replacedAll)
	fmt.Printf("Replace:   %s\n", replacedAll)

	replacedOnce := strings.Replace(nil, "foo", "o", "0", 1)
	defer mem.FreeString(nil, replacedOnce)
	fmt.Printf("Replace:   %s\n", replacedOnce)

	splitted := strings.Split(nil, "a-b-c-d-e", "-")
	defer mem.FreeSlice(nil, splitted)
	for i, s := range splitted {
		fmt.Printf("Split %d    %s\n", i, s)
	}

	lowered := strings.ToLower(nil, "TEST")
	defer mem.FreeString(nil, lowered)
	fmt.Printf("ToLower:   %s\n", lowered)

	uppered := strings.ToUpper(nil, "test")
	defer mem.FreeString(nil, uppered)
	fmt.Printf("ToUpper:   %s\n", uppered)
}
