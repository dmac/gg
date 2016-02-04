// +build !js

package main

import (
	"log"
	"runtime"

	_ "github.com/dmac/gg/v2.1"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func main() {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Triangle Demo", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	const vertShader = `#version 120
uniform mat4 proj;
attribute vec3 vertex_position;

void main() {
	gl_Position = proj * vec4(vertex_position, 1);
}
`

	const fragShader = `#version 120
uniform vec4 color;

void main() {
	gl_FragColor = color;
}
`
	scene, err := NewScene(vertShader, fragShader)
	if err != nil {
		log.Fatal(err)
	}

	for !window.ShouldClose() {
		scene.Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
