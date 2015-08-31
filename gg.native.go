// +build !js

package gg

import (
	"fmt"
	"image"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

const vertexShader = `#version 330

uniform mat4 proj, view, model;
in vec3 vertex_position;
in vec2 vertex_texture;
out vec2 texture_coordinates;

void main() {
	gl_Position = proj * view * model * vec4(vertex_position, 1);
	texture_coordinates = vertex_texture;
}
`

const fragmentShader = `#version 330

uniform sampler2D tex_loc;
in vec2 texture_coordinates;
out vec4 color;

void main() {
	color = texture(tex_loc, texture_coordinates);
}
`

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

func NewTextureFromImage(img image.Image) *Texture {
	var buf []uint8
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			buf = append(buf, uint8(r/256))
			buf = append(buf, uint8(g/256))
			buf = append(buf, uint8(b/256))
			buf = append(buf, uint8(a/256))
		}
	}
	tex := &Texture{
		W: img.Bounds().Dx(),
		H: img.Bounds().Dy(),
	}
	gl.GenTextures(1, &tex.t)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.t)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(img.Bounds().Dx()),
		int32(img.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(buf),
	)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	return tex
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
