package reactk

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestDerive(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		res := make(chan int, 100)

		writable := NewWritable(0)
		multiply := Derive[int](writable, func(v int) int { return v * 2 })
		multiply.Subscribe(func(v int) { res <- v })

		writable.Set(1)
		writable.Set(2)

		expects := map[int]int{
			0: 0,
			2: 0,
			4: 0,
		}

		for range expects {
			v := <-res
			_, ok := expects[v]
			assert.True(t, ok)
			if ok {
				expects[v]++
			}
		}

		// There's actually no guarantee that everything is called exactly once!
		// Set will opportunistically call the subscribers with the latest
		// value.
		for k, v := range expects {
			t.Logf("expect %d to be called %d times", k, v)
		}
	})
}
