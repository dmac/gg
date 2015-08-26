// +build !js

package main

import (
	"log"
	"runtime"
	"time"

	"github.com/dmac/gg"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func main() {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	window, err := glfw.CreateWindow(640, 480, "Hello", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	if err := gg.Init(); err != nil {
		log.Fatal(err)
	}

	triangle := gg.NewPoly([]gg.Vec3{
		{0.0, 0.5, 0.0},
		{0.5, -0.5, 0.0},
		{-0.5, -0.5, 0.0},
	})

	frames := 0
	last := time.Now()
	for !window.ShouldClose() {
		frames++
		if time.Since(last) > time.Second {
			//fmt.Println(frames)
			frames = 0
			last = time.Now()
		}

		gl.ClearColor(0.5, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		triangle.Draw()

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
