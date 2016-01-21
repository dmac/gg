// +build js

package main

import (
	"log"
	"time"

	"github.com/dmac/gg"
	ggwebgl "github.com/dmac/gg/webgl/gg"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
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

	const vertShader = `

uniform mat4 proj;
attribute vec3 vertex_position;
attribute vec2 vertex_texture;
varying highp vec2 texture_coordinates;

void main() {
	gl_Position = proj * vec4(vertex_position, 1);
	texture_coordinates = vertex_texture;
}
`

	const fragShader = `#version 100

uniform sampler2D tex_loc;
varying highp vec2 texture_coordinates;

void main() {
	gl_FragColor = texture2D(tex_loc, texture_coordinates);
}
`

	texture := newImageTexture("sq.png")
	scene, err := NewScene(vertShader, fragShader, texture)
	if err != nil {
		log.Fatal(err)
	}

	for {
		scene.Draw()
		time.Sleep(16 * time.Millisecond)
	}
}

// TODO(dmac) Move into gg helpers package
func newImageTexture(path string) *gg.Texture {
	img := js.Global.Get("Image").New()
	img.Set("src", path)
	img.Set("crossOrigin", "")
	c := make(chan struct{})
	img.Call("addEventListener", "load", func() { close(c) }, false)
	<-c

	tex := gg.CreateTexture()
	gg.ActiveTexture(gg.TEXTURE0)
	gg.BindTexture(gg.TEXTURE_2D, tex)
	gg.TexImage2D(
		gg.TEXTURE_2D,
		0,
		gg.RGBA,
		-1, -1, -1, // These args are ignored in the webgl backend.
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
