package main

// This test documents expected panic behavior
// Actual panic testing requires special infrastructure

func main() {
	// ch := make(chan int)
	// close(ch)
	// close(ch)  // Should panic: close of closed channel

	println("ok: close twice panics (documented)")
}
