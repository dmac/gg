// +build js

package main

import (
	"log"
	"time"

	"github.com/dmac/gg"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"honnef.co/go/js/dom"
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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	if err := gg.Init(CanvasWidth, CanvasHeight, gl); err != nil {
		log.Fatal(err)
	}

	loadTextures()

	tetris := NewTetris(CanvasWidth, CanvasHeight)
	dom.GetWindow().Document().AddEventListener("keypress", false, tetris.handleKeyPress)

	for {
		gl.ClearColor(0.5, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		tetris.Draw()

		time.Sleep(16 * time.Millisecond)
	}
}

func (t *Tetris) handleKeyPress(e dom.Event) {
	event, ok := e.(*dom.KeyboardEvent)
	if !ok {
		return
	}
	switch event.CharCode {
	case 119:
		t.HandleInput(inputUp)
	case 100:
		t.HandleInput(inputRight)
	case 115:
		t.HandleInput(inputDown)
	case 97:
		t.HandleInput(inputLeft)
	case 32:
		t.HandleInput(inputSpace)
	}
}

var textures map[string]*gg.Texture

func loadTextures() {
	textures = make(map[string]*gg.Texture)
	textures["orange"] = gg.NewTextureFromImage(openImage("orange.png"))
}

func openImage(path string) *js.Object {
	img := js.Global.Get("Image").New()
	img.Set("src", path)
	img.Set("crossOrigin", "")
	c := make(chan struct{})
	img.Call("addEventListener", "load", func() { close(c) }, false)
	<-c
	return img
}
