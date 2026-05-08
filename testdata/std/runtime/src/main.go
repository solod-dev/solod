package main

import "solod.dev/so/runtime"

func main() {
	println(runtime.Version(), runtime.GOOS, runtime.GOARCH)
	{
		// Version.
		v := runtime.Version()
		if len(v) == 0 {
			panic("Empty version")
		}
	}
	{
		// GOOS.
		os := runtime.GOOS
		if os != "darwin" && os != "linux" && os != "windows" {
			panic("Unexpected GOOS")
		}
	}
	{
		// GOARCH.
		arch := runtime.GOARCH
		if arch != "amd64" && arch != "arm64" && arch != "386" && arch != "riscv64" {
			panic("Unexpected GOARCH")
		}
	}
}
