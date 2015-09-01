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
attribute vec2 vertex_texture;
// out
varying highp vec2 texture_coordinates;

void main() {
	gl_Position = proj * view * model * vec4(vertex_position, 1);
	texture_coordinates = vertex_texture;
}
`

const fragmentShader = `
uniform sampler2D tex_loc;
uniform highp float mix_value;
uniform highp vec4 color;
// in
varying highp vec2 texture_coordinates;

void main() {
	gl_FragColor = mix(
		color,
		texture2D(tex_loc, texture_coordinates),
		mix_value
	);
}
`

type polyPlatformData struct {
	program *js.Object
	vvbo    *js.Object
}

type spritePlatformData struct {
	program *js.Object
	vvbo    *js.Object
	tvbo    *js.Object
}

type texturePlatformData struct {
	t *js.Object
}

func Init(canvasWidth, canvasHeight int, glContext *webgl.Context) error {
	gl = glContext

	var err error
	program, err = linkProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}
	gl.UseProgram(program)

	proj := mgl.Ortho(0, float32(canvasWidth), float32(canvasHeight), 0, 0, 1)
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
		mesh = append(mesh, v[0] - p.position[0])
		mesh = append(mesh, v[1] - p.position[1])
		mesh = append(mesh, 0)
	}

	p.vvbo = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, p.vvbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		mesh,
		gl.STATIC_DRAW,
	)
}

func (p *Poly) Draw() {
	gl.UseProgram(p.program)

	model := p.transform()
	modelUniform := gl.GetUniformLocation(p.program, "model")
	gl.UniformMatrix4fv(modelUniform, false, model[:])

	gl.BindBuffer(gl.ARRAY_BUFFER, p.vvbo)
	vattrib := gl.GetAttribLocation(program, "vertex_position")
	gl.EnableVertexAttribArray(vattrib)
	gl.VertexAttribPointer(vattrib, 3, gl.FLOAT, false, 0, 0)

	mixUniform := gl.GetUniformLocation(p.program, "mix_value")
	gl.Uniform1f(mixUniform, 0.0)

	colorUniform := gl.GetUniformLocation(p.program, "color")
	gl.Uniform4f(colorUniform, p.color[0], p.color[1], p.color[2], p.color[3])

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, int(p.n))
}

func (s *Sprite) init() {
	s.program = program

	vmesh := []float32{
		0, 0, 0,
		0, s.W, 0,
		s.W, s.H, 0,
		s.W, 0, 0,
	}
	s.vvbo = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vvbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		vmesh,
		gl.STATIC_DRAW,
	)

	tmesh := []float32{
		0, 0,
		0, 1,
		1, 1,
		1, 0,
	}
	s.tvbo = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, s.tvbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		tmesh,
		gl.STATIC_DRAW,
	)
}

func (s *Sprite) Draw() {
	gl.UseProgram(s.program)

	model := s.transform()
	modelUniform := gl.GetUniformLocation(s.program, "model")
	gl.UniformMatrix4fv(modelUniform, false, model[:])

	gl.BindBuffer(gl.ARRAY_BUFFER, s.vvbo)
	vattrib := gl.GetAttribLocation(program, "vertex_position")
	gl.EnableVertexAttribArray(vattrib)
	gl.VertexAttribPointer(vattrib, 3, gl.FLOAT, false, 0, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.tex.t)
	textureUniform := gl.GetUniformLocation(program, "tex_loc")
	gl.Uniform1i(textureUniform, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.tvbo)
	tattrib := gl.GetAttribLocation(program, "vertex_texture")
	gl.EnableVertexAttribArray(tattrib)
	gl.VertexAttribPointer(tattrib, 2, gl.FLOAT, false, 0, 0)

	mixUniform := gl.GetUniformLocation(s.program, "mix_value")
	gl.Uniform1f(mixUniform, 1.0)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}

func NewTextureFromImage(img *js.Object) *Texture {
	tex := &Texture{
		W: img.Get("width").Int(),
		H: img.Get("height").Int(),
		texturePlatformData: texturePlatformData{
			t: gl.CreateTexture(),
		},
	}
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex.t)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		img,
	)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	return tex
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
