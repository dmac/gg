// +build !js

package main

import (
	"image/png"
	"log"
	"os"
	"runtime"

	"github.com/dmac/gg"
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

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Texture Demo", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	const vertShader = `#version 120

uniform mat4 proj;
attribute vec3 vertex_position;
attribute vec2 vertex_texture;
varying vec2 texture_coordinates;

void main() {
	gl_Position = proj * vec4(vertex_position, 1);
	texture_coordinates = vertex_texture;
}
`

	const fragShader = `#version 120

uniform sampler2D tex_loc;
varying vec2 texture_coordinates;

void main() {
	gl_FragColor = texture2D(tex_loc, texture_coordinates);
}
`

	texture, err := newImageTexture("sq.png")
	if err != nil {
		log.Fatal(err)
	}
	scene, err := NewScene(vertShader, fragShader, texture)
	if err != nil {
		log.Fatal(err)
	}

	for !window.ShouldClose() {
		scene.Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

// TODO(dmac) Move into gg helpers package
func newImageTexture(filename string) (*gg.Texture, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	var buf []byte
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			buf = append(buf, byte(r/256))
			buf = append(buf, byte(g/256))
			buf = append(buf, byte(b/256))
			buf = append(buf, byte(a/256))
		}
	}
	tex := gg.CreateTexture()
	gg.ActiveTexture(gg.TEXTURE0)
	gg.Enable(gg.TEXTURE_2D)
	gg.BindTexture(gg.TEXTURE_2D, tex)
	gg.TexImage2D(
		gg.TEXTURE_2D,
		0,
		gg.RGBA,
		img.Bounds().Dx(),
		img.Bounds().Dy(),
		0,
		gg.RGBA,
		gg.UNSIGNED_BYTE,
		buf,
	)
	gg.TexParameteri(gg.TEXTURE_2D, gg.TEXTURE_WRAP_S, gg.CLAMP_TO_EDGE)
	gg.TexParameteri(gg.TEXTURE_2D, gg.TEXTURE_WRAP_T, gg.CLAMP_TO_EDGE)
	gg.TexParameteri(gg.TEXTURE_2D, gg.TEXTURE_MAG_FILTER, gg.LINEAR)
	gg.TexParameteri(gg.TEXTURE_2D, gg.TEXTURE_MIN_FILTER, gg.LINEAR)

	return tex, nil
}
