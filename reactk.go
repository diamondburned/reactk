// Package reactk implements a small library for reactive GUI programming. It is
// inspired by Svelte's stores, without the goodies.
package reactk

// EnforceGoroutine determines whether or not reactk should enforce non-blocking
// goroutine behaviors for Set and Subscribe. This is useful in most Go
// applications, but in GTK applications, it is wasteful, since GTK uses its own
// main thread which already enforces non-blocking behavior.
//
// Most applications should leave this as true. If greactk is used, it will
// automatically set this to false.
var EnforceGoroutine = true

func enforcego(f func()) {
	if EnforceGoroutine {
		go f()
	} else {
		f()
	}
}

// UnsubscribeFunc is a function that can be called to unsubscribe from a store.
type UnsubscribeFunc func()

// Subscriber describes reactive stores that can be subscribed to.
type Subscriber[T any] interface {
	// Subscribe subscribes the given function to the store. The function will
	// be called whenever the store's value changes. The function will be called
	// in an unspecified goroutine; manual synchronization is required if
	// needed.
	Subscribe(fn func(T)) UnsubscribeFunc
}

// Getter describes reactive stores that can be read from.
type Getter[T any] interface {
	// Get returns the current value of the store.
	Get() T
}

// Setter describes reactive stores that can be written to.
type Setter[T any] interface {
	// Set sets the value of the store. It will trigger all subscribers.
	Set(T)
}
