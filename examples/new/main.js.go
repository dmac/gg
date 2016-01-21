// +build js

package main

import (
	"log"
	"time"

	webgg "github.com/dmac/gg/webgl/gg"
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
	webgg.Init(gl)

const vertShader = `#version 100

uniform mat4 proj;
attribute vec3 vertex_position;

void main() {
	gl_Position = proj * vec4(vertex_position, 1);
}
`

const fragShader = `#version 100

uniform highp vec4 color;

void main() {
	gl_FragColor = color;
}
`
	game, err := NewGame(vertShader, fragShader)
	if err != nil {
		log.Fatal(err)
	}

	for {
		game.Draw()
		time.Sleep(16 * time.Millisecond)
	}
}
