// Package assert wraps the C <assert.h> header.
// It offers functions for checking conditions that should always be true.
package assert

import _ "embed"

// Enabled reports whether assertions are enabled.
// Cannot be changed at runtime.
//
//so:extern
var Enabled bool

//so:embed assert.h
var assert_h string

// Assert aborts the program if the condition is not true.
// If [Enabled] is false, does nothing.
//
//so:extern
func Assert(cond bool) {}

// Assertf aborts the program with the given message
// if the condition is not true.
// If [Enabled] is false, does nothing.
//
//so:extern
func Assertf(cond bool, msg string) {}
