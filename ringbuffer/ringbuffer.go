package ringbuffer

import (
	"errors"
	"sync"
)

var (
	ErrTooManyDataToWrite = errors.New("too many buf to write")
	ErrIsFull             = errors.New("ring buffer is full")
	ErrIsEmpty            = errors.New("ring buffer is empty")
)

// RingBuffer ring buffer
type RingBuffer struct {
	m      sync.RWMutex
	buf    []int
	size   int
	isFull bool // 区分空buf和满buf两种情况，两种情况下 r == w

	r int // 读指针
	w int // 写指针
}

// NewRingBuffer return an instance of RingBuffer
func NewRingBuffer(size int) *RingBuffer {
	r := RingBuffer{
		buf:  make([]int, size),
		size: size,
	}
	return &r
}

// IsFull is full
func (r *RingBuffer) IsFull() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.isFull
}

// IsEmpty is empty
func (r *RingBuffer) IsEmpty() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.r == r.w && !r.isFull
}

// Read from buf
func (r *RingBuffer) Read(p []int) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	r.m.Lock()
	n, err := r.read(p)
	r.m.Unlock()
	return n, err
}

// Write to buf
func (r *RingBuffer) Write(p []int) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	r.m.Lock()
	n, err := r.write(p)
	r.m.Unlock()
	return n, err
}

// read all data from buf
func (r *RingBuffer) read(p []int) (n int, err error) {
	if r.w == r.r && !r.isFull {
		return 0, ErrIsEmpty
	}
	defer func() {
		r.isFull = false
		r.r = (r.r + n) % r.size
	}()

	if r.w > r.r {
		n = r.w - r.r
		if n > len(p) {
			n = len(p)
		}
		copy(p, r.buf[r.r:r.r+n])
		return
	}

	n = r.size - r.r + r.w
	if n > len(p) {
		n = len(p)
	}
	if r.r+n < r.size {
		copy(p, r.buf[:n])
		return
	}

	s1 := r.size - r.r
	copy(p, r.buf[r.r:r.size])
	copy(p[s1:], r.buf[:(n-s1)])

	return n, err
}

// write data to buf
func (r *RingBuffer) write(p []int) (n int, err error) {
	if r.isFull {
		return 0, ErrIsFull
	}

	var avail int
	if r.w >= r.r {
		avail = r.size - r.w + r.r
	} else {
		avail = r.r - r.w
	}

	if len(p) > avail {
		err = ErrTooManyDataToWrite
		p = p[:avail]
	}

	n = len(p)
	if r.w < r.r {
		copy(r.buf[r.w:], p)
		r.w += n
	} else {
		if r.w+n <= r.size {
			copy(r.buf[r.w:], p)
			r.w = (r.w + n) % r.size
		} else {
			s1 := r.size - r.w
			copy(r.buf[r.w:], p[:s1])
			copy(r.buf[0:], p[s1:])
			r.w = n - s1
		}
	}
	if r.w == r.r {
		r.isFull = true
	}

	return n, err
}

// ReadOne read next item
func (r *RingBuffer) ReadOne() (val int, err error) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.r == r.w && !r.isFull {
		return 0, ErrIsEmpty
	}
	val = r.buf[r.r]
	r.r = (r.r + 1) % r.size
	r.isFull = false
	return
}

// WriteOne one to buffer
func (r *RingBuffer) WriteOne(val int) error {
	r.m.Lock()
	defer r.m.Unlock()
	if r.isFull {
		return ErrIsFull
	}
	r.buf[r.w] = val
	r.w = (r.w + 1) % r.size
	if r.w == r.r {
		r.isFull = true
	}
	return nil
}
