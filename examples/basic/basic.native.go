// +build !js

package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/dmac/gg/v2.1/gg"
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
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Basic", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	if err := gg.Init(WindowWidth, WindowHeight); err != nil {
		log.Fatal(err)
	}

	triangle := gg.NewPoly([][2]float32{
		{320, 100},
		{240, 200},
		{400, 200},
	})
	triangle.SetColor(1, 0, 1, 1)

	img1, err := NewImageFromFile("test.png")
	if err != nil {
		log.Fatal(err)
	}
	img2, err := NewImageFromFile("test2.png")
	if err != nil {
		log.Fatal(err)
	}

	tex1 := gg.NewTextureFromImage(img1)
	tex2 := gg.NewTextureFromImage(img2)
	spr1 := gg.NewSpriteFromTexture(tex1)
	spr2 := gg.NewSpriteFromTexture(tex2)
	spr1.SetPosition(240, 230)
	spr2.SetPosition(340, 230)

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
		spr1.Draw()
		spr2.Draw()

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func NewImageFromFile(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return png.Decode(f)
}
