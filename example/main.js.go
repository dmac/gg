// +build js

package main

import (
	"log"
	"time"

	"github.com/dmac/gg"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

func main() {
	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	document.Get("body").Call("appendChild", canvas)
	canvas.Call("setAttribute", "id", "canvas")
	canvas.Call("setAttribute", "width", 640)
	canvas.Call("setAttribute", "height", 480)

	attrs := webgl.DefaultAttributes()
	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		log.Fatal(err)
	}

	if err := gg.Init(gl); err != nil {
		log.Fatal(err)
	}

	triangle := gg.NewPoly([]gg.Vec2{
		{150, 50},
		{50, 50},
		{50, 150},
	})

	for {
		gl.ClearColor(0.5, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		triangle.Draw()

		time.Sleep(16 * time.Millisecond)
		break
	}
}
