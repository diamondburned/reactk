package syncg

import "sync/atomic"

// AtomicValue is a generic wrapper around atomic.Value.
type AtomicValue[T any] atomic.Value

// Load returns the value stored in the AtomicValue.
func (v *AtomicValue[T]) Load() (T, bool) {
	val, ok := (*atomic.Value)(v).Load().(T)
	if ok {
		return val, true
	}
	var z T
	return z, false
}

// Store stores the given value in the AtomicValue.
func (v *AtomicValue[T]) Store(val T) {
	(*atomic.Value)(v).Store(val)
}

// Swap swaps the value stored in the AtomicValue with the given value and
// returns the old value.
func (v *AtomicValue[T]) Swap(new T) T {
	return (*atomic.Value)(v).Swap(new).(T)
}

// CompareAndSwap compares the value stored in the AtomicValue with the given
// value and, if they are equal, stores the new value and returns true. If they
// are not equal, it returns false.
func (v *AtomicValue[T]) CompareAndSwap(old, new T) bool {
	return (*atomic.Value)(v).CompareAndSwap(old, new)
}

// AtomicBool is an atomic wrapper for bool values.
type AtomicBool uint32

// Load returns the value stored in the Bool.
func (b *AtomicBool) Load() bool {
	return atomic.LoadUint32((*uint32)(b)) != 0
}

// Store stores the given value in the Bool.
func (b *AtomicBool) Store(val bool) {
	if val {
		atomic.StoreUint32((*uint32)(b), 1)
	} else {
		atomic.StoreUint32((*uint32)(b), 0)
	}
}

// Swap swaps the value stored in the Bool with the given value and returns the
// old value.
func (b *AtomicBool) Swap(new bool) bool {
	if new {
		return atomic.SwapUint32((*uint32)(b), 1) != 0
	}
	return atomic.SwapUint32((*uint32)(b), 0) != 0
}

// CompareAndSwap compares the value stored in the Bool with the given value and,
// if they are equal, stores the new value and returns true. If they are not
// equal, it returns false.
func (b *AtomicBool) CompareAndSwap(old, new bool) bool {
	var oldn, newn uint32
	if old {
		oldn = 1
	}
	if new {
		newn = 1
	}
	return atomic.CompareAndSwapUint32((*uint32)(b), oldn, newn)
}
