package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/os"
)

func dirTest() {
	{
		// ReadDir on a directory with known contents.
		dirName := "test_readdir"
		os.Mkdir(dirName, 0o755)
		defer os.Remove(dirName)

		os.WriteFile(dirName+"/aaa.txt", []byte("hello"), 0o666)
		defer os.Remove(dirName + "/aaa.txt")

		os.WriteFile(dirName+"/bbb.txt", []byte("world"), 0o666)
		defer os.Remove(dirName + "/bbb.txt")

		os.Mkdir(dirName+"/subdir", 0o755)
		defer os.Remove(dirName + "/subdir")

		entries, err := os.ReadDir(nil, dirName)
		if err != nil {
			panic("ReadDir failed")
		}
		defer os.FreeDirEntry(nil, entries)

		if len(entries) != 3 {
			fmt.Printf("ReadDir: expected 3 entries, got %d\n", len(entries))
			panic("ReadDir: wrong count")
		}

		// Check that we find each expected entry.
		foundFile := false
		foundDir := false
		for _, entry := range entries {
			if entry.Name == "aaa.txt" || entry.Name == "bbb.txt" {
				foundFile = true
				if entry.IsDir {
					panic("ReadDir: file should not be dir")
				}
			}
			if entry.Name == "subdir" {
				foundDir = true
				if !entry.IsDir {
					panic("ReadDir: subdir should be dir")
				}
				if entry.Type&os.ModeDir == 0 {
					panic("ReadDir: subdir Type should have ModeDir")
				}
			}
		}
		if !foundFile {
			panic("ReadDir: did not find file entries")
		}
		if !foundDir {
			panic("ReadDir: did not find subdir")
		}
	}
	{
		// ReadDir on nonexistent directory.
		_, err := os.ReadDir(nil, "nonexistent_dir_xyz")
		if err != os.ErrNotExist {
			panic("ReadDir nonexistent: wrong error")
		}
	}
}
