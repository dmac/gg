// +build js

package main

import (
	"log"
	"time"

	"github.com/dmac/gg"
	ggwebgl "github.com/dmac/gg/webgl/gg"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"honnef.co/go/js/dom"
)

func main() {
	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	document.Get("body").Call("appendChild", canvas)
	canvas.Call("setAttribute", "id", "canvas")
	canvas.Call("setAttribute", "width", WindowWidth)
	canvas.Call("setAttribute", "height", WindowHeight)

	attrs := webgl.DefaultAttributes()
	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		log.Fatal(err)
	}
	ggwebgl.Init(gl)

	const vshader = `#version 100

uniform mat4 proj, model;
attribute vec3 vertex_position;
attribute vec2 vertex_texture;
varying highp vec2 texture_coordinates;

void main() {
	gl_Position = proj * model * vec4(vertex_position, 1);
	texture_coordinates = vertex_texture;
}
`

	const fshader = `#version 100

uniform sampler2D tex_loc;
varying highp vec2 texture_coordinates;

void main() {
	gl_FragColor = texture2D(tex_loc, texture_coordinates);
}
`

	tetris, err := NewTetris(vshader, fshader)
	if err != nil {
		log.Fatal(err)
	}

	dom.GetWindow().Document().AddEventListener("keypress", false, tetris.handleKeyPress)

	for {
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

func LoadTextures() (map[string]*gg.Texture, error) {
	textures := make(map[string]*gg.Texture)
	textures["bg"] = NewTextureFromImage(openImage("bg.png"))
	textures["board"] = NewTextureFromImage(openImage("board.png"))
	textures["blue"] = NewTextureFromImage(openImage("blue.png"))
	textures["cyan"] = NewTextureFromImage(openImage("cyan.png"))
	textures["green"] = NewTextureFromImage(openImage("green.png"))
	textures["orange"] = NewTextureFromImage(openImage("orange.png"))
	textures["purple"] = NewTextureFromImage(openImage("purple.png"))
	textures["red"] = NewTextureFromImage(openImage("red.png"))
	textures["yellow"] = NewTextureFromImage(openImage("yellow.png"))
	return textures, nil
}

func NewTextureFromImage(img *js.Object) *gg.Texture {
	tex := gg.CreateTexture()
	gg.ActiveTexture(gg.TEXTURE0)
	gg.BindTexture(gg.TEXTURE_2D, tex)
	gg.TexImage2D(
		gg.TEXTURE_2D,
		0,
		gg.RGBA,
		-1, -1, -1, // ignored by gg/webgl
		gg.RGBA,
		gg.UNSIGNED_BYTE,
		img,
	)
	gg.TexParameteri(gg.TEXTURE_2D, gg.TEXTURE_WRAP_S, gg.CLAMP_TO_EDGE)
	gg.TexParameteri(gg.TEXTURE_2D, gg.TEXTURE_WRAP_T, gg.CLAMP_TO_EDGE)
	gg.TexParameteri(gg.TEXTURE_2D, gg.TEXTURE_MAG_FILTER, gg.LINEAR)
	gg.TexParameteri(gg.TEXTURE_2D, gg.TEXTURE_MIN_FILTER, gg.LINEAR)

	return tex
}

func openImage(path string) *js.Object {
	img := js.Global.Get("Image").New()
	img.Set("src", "images/"+path)
	img.Set("crossOrigin", "")
	c := make(chan struct{})
	img.Call("addEventListener", "load", func() { close(c) }, false)
	<-c
	return img
}
