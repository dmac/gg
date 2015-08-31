// +build js

package gg

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

var gl *webgl.Context
var program *js.Object

const vertexShader = `
uniform mat4 proj, view, model;
// in
attribute vec3 vertex_position;
//varying vec2 vertex_texture;
// out
//varying vec2 texture_coordinates;

void main() {
	gl_Position = proj * view * model * vec4(vertex_position, 1);
	//texture_coordinates = vertex_texture;
}
`

const fragmentShader = `
//uniform sampler2D tex_loc;
// in
//varying vec2 texture_coordinates;

void main() {
	//gl_FragColor = texture(tex_loc, texture_coordinates);
	gl_FragColor = vec4(1.0, 0.0, 1.0, 1.0);
}
`

type polyPlatformData struct {
	program *js.Object
	vvbo    *js.Object
}

type spritePlatformData struct {
	program *js.Object
}

func Init(glContext *webgl.Context) error {
	gl = glContext

	var err error
	program, err = linkProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}
	gl.UseProgram(program)

	proj := mgl.Ortho(0, 640, 480, 0, 0, 1)
	projUniform := gl.GetUniformLocation(program, "proj")
	gl.UniformMatrix4fv(projUniform, false, proj[:])

	view := mgl.Ident4()
	viewUniform := gl.GetUniformLocation(program, "view")
	gl.UniformMatrix4fv(viewUniform, false, view[:])

	return nil
}

func (p *Poly) init(vertices []Vec2) {
	p.program = program

	var mesh []float32
	for _, v := range vertices {
		mesh = append(mesh, v[0])
		mesh = append(mesh, v[1])
		mesh = append(mesh, 0)
	}

	p.vvbo = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, p.vvbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		mesh,
		gl.STATIC_DRAW,
	)

	vattrib := gl.GetAttribLocation(program, "vertex_position")
	gl.EnableVertexAttribArray(vattrib)
	gl.VertexAttribPointer(vattrib, 3, gl.FLOAT, false, 0, 0)
}

func (p *Poly) Draw() {
	gl.UseProgram(p.program)

	model := p.transform()
	modelUniform := gl.GetUniformLocation(p.program, "model")
	gl.UniformMatrix4fv(modelUniform, false, model[:])

	gl.BindBuffer(gl.ARRAY_BUFFER, p.vvbo)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, int(p.n))
}

func (s *Sprite) init() {
}

func linkProgram(vertexShaderSource, fragmentShaderSource string) (*js.Object, error) {
	vshader, err := compileShader(vertexShader, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fshader, err := compileShader(fragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	program := gl.CreateProgram()
	gl.AttachShader(program, vshader)
	gl.AttachShader(program, fshader)
	gl.LinkProgram(program)
	log := gl.GetProgramInfoLog(program)
	if log == "" {
		return program, nil
	}
	return nil, fmt.Errorf("link program: %s", log)
}

func compileShader(source string, shaderType int) (*js.Object, error) {
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, source)
	gl.CompileShader(shader)
	log := gl.GetShaderInfoLog(shader)
	if log == "" {
		return shader, nil
	}
	return nil, fmt.Errorf("compile shader: %s%s", source, log)
}
