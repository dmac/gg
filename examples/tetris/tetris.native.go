// +build !js

package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"

	"github.com/dmac/gg"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const WindowWidth = 640
const WindowHeight = 480

func main() {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	//glfw.WindowHint(glfw.ContextVersionMajor, 3)
	//glfw.WindowHint(glfw.ContextVersionMinor, 2)
	//glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Tetris", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	if err := gg.Init(WindowWidth, WindowHeight); err != nil {
		log.Fatal(err)
	}

	loadTextures()

	tetris := NewTetris(WindowWidth, WindowHeight)

	window.SetKeyCallback(tetris.handleKeyInput)

	for !window.ShouldClose() {
		gl.ClearColor(0.5, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

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
	}
}

var textures map[string]*gg.Texture

func loadTextures() {
	textures = make(map[string]*gg.Texture)
	for _, color := range []string{"blue", "cyan", "green", "orange", "purple", "red", "yellow"} {
		img, err := openPNG(color + ".png")
		if err != nil {
			fmt.Println(err)
			continue
		}
		textures[color] = gg.NewTextureFromImage(img)
	}
}

func openPNG(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return png.Decode(f)
}
