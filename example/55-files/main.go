// So provides file I/O through the `so/os` package
// and formatted output through the `so/fmt` package.
// This example shows how to write and read files.
package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/io"
	"solod.dev/so/os"
)

func main() {
	{
		// Write a file using WriteString and Fprintf.
		f, err := os.Create("/tmp/so-example.txt")
		if err != nil {
			panic("failed to create file")
		}
		defer f.Close()
		f.WriteString("hello, file!\n")
		fmt.Fprintf(&f, "the answer is %d\n", 42)
	}

	{
		// Read the file back.
		f, err := os.Open("/tmp/so-example.txt")
		if err != nil {
			panic("failed to open file for reading")
		}
		defer f.Close()

		// Use Seek with SeekCurrent to check the starting position.
		pos, _ := f.Seek(0, io.SeekCurrent)
		fmt.Printf("starting position: %d\n", pos)

		// Read and print the file contents.
		buf := make([]byte, 256)
		fmt.Println("file contents:")
		for {
			n, err := f.Read(buf)
			if n > 0 {
				fmt.Print(string(buf[:n]))
			}
			if err == io.EOF {
				fmt.Println("reached end of file")
				break
			}
			if err != nil {
				fmt.Println("error reading file")
				break
			}
		}

		// Seek back to the beginning and read the file again.
		buf = make([]byte, 256)
		f.Seek(0, io.SeekStart)
		n, _ := f.Read(buf)
		fmt.Println("content after seek:")
		fmt.Print(string(buf[:n]))
		f.Close()
	}

	{
		// Delete the file.
		err := os.Remove("/tmp/so-example.txt")
		if err != nil {
			panic("failed to delete file")
		}
		fmt.Println("file deleted")
	}
}
