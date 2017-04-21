package buffer

import "time"

type Buffer struct {
	index int
	ring  []time.Time
}

func New(size int) *Buffer {
	if size < 1 {
		size = 1
	}
	return &Buffer{
		index: 0,
		ring:  make([]time.Time, size),
	}
}

func (b *Buffer) AddTime(t time.Time) {
	b.ring[b.index] = t
	if b.index+1 == len(b.ring) {
		b.index = 0
	} else {
		b.index++
	}
}

func (b *Buffer) AddNow() {
	b.AddTime(time.Now())
}

func (b *Buffer) IsOkay(t time.Time, d time.Duration) bool {
	ct := b.ring[b.index]
	if ct.IsZero() {
		return true
	}
	return t.Sub(ct) > d
}
