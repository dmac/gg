package main

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/dmac/gg"
	mgl "github.com/go-gl/mathgl/mgl32"
)

const WindowWidth = 640
const WindowHeight = 480

type Scene struct {
	sprite *Sprite
}

func NewScene(vertShader, fragShader string, texture *gg.Texture) (*Scene, error) {
	gg.Enable(gg.DEPTH_TEST)
	gg.Enable(gg.CULL_FACE)
	gg.DepthFunc(gg.LESS)

	gg.Enable(gg.BLEND)
	gg.BlendFunc(gg.SRC_ALPHA, gg.ONE_MINUS_SRC_ALPHA)

	// Compile the global shader program
	vshader, err := gg.CreateShader([]byte(vertShader), gg.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fshader, err := gg.CreateShader([]byte(fragShader), gg.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	program := gg.CreateProgram()
	gg.AttachShader(program, vshader)
	gg.AttachShader(program, fshader)
	if err := gg.LinkProgram(program); err != nil {
		return nil, err
	}

	// Set the global projection matrix
	gg.UseProgram(program)
	proj := mgl.Ortho(0, float32(WindowWidth), float32(WindowHeight), 0, 0, 1)
	projUniform, err := gg.GetUniformLocation(program, "proj")
	if err != nil {
		return nil, err
	}
	gg.UniformMatrix4fv(projUniform, proj[:])

	vertices := []float32{
		float32(WindowWidth)/2 - 50, float32(WindowHeight)/2 - 50, 0,
		float32(WindowWidth)/2 - 50, float32(WindowHeight)/2 + 50, 0,
		float32(WindowWidth)/2 + 50, float32(WindowHeight)/2 + 50, 0,
		float32(WindowWidth)/2 + 50, float32(WindowHeight)/2 - 50, 0,
	}

	// Initialize sprite
	sprite, err := NewSprite(vertices, program, texture)
	if err != nil {
		return nil, err
	}

	return &Scene{sprite: sprite}, nil
}

func (s *Scene) Draw() {
	gg.ClearColor(0.5, 0.5, 0.5, 1.0)
	gg.Clear(gg.COLOR_BUFFER_BIT | gg.DEPTH_BUFFER_BIT)
	s.sprite.Draw()
}

type Sprite struct {
	pvbo    *gg.Buffer
	tvbo    *gg.Buffer
	program *gg.Program
	tex     *gg.Texture
}

func NewSprite(vertices []float32, program *gg.Program, texture *gg.Texture) (*Sprite, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, vertices); err != nil {
		return nil, err
	}
	pvbo := gg.CreateBuffer()
	gg.BindBuffer(gg.ARRAY_BUFFER, pvbo)
	gg.BufferData(gg.ARRAY_BUFFER, buf.Bytes(), gg.STATIC_DRAW)

	texVertices := []float32{
		0, 0,
		0, 1,
		1, 1,
		1, 0,
	}
	buf.Reset()
	if err := binary.Write(buf, binary.LittleEndian, texVertices); err != nil {
		return nil, err
	}
	tvbo := gg.CreateBuffer()
	gg.BindBuffer(gg.ARRAY_BUFFER, tvbo)
	gg.BufferData(gg.ARRAY_BUFFER, buf.Bytes(), gg.STATIC_DRAW)

	return &Sprite{
		pvbo:    pvbo,
		tvbo:    tvbo,
		program: program,
		tex:     texture,
	}, nil
}

func (s *Sprite) Draw() {
	gg.UseProgram(s.program)

	vattrib, err := gg.GetAttribLocation(s.program, "vertex_position")
	if err != nil {
		log.Fatal(err)
	}
	gg.EnableVertexAttribArray(vattrib)
	gg.BindBuffer(gg.ARRAY_BUFFER, s.pvbo)
	gg.VertexAttribPointer(vattrib, 3, gg.FLOAT, false, 0, 0)

	tattrib, err := gg.GetAttribLocation(s.program, "vertex_texture")
	if err != nil {
		log.Fatal(err)
	}
	gg.EnableVertexAttribArray(tattrib)
	gg.BindBuffer(gg.ARRAY_BUFFER, s.tvbo)
	gg.VertexAttribPointer(tattrib, 2, gg.FLOAT, false, 0, 0)

	gg.ActiveTexture(gg.TEXTURE0)
	gg.BindTexture(gg.TEXTURE_2D, s.tex)
	texUniform, err := gg.GetUniformLocation(s.program, "tex_loc")
	if err != nil {
		log.Fatal(err)
	}
	gg.Uniform1i(texUniform, 0)

	gg.DrawArrays(gg.TRIANGLE_FAN, 0, 4)
}
