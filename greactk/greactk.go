package greactk

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"libdb.so/reactk"
)

var DefaultPriority = glib.PriorityDefault

func init() {
	reactk.EnforceGoroutine = false
}

// Widget describes a widget that can be used to bind to a store. It is an
// interface that both GTK3 and GTK4 widgets implement.
type Widget interface {
	ConnectRealize(f func()) glib.SignalHandle
	ConnectUnrealize(f func()) glib.SignalHandle
}

// These methods above are type asserted. Source: trust me bro, I didn't just
// made it the fuck up.

// Readable extends reactk.Readable with widget-specific methods.
type Readable[T any] interface {
	reactk.Readable[T]
	// SubscribeWidget subscribes the given function to the store. The function
	// will be called in the GLib main thread whenever the store's value
	// changes.
	SubscribeWidget(Widget, func(T))
}

// Subscribe subscribes to the given store. The given callback will be called in
// the GLib main loop.
func Subscribe[T any](readable reactk.Readable[T], f func(T)) reactk.UnsubscribeFunc {
	return readable.Subscribe(func(v T) {
		glib.IdleAddPriority(DefaultPriority, func() { f(v) })
	})
}

// SubscribeWidget implements Readable's SubscribeWidget.
func SubscribeWidget[T any](readable reactk.Readable[T], widget Widget, f func(T)) {
	var unsub reactk.UnsubscribeFunc
	widget.ConnectRealize(func() {
		// Note: the store is binded for the lifetime of the widget, which is
		// bad.
		unsub = readable.Subscribe(f)
	})
	widget.ConnectUnrealize(func() {
		unsub()
	})
}

type readable[T any] struct{ reactk.Readable[T] }

func (w readable[T]) Subscribe(f func(T)) reactk.UnsubscribeFunc {
	return Subscribe(w.Readable, f)
}

func (w readable[T]) SubscribeWidget(widget Widget, f func(T)) {
	SubscribeWidget(w.Readable, widget, f)
}

// NewReadable creates a new readable store. The store's subscribers will be
// called in the GLib main loop.
func NewReadable[T any](initial T, initFn reactk.ReadableInitFunc[T]) Readable[T] {
	return readable[T]{reactk.NewReadable(initial, initFn)}
}

// Writable extends reactk.Writable with widget-specific methods.
type Writable[T any] interface {
	Readable[T]
	reactk.Writable[T]
}

type writable[T any] struct{ reactk.Writable[T] }

func (w writable[T]) Subscribe(f func(T)) reactk.UnsubscribeFunc {
	return Subscribe[T](w.Writable, f)
}

func (w writable[T]) SubscribeWidget(widget Widget, f func(T)) {
	SubscribeWidget[T](w.Writable, widget, f)
}

// NewWritable creates a new writable store. The store's subscribers will be
// called in the GLib main loop.
func NewWritable[T any](initial T) reactk.Writable[T] {
	return writable[T]{reactk.NewWritable(initial)}
}

// Derive wraps reactk.Derive like NewReadable.
func Derive[SourceT, DestinationT any](store reactk.Readable[SourceT], derive func(SourceT) DestinationT) Readable[DestinationT] {
	return readable[DestinationT]{reactk.Derive(store, derive)}
}

// DeriveMultiple wraps reactk.DeriveMultiple like NewReadable.
func DeriveMultiple[T any](stores []reactk.AnyReadable, deriver any) Readable[T] {
	return readable[T]{reactk.DeriveMultiple[T](stores, deriver)}
}
