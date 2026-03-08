// So provides file I/O through the `so/c/stdio` package,
// which wraps the C standard I/O functions. This example
// shows how to write and read files.
package main

import (
	"github.com/nalgeon/solod/so/c"
	"github.com/nalgeon/solod/so/c/stdio"
)

func main() {
	var f *stdio.File

	{
		// Write a file using Fputs and Fprintf.
		f = stdio.Fopen("/tmp/so-example.txt", "w")
		if f == nil {
			panic("failed to open file for writing")
		}
		stdio.Fputs("hello, file!\n", f)
		stdio.Fprintf(f, "the answer is %d\n", 42)
		stdio.Fclose(f)
	}

	{
		// Read the file back line by line using Fgets.
		f = stdio.Fopen("/tmp/so-example.txt", "r")
		if f == nil {
			panic("failed to open file for reading")
		}

		// Use Ftell to check the starting position.
		pos := stdio.Ftell(f)
		println("starting position:", pos)

		// Read and print each line until EOF.
		// Use `c.String` to convert a *byte buffer to a So string.
		var buf [256]byte
		println("file contents:")
		for stdio.Fgets(&buf[0], 256, f) != nil {
			line := c.String(&buf[0])
			stdio.Fputs(line, stdio.Stdout)
		}

		// Check end-of-file and error status.
		if stdio.Feof(f) {
			println("reached end of file")
		}
		if stdio.Ferror(f) {
			println("error reading file")
		}

		// Not closing the file here because
		// we want to seek and read again below.
	}

	{
		// Seek back to the beginning and read the first line again.
		var buf [256]byte
		stdio.Fseek(f, 0, stdio.SeekSet)
		stdio.Fgets(&buf[0], 256, f)
		line := c.String(&buf[0])
		stdio.Fputs("first line after seek: ", stdio.Stdout)
		stdio.Fputs(line, stdio.Stdout)
		stdio.Fclose(f)
	}
}
