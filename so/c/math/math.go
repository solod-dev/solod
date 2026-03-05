// Package math wraps the C <math.h> header.
// It offers mathematical functions and constants.
package math

import _ "embed"

//so:embed math.h
var math_h string

// Pi is the ratio of a circle's circumference to its diameter.
//
//so:extern
var Pi float64

// E is the base of natural logarithms.
//
//so:extern
var E float64

// Inf is positive infinity.
//
//so:extern
var Inf float64

// Abs returns the absolute value of x.
//
//so:extern
func Abs(x float64) float64 { return 0 }

// Sqrt returns the square root of x.
//
//so:extern
func Sqrt(x float64) float64 { return 0 }

// Pow returns x raised to the power y.
//
//so:extern
func Pow(x float64, y float64) float64 { return 0 }

// Floor returns the largest integer value less than or equal to x.
//
//so:extern
func Floor(x float64) float64 { return 0 }

// Ceil returns the smallest integer value greater than or equal to x.
//
//so:extern
func Ceil(x float64) float64 { return 0 }

// Round returns the nearest integer value, rounding halfway cases away from zero.
//
//so:extern
func Round(x float64) float64 { return 0 }

// Log returns the natural logarithm of x.
//
//so:extern
func Log(x float64) float64 { return 0 }

// Log2 returns the base-2 logarithm of x.
//
//so:extern
func Log2(x float64) float64 { return 0 }

// Log10 returns the base-10 logarithm of x.
//
//so:extern
func Log10(x float64) float64 { return 0 }

// Exp returns e raised to the power x.
//
//so:extern
func Exp(x float64) float64 { return 0 }

// Sin returns the sine of x (in radians).
//
//so:extern
func Sin(x float64) float64 { return 0 }

// Cos returns the cosine of x (in radians).
//
//so:extern
func Cos(x float64) float64 { return 0 }

// Atan2 returns the arc tangent of y/x, using the signs of the two
// to determine the quadrant of the return value.
//
//so:extern
func Atan2(y float64, x float64) float64 { return 0 }

// Fmin returns the smaller of x and y.
//
//so:extern
func Fmin(x float64, y float64) float64 { return 0 }

// Fmax returns the larger of x and y.
//
//so:extern
func Fmax(x float64, y float64) float64 { return 0 }

// Fmod returns the floating-point remainder of x/y.
//
//so:extern
func Fmod(x float64, y float64) float64 { return 0 }
