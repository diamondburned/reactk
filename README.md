# reactk

Package `reactk` implements a small library for reactive GUI programming. It is
inspired by Svelte's stores, without the goodies.

It contains package `greactk` which contains wrapper for GTK widgets.

## Examples

```go
number := greactk.NewWritable(0.0)
multiplier := greactk.Derive[float64](number, func(n float64) float64 {
	return n * 2
})

spin := gtk.NewSpinButtonWithRange(-100, 100, 0.1)
spin.SetValue(number.Get())
spin.ConnectChanged(func() { number.Set(spin.Value()) })

label := gtk.NewLabel("")
label.SetJustify(gtk.JustifyLeft)
multiplier.SubscribeWidget(label, func(n float64) {
	label.SetText(fmt.Sprintf(
		"x1: %f\nx2: %f",
		number.Get(), n,
	))
})

box := gtk.NewBox(gtk.OrientationVertical, 0)
box.Append(spin)
box.Append(label)

win := gtk.NewApplicationWindow(app)
win.SetTitle("greactk example")
win.SetChild(box)
win.Show()
```

See the [./examples](examples) directory for more.
