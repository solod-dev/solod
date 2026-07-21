package mem

// Arena is a memory allocator that bump-allocates linearly within a fixed
// buffer. [Arena.Free] reclaims the last allocation if the pointer matches;
// otherwise it is a no-op. Use [Arena.Reset] to reclaim all memory at once.
//
// Arena is not thread-safe.
type Arena struct {
	buf       []byte
	offset    int
	lastStart int
}

// NewArena creates an arena allocator backed by the given buffer.
func NewArena(buf []byte) Arena {
	return Arena{buf: buf}
}

func (a *Arena) Alloc(size int, align int) (any, error) {
	if size <= 0 {
		panic("mem: invalid allocation size")
	}
	if align <= 0 || (align&(align-1)) != 0 {
		panic("mem: invalid alignment")
	}
	aligned := (a.offset + align - 1) & ^(align - 1)
	if aligned+size > len(a.buf) {
		return nil, ErrOutOfMemory
	}
	clear(a.buf[aligned : aligned+size])
	a.lastStart = aligned
	a.offset = aligned + size
	return &a.buf[aligned], nil
}

func (a *Arena) Realloc(ptr any, oldSize int, newSize int, align int) (any, error) {
	if oldSize <= 0 || newSize <= 0 {
		panic("mem: invalid allocation size")
	}
	if align <= 0 || (align&(align-1)) != 0 {
		panic("mem: invalid alignment")
	}
	// Last allocation: resize in place.
	if ptr == &a.buf[a.lastStart] && a.lastStart+oldSize == a.offset {
		newEnd := a.lastStart + newSize
		if newEnd > len(a.buf) {
			return nil, ErrOutOfMemory
		}
		if newSize > oldSize {
			clear(a.buf[a.offset:newEnd])
		}
		a.offset = newEnd
		return ptr, nil
	}
	// Not the last allocation, shrinking: return same pointer.
	if newSize <= oldSize {
		return ptr, nil
	}
	// Not the last allocation, growing: allocate new and copy.
	newPtr, err := a.Alloc(newSize, align)
	if err != nil {
		return nil, err
	}
	memmove(newPtr, ptr, uintptr(oldSize))
	return newPtr, nil
}

func (a *Arena) Free(ptr any, size int, align int) {
	_ = align
	// Last allocation: reclaim the space.
	if ptr == &a.buf[a.lastStart] && a.lastStart+size == a.offset {
		a.offset = a.lastStart
	}
}

// Reset reclaims all allocated memory, allowing the arena to be reused.
func (a *Arena) Reset() {
	a.offset = 0
	a.lastStart = 0
}
