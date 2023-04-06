package reactk

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestWritable(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		res := make(chan int, 1)

		w := NewWritable(0)
		w.Subscribe(func(v int) { res <- v })
		assert.Equal(t, 0, <-res)

		w.Set(1)
		assert.Equal(t, 1, <-res)

		v, ok := w.Get()

		assert.True(t, ok)
		assert.Equal(t, 1, v)
	})
}
