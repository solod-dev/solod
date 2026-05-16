package os

import "solod.dev/so/time"

// A FileInfo describes a file and is returned by [Stat] and [Lstat].
type FileInfo struct {
	name    string
	size    int64
	mode    FileMode
	modTime time.Time
	dev     uint64
	ino     uint64
}

// Name returns the base name of the file.
func (fi *FileInfo) Name() string { return fi.name }

// Size returns the length in bytes for regular files; system-dependent for others.
func (fi *FileInfo) Size() int64 { return fi.size }

// Mode returns the file mode bits.
func (fi *FileInfo) Mode() FileMode { return fi.mode }

// ModTime returns the modification time.
func (fi *FileInfo) ModTime() time.Time { return fi.modTime }

// IsDir reports whether the file is a directory.
func (fi *FileInfo) IsDir() bool { return fi.mode.IsDir() }

// baseName returns the last element of the path.
func baseName(path string) string {
	i := len(path) - 1
	// Strip trailing slashes.
	for i > 0 && path[i] == '/' {
		i--
	}
	end := i + 1
	// Find the last slash.
	for i >= 0 && path[i] != '/' {
		i--
	}
	if end == 0 {
		return "."
	}
	return path[i+1 : end]
}

// Stat returns a [FileInfo] describing the named file.
func Stat(name string) (FileInfo, error) {
	r := os_stat(name)
	if !r.ok {
		return FileInfo{}, mapError()
	}
	fmode := mode_t(r.mode).toFileMode()
	return FileInfo{
		name:    baseName(name),
		size:    r.size,
		mode:    fmode,
		modTime: time.Unix(r.modSec, r.modNsec),
		dev:     r.dev,
		ino:     r.ino,
	}, nil
}

// Lstat returns a [FileInfo] describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link. Lstat makes no attempt to follow the link.
func Lstat(name string) (FileInfo, error) {
	r := os_lstat(name)
	if !r.ok {
		return FileInfo{}, mapError()
	}
	fmode := mode_t(r.mode).toFileMode()
	return FileInfo{
		name:    baseName(name),
		size:    r.size,
		mode:    fmode,
		modTime: time.Unix(r.modSec, r.modNsec),
		dev:     r.dev,
		ino:     r.ino,
	}, nil
}
