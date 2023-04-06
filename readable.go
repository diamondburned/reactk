package reactk

import (
	"sync"

	"libdb.so/reactk/internal/syncg"
)

// Readable describes reactive stores that can be read from and subscribed to.
type Readable[T any] interface {
	Subscriber[T]
	Getter[T]
}

// NewReadable creates a new readable store. The store can only be set using
// initFn, which is called when the store is first subscribed to. If initFn
// returns a non-nil UnsubscribeFunc, it will be called when the store has no
// more subscribers.
//
// Note that if initFn is nil, then it is assumed that the store will never
// update and the returning implementation will be a static store.
func NewReadable[T any](initial T, initFn ReadableInitFunc[T]) Readable[T] {
	if initFn == nil {
		return static[T]{initial}
	}
	return newReadable(initial, initFn)
}

type static[T any] struct{ v T }

func (r static[T]) Subscribe(fn func(T)) UnsubscribeFunc {
	fn(r.v)
	return func() {}
}

func (r static[T]) Get() T {
	return r.v
}

// ReadableInitFunc is a function that initializes a Readable store.
// See NewReadable for more information.
type ReadableInitFunc[T any] func(set func(T)) UnsubscribeFunc

type readable[T any] struct {
	subs   syncg.Map[*callback[T], struct{}]
	init   syncg.AtomicValue[*initOnce[T]]
	initFn ReadableInitFunc[T]
}

type initOnce[T any] struct {
	once  sync.Once
	value syncg.AtomicValue[T]
	unsub UnsubscribeFunc
}

func (init *initOnce[T]) set(v T, r *readable[T]) {
	init.value.Store(v)
	enforcego(func() {
		v, _ := init.value.Load()
		r.subs.Range(func(cb *callback[T], _ struct{}) bool {
			cb.fn(v)
			return true
		})
	})
}

func (init *initOnce[T]) do(r *readable[T]) bool {
	var ran bool
	init.once.Do(func() {
		ran = true
		v := r.initFn(func(v T) { init.set(v, r) })
		// As long as we call unsub only after init.once finishes, we should be
		// good. I think.
		if v != nil {
			var unsubOnce sync.Once
			init.unsub = func() { unsubOnce.Do(v) }
		} else {
			init.unsub = func() {}
		}
	})
	return ran
}

func newReadable[T any](initial T, initFn ReadableInitFunc[T]) *readable[T] {
	var init syncg.AtomicValue[*initOnce[T]]
	init.Store(&initOnce[T]{})

	return &readable[T]{
		initFn: initFn,
		init:   init,
	}
}

func (r *readable[T]) Subscribe(fn func(T)) UnsubscribeFunc {
	init, ok := r.init.Load()
	if !ok {
		panic("bug: cannot subscribe: r.init is nil")
	}

	callback := &callback[T]{fn}
	r.subs.Store(callback, struct{}{})

	enforcego(func() {
		// NOTE: Subscribe has a weird behavior where if Subscribe is first
		// called, the given callback may be called twice. This guard prevents
		// this behavior.
		if !init.do(r) {
			v, _ := init.value.Load()
			fn(v)
		}
	})

	return func() {
		r.subs.Delete(callback)
		if r.subs.IsEmpty() {
			r.init.CompareAndSwap(init, &initOnce[T]{})
			// Successfully cleared r.init while guaranteeing that there's only
			// one winner. Finish our cleanup.
			enforcego(func() {
				// Ensure that init is called before we call the unsubscriber.
				init.do(r)
				init.unsub()
			})
		}
	}
}

func (r *readable[T]) Get() T {
	init, ok := r.init.Load()
	if !ok {
		panic("bug: cannot subscribe: r.init is nil")
	}

	v, _ := init.value.Load()
	return v
}
