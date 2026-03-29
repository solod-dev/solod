package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/os"
)

func fileTest() {
	{
		// OpenFile with O_CREATE | O_WRONLY | O_TRUNC.
		name := "test_openfile.txt"
		f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			panic("OpenFile create failed")
		}
		defer os.Remove(name)
		f.Write([]byte("openfile"))
		f.Close()

		b, err := os.ReadFile(nil, name)
		if err != nil {
			panic("ReadFile after OpenFile failed")
		}
		defer mem.FreeSlice(nil, b)
		if string(b) != "openfile" {
			panic("OpenFile: wrong data")
		}
	}
	{
		// OpenFile with O_RDONLY.
		name := "test_openfile_rd.txt"
		os.WriteFile(name, []byte("readonly"), 0o666)
		defer os.Remove(name)

		f, err := os.OpenFile(name, os.O_RDONLY, 0)
		if err != nil {
			panic("OpenFile rdonly failed")
		}
		buf := make([]byte, 16)
		n, err := f.Read(buf)
		if err != nil {
			panic("Read from rdonly failed")
		}
		if string(buf[:n]) != "readonly" {
			panic("OpenFile rdonly: wrong data")
		}
		f.Close()
	}
	{
		// File.Name.
		name := "test_filename.txt"
		f, err := os.Create(name)
		if err != nil {
			panic("Create failed")
		}
		defer os.Remove(name)
		if f.Name() != name {
			panic("Name: wrong")
		}
		f.Close()
	}
	{
		// Link and Readlink (via symlink).
		target := "test_link_target.txt"
		os.WriteFile(target, []byte("linked"), 0o666)
		defer os.Remove(target)

		// Hard link.
		hard := "test_hard_link.txt"
		err := os.Link(target, hard)
		if err != nil {
			panic("Link failed")
		}
		defer os.Remove(hard)

		b, err := os.ReadFile(nil, hard)
		if err != nil {
			panic("ReadFile hard link failed")
		}
		defer mem.FreeSlice(nil, b)
		if string(b) != "linked" {
			panic("Hard link: wrong data")
		}
	}
	{
		// Symlink and Readlink.
		target := "test_sym_target.txt"
		os.WriteFile(target, []byte("sym"), 0o666)
		defer os.Remove(target)

		link := "test_sym_link"
		err := os.Symlink(target, link)
		if err != nil {
			panic("Symlink failed")
		}
		defer os.Remove(link)

		var rlBuf [os.MaxPathLen]byte
		dest, err := os.Readlink(rlBuf[:], link)
		if err != nil {
			panic("Readlink failed")
		}
		if dest != target {
			panic("Readlink: wrong target")
		}
	}
	{
		// Mkdir and Chdir.
		dir := "test_mkdir_dir"
		err := os.Mkdir(dir, 0o755)
		if err != nil {
			panic("Mkdir failed")
		}
		defer os.Remove(dir)

		// Get current dir.
		var wdBuf [os.MaxPathLen]byte
		origWd, err := os.Getwd(wdBuf[:])
		if err != nil {
			panic("Getwd failed")
		}

		// Change to new dir.
		err = os.Chdir(dir)
		if err != nil {
			panic("Chdir failed")
		}

		// Verify we moved.
		var wdBuf2 [os.MaxPathLen]byte
		newWd, err := os.Getwd(wdBuf2[:])
		if err != nil {
			panic("Getwd after Chdir failed")
		}
		if newWd == origWd {
			panic("Chdir: dir did not change")
		}

		// Change back.
		os.Chdir(origWd)
	}
	{
		// Truncate.
		name := "test_truncate.txt"
		os.WriteFile(name, []byte("abcdef"), 0o666)
		defer os.Remove(name)

		err := os.Truncate(name, 3)
		if err != nil {
			panic("Truncate failed")
		}
		b, err := os.ReadFile(nil, name)
		if err != nil {
			panic("ReadFile after Truncate failed")
		}
		defer mem.FreeSlice(nil, b)
		if string(b) != "abc" {
			panic("Truncate: wrong data")
		}
	}
	{
		// OpenFile with O_APPEND.
		name := "test_append.txt"
		os.WriteFile(name, []byte("hello"), 0o666)
		defer os.Remove(name)

		f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND, 0)
		if err != nil {
			panic("OpenFile append failed")
		}
		f.Write([]byte(" world"))
		f.Close()

		b, err := os.ReadFile(nil, name)
		if err != nil {
			panic("ReadFile after append failed")
		}
		defer mem.FreeSlice(nil, b)
		if string(b) != "hello world" {
			panic("Append: wrong data")
		}
	}
	{
		// Chtimes - just verify it doesn't error.
		name := "test_chtimes.txt"
		os.WriteFile(name, []byte("times"), 0o666)
		defer os.Remove(name)

		fi, err := os.Stat(name)
		if err != nil {
			panic("Stat for Chtimes failed")
		}
		mt := fi.ModTime()
		err = os.Chtimes(name, mt, mt)
		if err != nil {
			panic("Chtimes failed")
		}
	}
	{
		// Chown with -1, -1 (no change) - should succeed.
		name := "test_chown.txt"
		os.WriteFile(name, []byte("chown"), 0o666)
		defer os.Remove(name)

		err := os.Chown(name, -1, -1)
		if err != nil {
			panic("Chown failed")
		}
	}
	{
		// Lchown with -1, -1 (no change) - should succeed.
		name := "test_lchown.txt"
		os.WriteFile(name, []byte("lchown"), 0o666)
		defer os.Remove(name)

		err := os.Lchown(name, -1, -1)
		if err != nil {
			panic("Lchown failed")
		}
	}
	{
		// Remove.
		name := "test_remove.txt"
		err := os.WriteFile(name, []byte("tmp"), 0o666)
		if err != nil {
			panic("WriteFile failed")
		}
		err = os.Remove(name)
		if err != nil {
			panic("Remove failed")
		}
		_, err = os.Open(name)
		if err == nil {
			panic("Open after Remove should fail")
		}
	}
	{
		// Rename.
		oldName := "test_old.txt"
		newName := "test_new.txt"
		os.WriteFile(oldName, []byte("renamed"), 0o666)
		err := os.Rename(oldName, newName)
		if err != nil {
			panic("Rename failed")
		}
		defer os.Remove(newName)
		b, err := os.ReadFile(nil, newName)
		if err != nil {
			panic("ReadFile after Rename failed")
		}
		defer mem.FreeSlice(nil, b)
		if string(b) != "renamed" {
			panic("Rename: wrong data")
		}
	}
	{
		// ErrExist - try to create dir that already exists.
		name := "test_exist_dir"
		os.Mkdir(name, 0o755)
		defer os.Remove(name)

		err := os.Mkdir(name, 0o755)
		if err != os.ErrExist {
			panic("Mkdir existing: wrong error")
		}
	}
	{
		// ErrNotExist.
		_, err := os.Open("nonexistent_file.txt")
		if err != os.ErrNotExist {
			panic("Open nonexistent: wrong error")
		}
	}
	{
		// Verify OpenFile nonexistent returns ErrNotExist.
		_, err := os.OpenFile("nonexistent_open.txt", os.O_RDONLY, 0)
		if err != os.ErrNotExist {
			panic("OpenFile nonexistent: wrong error")
		}
	}
}
