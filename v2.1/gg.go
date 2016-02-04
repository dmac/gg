package gg_21

import (
	"fmt"
	"strings"

	"github.com/dmac/gg"
	"github.com/go-gl/gl/v2.1/gl"
)

type backend struct{}

var _ gg.Backend = (*backend)(nil)

func init() {
	gg.Register(&backend{})
}

func (*backend) Enable(c gg.Enum) {
	gl.Enable(uint32(c))
}

func (*backend) DepthFunc(f gg.Enum) {
	gl.DepthFunc(uint32(f))
}

func (*backend) BlendFunc(sfactor, dfactor gg.Enum) {
	gl.BlendFunc(uint32(sfactor), uint32(dfactor))
}

func (*backend) Clear(mask gg.Enum) {
	gl.Clear(uint32(mask))
}

func (*backend) ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (*backend) CreateBuffer() *gg.Buffer {
	var b uint32
	gl.GenBuffers(1, &b)
	return &gg.Buffer{Value: b}
}

func (*backend) BindBuffer(typ gg.Enum, b *gg.Buffer) {
	gl.BindBuffer(uint32(typ), b.Value.(uint32))
}

func (*backend) BufferData(typ gg.Enum, src []byte, usage gg.Enum) {
	gl.BufferData(uint32(typ), len(src), gl.Ptr(src), uint32(usage))
}

func (*backend) CreateShader(src []byte, typ gg.Enum) (*gg.Shader, error) {
	csrc := gl.Str(string(append([]byte(src), 0)))
	shader := gl.CreateShader(uint32(typ))
	gl.ShaderSource(shader, 1, &csrc, nil)
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.TRUE {
		return &gg.Shader{Value: shader}, nil
	}
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
	return nil, fmt.Errorf("compile shader: %s%s", src, log)
}

func (*backend) CreateProgram() *gg.Program {
	p := gl.CreateProgram()
	return &gg.Program{Value: p}
}

func (*backend) AttachShader(p *gg.Program, s *gg.Shader) {
	gl.AttachShader(p.Value.(uint32), s.Value.(uint32))
}

func (*backend) LinkProgram(p *gg.Program) error {
	pv := p.Value.(uint32)
	gl.LinkProgram(pv)
	var status int32
	gl.GetProgramiv(pv, gl.LINK_STATUS, &status)
	if status == gl.TRUE {
		return nil
	}
	var logLength int32
	gl.GetProgramiv(pv, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(pv, logLength, nil, gl.Str(log))
	return fmt.Errorf("link program: %s", log)
}

func (*backend) UseProgram(p *gg.Program) {
	gl.UseProgram(p.Value.(uint32))
}

func (*backend) GetUniformLocation(p *gg.Program, name string) (*gg.Uniform, error) {
	u := gl.GetUniformLocation(p.Value.(uint32), gl.Str(name+"\x00"))
	if u < 0 {
		return nil, fmt.Errorf("gg: no uniform named " + name)
	}
	return &gg.Uniform{Value: u}, nil
}

func (*backend) Uniform1f(u *gg.Uniform, v0 float32) {
	gl.Uniform1f(u.Value.(int32), v0)
}

func (*backend) Uniform1i(u *gg.Uniform, v0 int) {
	gl.Uniform1i(u.Value.(int32), int32(v0))
}

func (*backend) Uniform4f(u *gg.Uniform, v0, v1, v2, v3 float32) {
	gl.Uniform4f(u.Value.(int32), v0, v1, v2, v3)
}

func (*backend) UniformMatrix4fv(u *gg.Uniform, values []float32) {
	gl.UniformMatrix4fv(u.Value.(int32), 1, false, &values[0])
}

func (*backend) GetAttribLocation(p *gg.Program, name string) (*gg.Attribute, error) {
	a := gl.GetAttribLocation(p.Value.(uint32), gl.Str(name+"\x00"))
	if a < 0 {
		return nil, fmt.Errorf("gg: no attribute named " + name)
	}
	return &gg.Attribute{Value: uint32(a)}, nil
}

func (*backend) EnableVertexAttribArray(a *gg.Attribute) {
	gl.EnableVertexAttribArray(a.Value.(uint32))
}

func (*backend) VertexAttribPointer(a *gg.Attribute, size int, typ gg.Enum, normalized bool, stride, offset int) {
	gl.VertexAttribPointer(
		a.Value.(uint32),
		int32(size),
		uint32(typ),
		normalized,
		int32(stride),
		gl.PtrOffset(offset),
	)
}

func (*backend) CreateTexture() *gg.Texture {
	var t uint32
	gl.GenTextures(1, &t)
	return &gg.Texture{Value: t}
}

func (*backend) ActiveTexture(tex gg.Enum) {
	gl.ActiveTexture(uint32(tex))
}

func (*backend) BindTexture(target gg.Enum, texture *gg.Texture) {
	gl.BindTexture(uint32(target), texture.Value.(uint32))
}

func (*backend) TexImage2D(
	target gg.Enum, level int, internalFormat gg.Enum,
	width, height, border int,
	format, typ gg.Enum,
	data interface{},
) {
	gl.TexImage2D(
		gl.TEXTURE_2D, int32(level), gl.RGBA,
		int32(width), int32(height), int32(border),
		uint32(format), uint32(typ),
		gl.Ptr(data),
	)
}

func (*backend) TexParameteri(target gg.Enum, pname gg.Enum, param gg.Enum) {
	gl.TexParameteri(uint32(target), uint32(pname), int32(param))
}

func (*backend) DrawArrays(mode gg.Enum, first, count int) {
	gl.DrawArrays(uint32(mode), int32(first), int32(count))
}
