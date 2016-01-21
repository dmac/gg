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

type Game struct {
	triangle *Triangle
}

func NewGame(vertShader, fragShader string) (*Game, error) {
	gg.Enable(gg.DEPTH_TEST)
	gg.Enable(gg.CULL_FACE)
	gg.DepthFunc(gg.LESS)

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
		float32(WindowWidth) / 2, float32(WindowHeight)/2 - 50, 0,
		float32(WindowWidth)/2 - 50, float32(WindowHeight)/2 + 50, 0,
		float32(WindowWidth)/2 + 50, float32(WindowHeight)/2 + 50, 0,
	}

	// Initialize triangle
	triangle, err := NewTriangle(vertices, program)
	if err != nil {
		return nil, err
	}

	return &Game{triangle: triangle}, nil
}

func (g *Game) Draw() {
	gg.ClearColor(0.5, 0.5, 0.5, 1.0)
	gg.Clear(gg.COLOR_BUFFER_BIT | gg.DEPTH_BUFFER_BIT)
	g.triangle.Draw()
}

type Triangle struct {
	vbo     *gg.Buffer
	program *gg.Program
}

func NewTriangle(vertices []float32, program *gg.Program) (*Triangle, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, vertices); err != nil {
		return nil, err
	}
	vbo := gg.CreateBuffer()
	gg.BindBuffer(gg.ARRAY_BUFFER, vbo)
	gg.BufferData(gg.ARRAY_BUFFER, buf.Bytes(), gg.STATIC_DRAW)
	return &Triangle{
		vbo: vbo,
		program: program,
	}, nil
}

func (t *Triangle) Draw() {
	gg.UseProgram(t.program)
	colorUniform, err := gg.GetUniformLocation(t.program, "color")
	if err != nil {
		log.Fatal(err)
	}
	gg.Uniform4f(colorUniform, 1.0, 0.0, 1.0, 1.0)

	vattrib, err := gg.GetAttribLocation(t.program, "vertex_position")
	if err != nil {
		log.Fatal(err)
	}
	gg.EnableVertexAttribArray(vattrib)
	gg.BindBuffer(gg.ARRAY_BUFFER, t.vbo)
	gg.VertexAttribArrayPointer(vattrib, 3, gg.FLOAT, false, 0, 0)
	gg.DrawArrays(gg.TRIANGLE_FAN, 0, 3)
}
