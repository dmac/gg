// +build !js

package main

import (
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/dmac/gg"
	_ "github.com/dmac/gg/v2.1/gg"
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

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Tetris", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	const vshader = `#version 120

uniform mat4 proj, model;
attribute vec3 vertex_position;
attribute vec2 vertex_texture;
varying vec2 texture_coordinates;

void main() {
	gl_Position = proj * model * vec4(vertex_position, 1);
	texture_coordinates = vertex_texture;
}
`

	const fshader = `#version 120

uniform sampler2D tex_loc;
varying vec2 texture_coordinates;

void main() {
	gl_FragColor = texture2D(tex_loc, texture_coordinates);
}
`

	tetris, err := NewTetris(vshader, fshader)
	if err != nil {
		log.Fatal(err)
	}

	window.SetKeyCallback(tetris.handleKeyInput)

	for !window.ShouldClose() {
		tetris.Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func (t *Tetris) handleKeyInput(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press && action != glfw.Repeat {
		return
	}
	switch key {
	case glfw.KeyW:
		t.HandleInput(inputUp)
	case glfw.KeyD:
		t.HandleInput(inputRight)
	case glfw.KeyS:
		t.HandleInput(inputDown)
	case glfw.KeyA:
		t.HandleInput(inputLeft)
	case glfw.KeySpace:
		t.HandleInput(inputSpace)
	}
}

func LoadTextures() (map[string]*gg.Texture, error) {
	textures := make(map[string]*gg.Texture)
	for _, name := range []string{
		"bg", "board",
		"red", "orange", "yellow", "green", "blue", "cyan", "purple",
	} {
		t, err := newImageTexture(filepath.Join("images", name+".png"))
		if err != nil {
			return nil, err
		}
		textures[name] = t
	}
	return textures, nil
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
