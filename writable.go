package reactk

import "libdb.so/reactk/internal/syncg"

// Writable describes reactive stores that can be read from, written to, and
// subscribed to.
type Writable[T any] interface {
	Readable[T]
	Setter[T]
}

type writable[T any] struct {
	value syncg.AtomicValue[T]
	subs  syncg.Map[*callback[T], struct{}]
}

var _ Writable[any] = (*writable[any])(nil)

type callback[T any] struct {
	fn func(T)
}

func (w *writable[T]) Subscribe(fn func(T)) UnsubscribeFunc {
	cb := &callback[T]{fn}
	w.subs.Store(cb, struct{}{})

	v, _ := w.value.Load()
	enforcego(func() { fn(v) })

	return func() { w.subs.Delete(cb) }
}

func (w *writable[T]) Get() T {
	v, _ := w.value.Load()
	return v
}

func (w *writable[T]) Set(v T) {
	w.value.Store(v)
	enforcego(func() {
		v, _ := w.value.Load()
		w.subs.Range(func(k *callback[T], _ struct{}) bool {
			k.fn(v)
			return true
		})
	})
}

// NewWritable creates a new writable store.
func NewWritable[T any](initial T) Writable[T] {
	return newWritable(initial)
}

func newWritable[T any](initial T) *writable[T] {
	var value syncg.AtomicValue[T]
	value.Store(initial)

	return &writable[T]{
		value: value,
	}
}

type readonly[T any] struct{ w Readable[T] }

// Readonly returns a readonly version of the store. Store may be a
// Writable or a Readable. The methods will be overridden to prevent
// type assertions over the store.
func Readonly[T any](store Readable[T]) Readable[T] {
	return readonly[T]{store}
}

func (r readonly[T]) Subscribe(fn func(T)) UnsubscribeFunc {
	return r.w.Subscribe(fn)
}

func (r readonly[T]) Get() T {
	return r.w.Get()
}
