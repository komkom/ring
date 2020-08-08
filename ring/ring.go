package ring

import (
	"fmt"
)

var ErrOverflow = fmt.Errorf(`overflow`)

type R struct {
	buffer   []byte
	size     int64
	writePos int64
	headPos  int64
}

func New(size int) *R {

	return &R{
		buffer: make([]byte, size),
		size:   int64(size),
	}
}

func (r *R) Len() int {

	if r.writePos < r.headPos {
		return int(r.size - r.headPos + r.writePos)
	}
	return int(r.writePos - r.headPos)
}

func (r *R) Write(p []byte) (n int, err error) {

	length := int64(len(p))

	wp := r.writePos
	head := r.headPos

	if wp < head && head-wp <= length {
		return 0, ErrOverflow
	}

	if wp >= head && head+r.size-wp <= length {
		return 0, ErrOverflow
	}

	if wp < head || r.size-wp >= length { // no split needed

		copy(r.buffer[wp:], p)
		r.writePos = (wp + length) % r.size

		return int(length), nil

	} else { // split needed not enough contiguous memory

		split := r.size - wp
		left := length - split

		copy(r.buffer[wp:], p[:split])
		copy(r.buffer, p[split:])

		r.writePos = left
		return int(length), nil
	}
}

func (r *R) Read(p []byte) (n int, err error) {

	length := int64(len(p))

	head := r.headPos
	write := r.writePos

	if head == write {
		return 0, nil
	}

	var size int64
	if head < write {

		size = write - head
		if length < size {
			size = length
		}
		copy(p, r.buffer[head:head+size])

	} else if r.size-head >= length {

		size = length
		copy(p, r.buffer[head:head+size])

	} else { // split needed

		size = r.size + write - head
		if length < size {
			size = length
		}

		right := r.size - head
		left := size - right

		copy(p, r.buffer[head:])
		copy(p[right:], r.buffer[:left])
	}

	r.headPos = (head + size) % r.size

	return int(size), nil
}
