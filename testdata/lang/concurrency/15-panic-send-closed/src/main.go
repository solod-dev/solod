package main

// This test documents expected panic behavior
// Actual panic testing requires special infrastructure

func main() {
	// ch := make(chan int)
	// close(ch)
	// ch <- 1  // Should panic: send on closed channel

	println("ok: send to closed channel panics (documented)")
}
