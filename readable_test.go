package reactk

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestReadable(t *testing.T) {
	t.Run("initer", func(t *testing.T) {
		var turn int

		res := make(chan int, 1)
		unsubbed := make(chan struct{}, 1)

		r := NewReadable(0, func(set func(int)) UnsubscribeFunc {
			turn++
			set(turn)
			return func() {
				turn *= 2
				unsubbed <- struct{}{}
			}
		})

		unsub1 := r.Subscribe(func(v int) { res <- v })
		assert.Equal(t, 1, <-res)

		unsub1()
		<-unsubbed

		_, ok1 := r.Get()
		assert.False(t, ok1)
		assert.Equal(t, 2, turn)

		unsub2 := r.Subscribe(func(v int) { res <- v })
		assert.Equal(t, 3, <-res)

		unsub2()
		<-unsubbed

		_, ok2 := r.Get()
		assert.False(t, ok2)
		assert.Equal(t, 6, turn)
	})
}
