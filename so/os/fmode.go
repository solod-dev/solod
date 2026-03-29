package os

// FileMode represents a file's mode and permission bits.
// The bits have the same definition on all systems, so that
// information about files can be moved from one system
// to another portably.
type FileMode uint32

// The defined file mode bits are the most significant bits of the FileMode.
const (
	ModeDir        FileMode = 1 << (32 - 1 - 0)
	ModeSymlink    FileMode = 1 << (32 - 1 - 1)
	ModeNamedPipe  FileMode = 1 << (32 - 1 - 2)
	ModeSocket     FileMode = 1 << (32 - 1 - 3)
	ModeDevice     FileMode = 1 << (32 - 1 - 4)
	ModeCharDevice FileMode = 1 << (32 - 1 - 5)
	ModeSetuid     FileMode = 1 << (32 - 1 - 6)
	ModeSetgid     FileMode = 1 << (32 - 1 - 7)
	ModeSticky     FileMode = 1 << (32 - 1 - 8)
	ModeIrregular  FileMode = 1 << (32 - 1 - 9)
)

// ModePerm is the Unix permission bits.
const ModePerm FileMode = 0o777

// IsDir reports whether m describes a directory.
func (m FileMode) IsDir() bool {
	return m&ModeDir != 0
}

// IsRegular reports whether m describes a regular file.
func (m FileMode) IsRegular() bool {
	return m&(ModeDir|ModeSymlink|ModeNamedPipe|ModeSocket|ModeDevice|ModeCharDevice|ModeIrregular) == 0
}

// Perm returns the Unix permission bits in m.
func (m FileMode) Perm() FileMode {
	return m & ModePerm
}

// makePosixMode converts Go FileMode bits to POSIX mode_t bits.
func makePosixMode(fmode FileMode) mode_t {
	pmode := mode_t(fmode & 0777)
	if fmode&ModeSetuid != 0 {
		pmode |= sISUID
	}
	if fmode&ModeSetgid != 0 {
		pmode |= sISGID
	}
	if fmode&ModeSticky != 0 {
		pmode |= sISVTX
	}
	return pmode
}

// toFileMode converts POSIX mode_t bits to Go FileMode bits.
// Go FileMode layout: high bits are type/special, low 9 bits are permissions.
func (m mode_t) toFileMode() FileMode {
	fmode := FileMode(m & 0777) // permission bits pass through
	switch m & sIFMT {
	case sIFDIR:
		fmode |= ModeDir
	case sIFLNK:
		fmode |= ModeSymlink
	case sIFIFO:
		fmode |= ModeNamedPipe
	case sIFSOCK:
		fmode |= ModeSocket
	case sIFBLK:
		fmode |= ModeDevice
	case sIFCHR:
		fmode |= ModeCharDevice
	case sIFREG:
		// no special bit for regular files
	}
	if m&sISUID != 0 {
		fmode |= ModeSetuid
	}
	if m&sISGID != 0 {
		fmode |= ModeSetgid
	}
	if m&sISVTX != 0 {
		fmode |= ModeSticky
	}
	return fmode
}
