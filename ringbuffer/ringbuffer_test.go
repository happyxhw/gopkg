package ringbuffer

import (
	"fmt"
	"testing"
	"time"
)

func TestRingBuffer_Write_Read(t *testing.T) {
	size := 20
	r := NewRingBuffer(size)
	val := make([]int, 100)
	var i int
	for {
		if !r.IsFull() {
			i++
			err := r.WriteOne(i)
			if err != nil {
				t.Fatal(err)
			}
			time.Sleep(time.Millisecond * 100)
			continue
		} else {
			n, err := r.Read(val)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(val[:n])
		}
	}
}

func BenchmarkRingBuffer_Read1(b *testing.B) {
	size := 200
	N := 100000000
	r := NewRingBuffer(size)
	val := make([]int, size)
	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			var timeout bool
			select {
			case <-t.C:
				timeout = true
			default:
				_ = r.WriteOne(j)
			}
			if (timeout || r.isFull) && !r.IsEmpty() {
				_, _ = r.Read(val)
			}
		}
	}
}

// x4 performance than ringbuffer
func BenchmarkRingBuffer_Read2(b *testing.B) {
	size := 200
	N := 100000000
	r := make([]int, 0, size)
	val := make([]int, size)
	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			var timeout bool
			select {
			case <-t.C:
				timeout = true
			default:
				if len(r) < size {
					r = append(r, j)
				}
			}
			if (timeout || len(r) == size) && len(r) > 0 {
				copy(val, r)
				r = r[:0]
			}
		}
	}
}
