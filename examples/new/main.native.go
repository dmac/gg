package main

import (
	"log"
	"runtime"

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

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Basic", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}

	for !window.ShouldClose() {
		game.Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
