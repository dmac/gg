package main

import (
	"log"
	"runtime"
	"testing"

	"github.com/dmac/gg/v2.1/gg"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const WindowWidth = 640
const WindowHeight = 480

func TestInit(t *testing.T) {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Test", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	if err := gg.Init(WindowWidth, WindowHeight); err != nil {
		log.Fatal(err)
	}
}
