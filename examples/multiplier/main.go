package main

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"libdb.so/reactk/greactk"
)

func main() {
	number := greactk.NewWritable(0.0)
	multiplier := greactk.Derive[float64](number, func(n float64) float64 {
		return n * 2
	})

	app := gtk.NewApplication("com.github.diamondburned.greactk.example", 0)
	app.ConnectActivate(func() {
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
	})
	app.Run(nil)
}
