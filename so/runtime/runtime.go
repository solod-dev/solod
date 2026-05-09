// Package runtime provides information about
// the environment where the program was compiled.
package runtime

//so:embed runtime.h
var runtime_h string

// GOOS is the running program's operating system target:
// one of darwin, linux, windows, and so on.
//
//so:extern
const GOOS string = "unknown"

// GOARCH is the running program's architecture target:
// one of amd64, arm64, and so on.
//
//so:extern
const GOARCH string = "unknown"

// Recognized GOOS/GOARCH pairs are:
// GOOS		GOARCH
// bare		amd64
// bare		arm64
// bare		386
// bare		riscv64
// bare		wasm
// darwin	amd64
// darwin	arm64
// linux	amd64
// linux	arm64
// linux	386
// linux	riscv64
// wasip1	wasm
// windows	amd64
// windows	arm64

//so:extern runtime_buildVersion
var buildVersion string

// Version returns the So tree's version string.
// It is either the commit hash and date at the time of the build or,
// when possible, a release tag like "v0.1.0".
func Version() string {
	return buildVersion
}

// Seed returns a random 64-bit seed.
// It's cryptographically secure on macOS/Linux (arc4random_buf/getrandom)
// and falls back to a time-based seed on other platforms.
//
//so:extern
func Seed() uint64 {
	return 0
}
