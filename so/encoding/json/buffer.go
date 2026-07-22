package json

import (
	"solod.dev/so/io"
	"solod.dev/so/mem"
)

const (
	defaultBufSize = 4096 // size a streaming buffer starts at by default
	minBufSize     = 16   // smallest possible buffer to hand out
	maxEmptyRead   = 100  // give up on a reader that never makes progress
)

// buffer is the byte source a [Decoder] scans. It gives the decoder a cursor
// (peek/advance) and a way to bracket a token (hold/take), hiding where the
// bytes come from.
//
// There are two modes:
//
//   - Fixed: the whole document is already in memory. The buffer borrows
//     the caller's slice, never allocates, never refills, and never writes
//     to it. EOF is the end of the slice.
//   - Streaming: bytes are pulled from an [io.Reader] into an owned buffer,
//     which starts at minSize, is compacted as it goes, and grows up to maxSize.
//     A token that does not fit yields [ErrTooLong].
//
// The mark is the key invariant: while a token is held, the bytes from the mark
// to the cursor survive refills, so take can return them as one contiguous
// slice. Everything before the mark may be discarded.
//
//so:promote
type buffer struct {
	alloc   mem.Allocator
	rd      io.Reader // nil in fixed mode
	buf     []byte
	cur     int   // current cursor position
	end     int   // end of valid data
	mark    int   // start of the held token; <= cur
	minSize int   // initial buffer size
	maxSize int   // largest the buffer may grow to
	holding bool  // a token is being scanned
	err     error // sticky refill error (read failure or ErrTooLong)
}

// fixedBuffer returns a buffer over a complete in-memory document. The document
// is borrowed, not copied; alloc is only used for the decoder's scratch.
func fixedBuffer(alloc mem.Allocator, doc []byte) buffer {
	return buffer{alloc: alloc, buf: doc, end: len(doc)}
}

// streamBuffer returns a buffer that pulls from rd into a buffer of its own:
// minSize bytes, allocated on first use and doubled as tokens demand it, up to
// maxSize. A token that does not fit yields [ErrTooLong].
//
// maxSize is max(minSize, maxTok+1): a maxTok-byte token needs one extra byte
// to detect its end, but the buffer is never shrunk below minSize to enforce a
// smaller maxTok.
func streamBuffer(alloc mem.Allocator, rd io.Reader, minSize, maxTok int) buffer {
	maxTok = max(maxTok, minBufSize)
	minSize = max(minSize, minBufSize)
	maxSize := max(minSize, maxTok+1) // +1 to confirm the token is complete
	if minSize >= maxTok {
		minSize = maxSize // large enough for any token, so it never grows
	}
	return buffer{alloc: alloc, rd: rd, minSize: minSize, maxSize: maxSize}
}

// peek returns the byte at the cursor as an int in 0..255, or -1 at EOF,
// refilling first if the cursor has reached the end of the valid data.
func (b *buffer) peek() int {
	if b.cur < b.end {
		return int(b.buf[b.cur])
	}
	if !b.fill() {
		return -1
	}
	return int(b.buf[b.cur])
}

// advance consumes the current byte. It is only called after a successful peek.
func (b *buffer) advance() {
	b.cur++
}

// hold anchors the start of a token at the cursor.
func (b *buffer) hold() {
	b.mark = b.cur
	b.holding = true
}

// take ends the held token and returns its bytes: the slice from the mark
// to the cursor. The result is a view into the buffer, valid until the next
// refill (i.e. until the next peek that needs more data).
func (b *buffer) take() []byte {
	tok := b.buf[b.mark:b.cur]
	b.holding = false
	b.mark = b.cur
	return tok
}

// fill makes at least one more byte available past the cursor.
// Returns false at EOF or on error. In fixed mode it's a no-op,
// always returning false.
func (b *buffer) fill() bool {
	if b.err != nil || b.rd == nil {
		return false
	}

	// While a token is held, keep it. Otherwise, drop everything already consumed.
	keep := b.cur
	if b.holding {
		keep = b.mark
	}

	if b.buf == nil {
		b.buf = mem.AllocSlice[byte](b.alloc, b.minSize, b.minSize)
	}

	// Compact: slide the live bytes to the front.
	if keep > 0 {
		copy(b.buf, b.buf[keep:b.end])
		b.end -= keep
		b.cur -= keep
		b.mark -= keep
	}

	// Grow if compaction did not free any room (a token filling the buffer).
	// A buffer that started at its maximum fails here instead of growing.
	if b.end == len(b.buf) {
		if len(b.buf) >= b.maxSize {
			b.err = ErrTooLong
			return false
		}
		newSize := min(len(b.buf)*2, b.maxSize)
		newBuf := mem.AllocSlice[byte](b.alloc, newSize, newSize)
		copy(newBuf, b.buf[0:b.end])
		mem.FreeSlice(b.alloc, b.buf)
		b.buf = newBuf
	}

	// Read, tolerating a bounded number of empty reads.
	for range maxEmptyRead {
		n, err := b.rd.Read(b.buf[b.end:])
		if n > 0 {
			b.end += n
			return true
		}
		if err != nil {
			b.err = err
			return false
		}
	}
	b.err = io.ErrNoProgress
	return false
}

// free releases an owned buffer. In fixed mode it's a no-op.
func (b *buffer) free() {
	if !b.owned() {
		return
	}
	mem.FreeSlice(b.alloc, b.buf)
	b.buf = nil
	b.rd = nil
	b.cur, b.end, b.mark = 0, 0, 0
	b.holding = false
}

// owned reports whether the bytes take returns belong to the buffer. In
// streaming mode they do, so the decoder may write to them. In fixed mode
// they are the caller's document, which it must leave untouched.
func (b *buffer) owned() bool {
	return b.rd != nil
}
