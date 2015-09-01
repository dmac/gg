// +build js

package main

import (
	"log"
	"time"

	"github.com/dmac/gg"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

const CanvasWidth = 640
const CanvasHeight = 480

func main() {
	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	document.Get("body").Call("appendChild", canvas)
	canvas.Call("setAttribute", "id", "canvas")
	canvas.Call("setAttribute", "width", CanvasWidth)
	canvas.Call("setAttribute", "height", CanvasHeight)

	attrs := webgl.DefaultAttributes()
	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		log.Fatal(err)
	}

	if err := gg.Init(CanvasWidth, CanvasHeight, gl); err != nil {
		log.Fatal(err)
	}

	triangle := gg.NewPoly([]gg.Vec2{
		{320, 100},
		{240, 200},
		{400, 200},
	})
	triangle.SetColor(1, 0, 1, 1)

	img1 := NewImageFromFile("test.png")
	img2 := NewImageFromFile("test2.png")
	tex1 := gg.NewTextureFromImage(img1)
	tex2 := gg.NewTextureFromImage(img2)
	spr1 := gg.NewSpriteFromTexture(tex1)
	spr2 := gg.NewSpriteFromTexture(tex2)
	spr1.SetPosition(240, 230)
	spr2.SetPosition(340, 230)

	for {
		gl.ClearColor(0.5, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		triangle.Draw()
		spr1.Draw()
		spr2.Draw()

		time.Sleep(16 * time.Millisecond)
		break
	}
}

func NewImageFromFile(filename string) *js.Object {
	img := js.Global.Get("Image").New()
	img.Set("src", filename)
	img.Set("crossOrigin", "")
	c := make(chan struct{})
	img.Call("addEventListener", "load", func() { close(c) }, false)
	<-c
	return img
}