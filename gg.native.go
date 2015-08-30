// +build !js

package gg

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

var program uint32

func Init() error {
	var err error
	program, err = linkProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}
	gl.UseProgram(program)

	dims := make([]float32, 4)
	gl.GetFloatv(gl.VIEWPORT, &dims[0])

	proj := mgl.Ortho(0, 640, 480, 0, 0, 1)
	projUniform := gl.GetUniformLocation(program, gl.Str("proj\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &proj[0])

	view := mgl.Ident4()
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	return nil
}

type polyPlatformData struct {
	vao     uint32
	program uint32
}

func (p *Poly) init(vertices []Vec2) {
	p.program = program

	var mesh []Vec3
	for _, vertex := range vertices {
		v := Vec3{
			vertex[0] - p.position[0],
			vertex[1] - p.position[1],
			0,
		}
		mesh = append(mesh, v)
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(mesh)*int(unsafe.Sizeof(Vec3{})),
		gl.Ptr(mesh),
		gl.STATIC_DRAW,
	)

	gl.GenVertexArrays(1, &p.vao)
	gl.BindVertexArray(p.vao)

	vattrib := uint32(gl.GetAttribLocation(p.program, gl.Str("vertex_position\x00")))
	gl.EnableVertexAttribArray(vattrib)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(vattrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
}

func (p *Poly) Draw() {
	gl.UseProgram(p.program)

	model := p.transform()
	modelUniform := gl.GetUniformLocation(p.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindVertexArray(p.vao)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, p.n)
}

func (p *Poly) transform() mgl.Mat4 {
	S := mgl.Scale2D(p.scale, p.scale).Mat4()
	R := mgl.Rotate2D(mgl.DegToRad(p.rotation)).Mat4()
	T := mgl.Translate3D(p.position[0], p.position[1], 0)
	return T.Mul4(R).Mul4(S)
}

type spritePlatformData struct {
	vao     uint32
	program uint32
}

func (s *Sprite) init() {
	s.program = program

	vmesh := []Vec3{
		{0, 0},
		{0, s.H},
		{s.W, s.H},
		{s.W, 0},
	}
	var vvbo uint32
	gl.GenBuffers(1, &vvbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vvbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(vmesh)*int(unsafe.Sizeof(Vec3{})),
		gl.Ptr(vmesh),
		gl.STATIC_DRAW,
	)

	tmesh := []Vec2{
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 0},
	}
	var tvbo uint32
	gl.GenBuffers(1, &tvbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, tvbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(tmesh)*int(unsafe.Sizeof(Vec2{})),
		gl.Ptr(tmesh),
		gl.STATIC_DRAW,
	)

	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	vattrib := uint32(gl.GetAttribLocation(s.program, gl.Str("vertex_position\x00")))
	gl.EnableVertexAttribArray(vattrib)
	gl.BindBuffer(gl.ARRAY_BUFFER, vvbo)
	gl.VertexAttribPointer(vattrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	tattrib := uint32(gl.GetAttribLocation(s.program, gl.Str("vertex_texture\x00")))
	gl.EnableVertexAttribArray(tattrib)
	gl.BindBuffer(gl.ARRAY_BUFFER, tvbo)
	gl.VertexAttribPointer(tattrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
}

func (s *Sprite) Draw() {
	gl.UseProgram(s.program)

	model := s.transform()
	modelUniform := gl.GetUniformLocation(s.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.tex.t)
	textureUniform := gl.GetUniformLocation(s.program, gl.Str("tex_loc\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}

func (s *Sprite) transform() mgl.Mat4 {
	S := mgl.Scale2D(s.scale, s.scale).Mat4()
	R := mgl.Rotate2D(mgl.DegToRad(s.rotation)).Mat4()
	T := mgl.Translate3D(s.position[0], s.position[1], 0)
	return T.Mul4(R).Mul4(S)
}

func linkProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vshader, err := compileShader(vertexShader, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fshader, err := compileShader(fragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	program := gl.CreateProgram()
	gl.AttachShader(program, vshader)
	gl.AttachShader(program, fshader)
	gl.LinkProgram(program)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.TRUE {
		return program, nil
	}
	var logLength int32
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
	return 0, fmt.Errorf("link program: %s", log)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	csource := gl.Str(string(append([]byte(source), 0)))
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.TRUE {
		return shader, nil
	}
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
	return 0, fmt.Errorf("compile shader: %s%s", source, log)
}
