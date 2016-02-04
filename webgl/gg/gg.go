package gg_webgl

import (
	"fmt"

	"github.com/dmac/gg"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

type backend struct {
	gl *webgl.Context
}

var _ gg.Backend = (*backend)(nil)

func Init(gl *webgl.Context) {
	gg.Register(&backend{gl: gl})
}

func (b *backend) Enable(c gg.Enum) {
	b.gl.Enable(int(c))
}

func (b *backend) DepthFunc(f gg.Enum) {
	b.gl.DepthFunc(int(f))
}

func (b *backend) BlendFunc(sfactor, dfactor gg.Enum) {
	b.gl.BlendFunc(int(sfactor), int(dfactor))
}

func (b *backend) Clear(mask gg.Enum) {
	b.gl.Clear(int(mask))
}

func (be *backend) ClearColor(r, g, b, a float32) {
	be.gl.ClearColor(r, g, b, a)
}

func (b *backend) CreateBuffer() *gg.Buffer {
	return &gg.Buffer{Value: b.gl.CreateBuffer()}
}

func (b *backend) BindBuffer(typ gg.Enum, buf *gg.Buffer) {
	b.gl.BindBuffer(int(typ), buf.Value.(*js.Object))
}

func (b *backend) BufferData(typ gg.Enum, src []byte, usage gg.Enum) {
	b.gl.BufferData(int(typ), src, int(usage))
}

func (b *backend) CreateShader(src []byte, typ gg.Enum) (*gg.Shader, error) {
	shader := b.gl.CreateShader(int(typ))
	b.gl.ShaderSource(shader, string(src))
	b.gl.CompileShader(shader)
	log := b.gl.GetShaderInfoLog(shader)
	if log != "" {
		return nil, fmt.Errorf("gg: compile shader: %s%s", src, log)
	}
	return &gg.Shader{Value: shader}, nil
}

func (b *backend) CreateProgram() *gg.Program {
	return &gg.Program{Value: b.gl.CreateProgram()}
}

func (b *backend) AttachShader(p *gg.Program, s *gg.Shader) {
	b.gl.AttachShader(p.Value.(*js.Object), s.Value.(*js.Object))
}

func (b *backend) LinkProgram(p *gg.Program) error {
	b.gl.LinkProgram(p.Value.(*js.Object))
	log := b.gl.GetProgramInfoLog(p.Value.(*js.Object))
	if log != "" {
		return fmt.Errorf("gg: link program: %s", log)
	}
	return nil
}

func (b *backend) UseProgram(p *gg.Program) {
	b.gl.UseProgram(p.Value.(*js.Object))
}

func (b *backend) GetUniformLocation(p *gg.Program, name string) (*gg.Uniform, error) {
	u := b.gl.GetUniformLocation(p.Value.(*js.Object), name)
	if u.Int() < 0 {
		return nil, fmt.Errorf("gg: no uniform named " + name)
	}
	return &gg.Uniform{Value: u}, nil
}

func (b *backend) Uniform1f(u *gg.Uniform, v0 float32) {
	b.gl.Uniform1f(u.Value.(*js.Object), v0)
}

func (b *backend) Uniform1i(u *gg.Uniform, v0 int) {
	b.gl.Uniform1i(u.Value.(*js.Object), v0)
}

func (b *backend) Uniform4f(u *gg.Uniform, v0, v1, v2, v3 float32) {
	b.gl.Uniform4f(u.Value.(*js.Object), v0, v1, v2, v3)
}

func (b *backend) UniformMatrix4fv(u *gg.Uniform, values []float32) {
	b.gl.UniformMatrix4fv(u.Value.(*js.Object), false, values)
}

func (b *backend) GetAttribLocation(p *gg.Program, name string) (*gg.Attribute, error) {
	a := b.gl.GetAttribLocation(p.Value.(*js.Object), name)
	if a < 0 {
		return nil, fmt.Errorf("gg: no attribute named " + name)
	}
	return &gg.Attribute{Value: a}, nil
}

func (b *backend) EnableVertexAttribArray(a *gg.Attribute) {
	b.gl.EnableVertexAttribArray(a.Value.(int))
}

func (b *backend) VertexAttribPointer(a *gg.Attribute, size int, typ gg.Enum, normalized bool, stride, offset int) {
	b.gl.VertexAttribPointer(a.Value.(int), size, int(typ), normalized, stride, offset)
}

func (b *backend) CreateTexture() *gg.Texture {
	t := b.gl.CreateTexture()
	return &gg.Texture{Value: t}
}

func (b *backend) ActiveTexture(tex gg.Enum) {
	b.gl.ActiveTexture(int(tex))
}

func (b *backend) BindTexture(target gg.Enum, texture *gg.Texture) {
	b.gl.BindTexture(int(target), texture.Value.(*js.Object))
}

func (b *backend) TexImage2D(
	target gg.Enum, level int, internalFormat gg.Enum,
	width, height, border int,
	format, typ gg.Enum,
	data interface{},
) {
	b.gl.TexImage2D(
		gg.TEXTURE_2D, int(level), gg.RGBA,
		gg.RGBA, gg.UNSIGNED_BYTE,
		data.(*js.Object),
	)
}

func (b *backend) TexParameteri(target gg.Enum, pname gg.Enum, param gg.Enum) {
	b.gl.TexParameteri(int(target), int(pname), int(param))
}

func (b *backend) DrawArrays(mode gg.Enum, first, count int) {
	b.gl.DrawArrays(int(mode), first, count)
}
