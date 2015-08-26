// +build js

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

func main() {
	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	document.Get("body").Call("appendChild", canvas)

	attrs := webgl.DefaultAttributes()
	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		log.Fatal(err)
	}

	lastSec := time.Now()
	last := time.Now()
	fps := 60.0
	targetMillisPerFrame := 1000/fps
	fmt.Println(targetMillisPerFrame)
	frames := 0
	for {
		frames++
		if time.Since(lastSec) > time.Second {
			fmt.Println(float64(frames) / float64(time.Since(lastSec).Nanoseconds()) * 1e9)
			frames = 0
			lastSec = time.Now()
		}

		gl.ClearColor(0.5, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		since := time.Since(last)
		time.Sleep(time.Duration(targetMillisPerFrame * 1e6) * time.Nanosecond - since)
		last = time.Now()
	}
}
