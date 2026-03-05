package ctype

import _ "embed"

//so:embed ctype.h
var ctype_h string

// IsAlnum reports whether ch is an alphanumeric character
// as classified by the current C locale.
//
// In the default locale, the following characters are alphanumeric:
//   - Digits (0123456789),
//   - Uppercase letters (ABCDEFGHIJKLMNOPQRSTUVWXYZ),
//   - Lowercase letters (abcdefghijklmnopqrstuvwxyz).
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsAlnum(ch int) bool { return false }

// IsAlpha reports whether ch is an alphabetic character, meaning either an
// uppercase letter (ABCDEFGHIJKLMNOPQRSTUVWXYZ), or a lowercase letter
// (abcdefghijklmnopqrstuvwxyz).
//
// In locales other than "C", an alphabetic character is a character for which
// [IsUpper] or [IsLower] returns true or any other character considered alphabetic
// by the locale. In any case, [IsCntrl], [IsDigit], [IsPunct] and [IsSpace] will
// return false for this character.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsAlpha(ch int) bool { return false }

// IsBlank reports whether the ch is a blank character in the current C locale.
// In the default C locale, only space (0x20) and horizontal tab (0x09) are
// classified as blank.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsBlank(ch int) bool { return false }

// IsCntrl reports whether ch is a control character (codes 0x00-0x1F and 0x7F).
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsCntrl(ch int) bool { return false }

// IsDigit reports whether ch is a decimal digit character (0-9).
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsDigit(ch int) bool { return false }

// IsGraph reports whether the ch has a graphical representation, meaning it is
// either a number (0123456789), an uppercase letter (ABCDEFGHIJKLMNOPQRSTUVWXYZ),
// a lowercase letter (abcdefghijklmnopqrstuvwxyz), or a punctuation character
// (!"#$%&'()*+,-./:;<=>?@[\]^_{|}~), or any graphical character specific to the
// current C locale.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsGraph(ch int) bool { return false }

// IsLower reports whether ch is classified as a lowercase character according to
// the current C locale. In the default "C" locale, islower returns true only for
// the lowercase letters (abcdefghijklmnopqrstuvwxyz).
//
// If IsLower returns true, it is guaranteed that [IsCntrl], [IsDigit], [IsPunct]
// and [IsSpace] return false for the same character in the same C locale.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsLower(ch int) bool { return false }

// IsPrint reports whether ch can be printed, meaning it is either a number
// (0123456789), an uppercase letter (ABCDEFGHIJKLMNOPQRSTUVWXYZ), a lowercase
// letter (abcdefghijklmnopqrstuvwxyz), a punctuation character
// (!"#$%&'()*+,-./:;<=>?@[\]^_{|}~), or space, or any character classified as
// printable by the current C locale.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsPrint(ch int) bool { return false }

// IsPunct reports whether the given character is a punctuation character in
// the current C locale. The default C locale classifies the characters
// !"#$%&'()*+,-./:;<=>?@[\]^_{|}~ as punctuation.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsPunct(ch int) bool { return false }

// IsSpace reports whether ch is either a standard white-space character:
//   - Space (0x20 = ' '),
//   - Form feed (0x0c = \f),
//   - Line feed (0x0a = \n),
//   - Carriage return (0x0d = \r),
//   - Horizontal tab (0x09 = \t),
//   - Vertical tab (0x0b = \v),
//   - Or a locale-specific white-space character.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsSpace(ch int) bool { return false }

// IsUpper reports whether the given character is an uppercase character according
// to the current C locale. In the default "C" locale, isupper returns true only
// for the uppercase letters (ABCDEFGHIJKLMNOPQRSTUVWXYZ).
//
// If IsUpper returns true, it is guaranteed that [IsCntrl], [IsDigit], [IsPunct]
// and [IsSpace] return false for the same character in the same C locale.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsUpper(ch int) bool { return false }

// IsXDigit reports whether the given character is a hexadecimal numeric character
// (0123456789abcdefABCDEF) or is classified as a hexadecimal character.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func IsXDigit(ch int) bool { return false }

// ToLower converts ch to lowercase according to the character conversion rules
// defined by the currently installed C locale.
//
// In the default "C" locale, the following uppercase letters ABCDEFGHIJKLMNOPQRSTUVWXYZ
// are replaced with respective lowercase letters abcdefghijklmnopqrstuvwxyz.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func ToLower(ch int) int { return 0 }

// ToUpper converts ch to uppercase according to the character conversion rules
// defined by the currently installed C locale.
//
// In the default "C" locale, the following lowercase letters abcdefghijklmnopqrstuvwxyz
// are replaced with respective uppercase letters ABCDEFGHIJKLMNOPQRSTUVWXYZ.
//
// The behavior is undefined if the value of ch is not representable
// as unsigned char and is not equal to [stdio.EOF].
//
//so:extern
func ToUpper(ch int) int { return 0 }
