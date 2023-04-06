package reactk

import (
	"fmt"
	"reflect"
	"runtime"
)

// AnyReadable is a Readable that can hold any value. It is used to perform
// reflect magic.
type AnyReadable struct {
	r interface {
		Subscribe(func(any)) UnsubscribeFunc
		Get() any
	}
	t reflect.Type
}

type anyAdapter[T any] struct {
	store Readable[T]
}

func (a anyAdapter[T]) Subscribe(fn func(any)) UnsubscribeFunc {
	return a.store.Subscribe(func(v T) { fn(v) })
}

func (a anyAdapter[T]) Get() any {
	return a.store.Get()
}

// AdaptReadable adapts a Readable[T] to AnyReadable. It is used to perform
// reflect magic.
func AdaptReadable[T any](readable Readable[T]) AnyReadable {
	var z T
	return AnyReadable{
		r: anyAdapter[T]{readable},
		t: reflect.TypeOf(z),
	}
}

type derived[T any] struct {
	main    Writable[T]
	stores  []AnyReadable
	unsubs  []UnsubscribeFunc
	deriver reflect.Value
}

// Derive creates a new derived store that is updated whenever the given store
// is updated. The derived store will be updated with the result of calling fn
// with the value of the store.
func Derive[SourceT, DestinationT any](store Readable[SourceT], derive func(SourceT) DestinationT) Readable[DestinationT] {
	return DeriveMultiple[DestinationT]([]AnyReadable{AdaptReadable(store)}, derive)
}

// DeriveMultiple creates a new derived store that is updated whenever any of
// the stores are updated. The derived store will be updated with the result of
// calling fn with the values of the stores.
//
// Most cases should use Derive instead, which is type-safe but only allows
// deriving from a single store.
func DeriveMultiple[T any](stores []AnyReadable, deriver any) Readable[T] {
	var z T
	destType := reflect.TypeOf(z)
	deriverType := reflect.TypeOf(deriver)
	if deriverType.Kind() != reflect.Func {
		panic("reactk: deriver is not a function")
	}
	if len(stores) != deriverType.NumIn() {
		panic("reactk: number of stores does not match number of deriver inputs")
	}
	if deriverType.NumOut() != 1 {
		panic("reactk: deriver must return exactly one value")
	}
	for i, store := range stores {
		fnIn := deriverType.In(i)
		if !store.t.AssignableTo(fnIn) {
			panic(fmt.Sprintf(
				"reactk: store %s is not assignable to deriver input %d (%s)",
				store.t, i, fnIn))
		}
	}
	if !deriverType.Out(0).AssignableTo(destType) {
		panic(fmt.Sprintf(
			"reactk: deriver output is not assignable to %s", destType))
	}

	d := &derived[T]{
		main:    NewWritable(z),
		stores:  stores,
		deriver: reflect.ValueOf(deriver),
	}

	runtime.SetFinalizer(d, func(d *derived[T]) {
		for _, unsub := range d.unsubs {
			unsub()
		}
	})

	doneCh := make(chan struct{})

	d.unsubs = make([]UnsubscribeFunc, len(stores))
	for i, store := range stores {
		i := i
		d.unsubs[i] = store.r.Subscribe(func(any) {
			select {
			case <-doneCh:
				d.derive()
			default:
			}
		})
	}

	// Ensure that all subscriptions are created before the first update.
	d.derive()
	close(doneCh)

	return d
}

func (d *derived[T]) Subscribe(fn func(T)) UnsubscribeFunc {
	return d.main.Subscribe(fn)
}

func (d *derived[T]) Get() T {
	return d.main.Get()
}

func (d *derived[T]) derive() {
	values := make([]reflect.Value, len(d.stores))
	for i, store := range d.stores {
		v := store.r.Get()
		values[i] = reflect.ValueOf(v)
	}

	result := d.deriver.Call(values)[0].Interface()
	d.main.Set(result.(T))
}

// func Derive()
//
// func newDerived[T any](stores []Readable[T])
