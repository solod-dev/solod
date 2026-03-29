package main

import "solod.dev/so/os"

func statTest() {
	{
		// Stat on a regular file.
		name := "test_stat.txt"
		os.WriteFile(name, []byte("hello"), 0o666)
		defer os.Remove(name)

		fi, err := os.Stat(name)
		if err != nil {
			panic("Stat failed")
		}
		if fi.Name() != "test_stat.txt" {
			panic("Stat: wrong name")
		}
		if fi.Size() != 5 {
			panic("Stat: wrong size")
		}
		if !fi.Mode().IsRegular() {
			panic("Stat: not regular")
		}
		if fi.IsDir() {
			panic("Stat: should not be dir")
		}
	}
	{
		// Stat on a directory.
		name := "test_stat_dir"
		os.Mkdir(name, 0o755)
		defer os.Remove(name)

		fi, err := os.Stat(name)
		if err != nil {
			panic("Stat dir failed")
		}
		if fi.Name() != "test_stat_dir" {
			panic("Stat dir: wrong name")
		}
		if !fi.IsDir() {
			panic("Stat dir: should be dir")
		}
		if fi.Mode().IsRegular() {
			panic("Stat dir: should not be regular")
		}
	}
	{
		// Lstat on a symlink.
		target := "test_lstat_target.txt"
		link := "test_lstat_link"
		os.WriteFile(target, []byte("target"), 0o666)
		defer os.Remove(target)
		os.Symlink(target, link)
		defer os.Remove(link)

		// Lstat returns info about the link itself.
		fi, err := os.Lstat(link)
		if err != nil {
			panic("Lstat failed")
		}
		if fi.Name() != "test_lstat_link" {
			panic("Lstat: wrong name")
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			panic("Lstat: should be symlink")
		}

		// Stat follows the link.
		fi2, err := os.Stat(link)
		if err != nil {
			panic("Stat through link failed")
		}
		if fi2.Size() != 6 {
			panic("Stat through link: wrong size")
		}
		if fi2.Mode()&os.ModeSymlink != 0 {
			panic("Stat through link: should not be symlink")
		}
	}
	{
		// SameFile.
		name := "test_samefile.txt"
		os.WriteFile(name, []byte("same"), 0o666)
		defer os.Remove(name)

		fi1, err := os.Stat(name)
		if err != nil {
			panic("Stat 1 failed")
		}
		fi2, err := os.Stat(name)
		if err != nil {
			panic("Stat 2 failed")
		}
		if !os.SameFile(fi1, fi2) {
			panic("SameFile: should be same")
		}

		name2 := "test_samefile2.txt"
		os.WriteFile(name2, []byte("other"), 0o666)
		defer os.Remove(name2)

		fi3, err := os.Stat(name2)
		if err != nil {
			panic("Stat 3 failed")
		}
		if os.SameFile(fi1, fi3) {
			panic("SameFile: should be different")
		}
	}
	{
		// Stat on nonexistent file.
		_, err := os.Stat("nonexistent_stat.txt")
		if err != os.ErrNotExist {
			panic("Stat nonexistent: wrong error")
		}
	}
	{
		// Chmod and permission check.
		name := "test_chmod.txt"
		os.WriteFile(name, []byte("chmod"), 0o666)
		defer os.Remove(name)

		err := os.Chmod(name, 0o644)
		if err != nil {
			panic("Chmod failed")
		}
		fi, err := os.Stat(name)
		if err != nil {
			panic("Stat after Chmod failed")
		}
		if fi.Mode().Perm() != 0o644 {
			panic("Chmod: wrong perm")
		}
	}
}
