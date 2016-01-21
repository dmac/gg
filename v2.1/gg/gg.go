package gg_21

import (
	"fmt"
	"strings"

	"github.com/dmac/gg"
	"github.com/go-gl/gl/v2.1/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
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

func (*backend) VertexAttribArrayPointer(a *gg.Attribute, size int, typ gg.Enum, normalized bool, stride, offset int) {
	gl.VertexAttribPointer(
		a.Value.(uint32),
		int32(size),
		uint32(typ),
		normalized,
		int32(stride),
		gl.PtrOffset(offset),
	)
}

func (*backend) DrawArrays(mode gg.Enum, first, count int) {
	gl.DrawArrays(uint32(mode), int32(first), int32(count))
}

type Rect struct {
	Min, Max [2]float32
}

type Poly struct {
	transformable
	polyBackend
	color [4]float32
	n     int32
}

func NewPoly(vertices [][2]float32) *Poly {
	aabb := computeAABB(vertices)
	p := &Poly{
		transformable: transformable{
			Position: aabb.Min,
			scale:    1,
		},
		color: [4]float32{0, 0, 0, 1},
		n:     int32(len(vertices)),
	}
	p.init(vertices)
	return p
}

func (p *Poly) SetColor(r, g, b, a float32) {
	p.color[0] = r
	p.color[1] = g
	p.color[2] = b
	p.color[3] = a
}

func (p *Poly) transform() mgl.Mat4 {
	S := mgl.Scale2D(p.scale, p.scale).Mat4()
	R := mgl.Rotate2D(mgl.DegToRad(p.rotation)).Mat4()
	T := mgl.Translate3D(p.Position[0], p.Position[1], 0)
	return T.Mul4(R).Mul4(S)
}

type Sprite struct {
	transformable
	spriteBackend
	W   float32
	H   float32
	tex *Texture
	// field for area of texture to draw, for e.g., spritesheet
}

func NewSpriteFromTexture(tex *Texture) *Sprite {
	s := &Sprite{
		transformable: transformable{
			scale: 1,
		},
		W:   float32(tex.W),
		H:   float32(tex.H),
		tex: tex,
	}
	s.init()
	return s
}

func (s *Sprite) transform() mgl.Mat4 {
	S := mgl.Scale2D(s.scale, s.scale).Mat4()
	R := mgl.Rotate2D(mgl.DegToRad(s.rotation)).Mat4()
	T := mgl.Translate3D(s.Position[0], s.Position[1], 0)
	return T.Mul4(R).Mul4(S)
}

type Texture struct {
	textureBackend
	W int
	H int
}

type transformable struct {
	// TODO(dmac) origin Vec2
	Position [2]float32
	rotation float32
	scale    float32
}

func (t *transformable) SetPosition(x, y float32) {
	t.Position[0] = x
	t.Position[1] = y
}

func (t *transformable) SetRotation(degrees float32) {
	t.rotation = degrees
}

func (t *transformable) SetScale(s float32) {
	t.scale = s
}

func (t *transformable) Move(x, y float32) {
	t.Position[0] += x
	t.Position[1] += y
}

func (t *transformable) Rotate(degrees float32) {
	t.rotation += degrees
}

func (t *transformable) Scale(factor float32) {
	t.scale *= factor
}

func computeAABB(vs [][2]float32) Rect {
	if len(vs) < 3 {
		panic(fmt.Errorf("can't compute bounding box for %d vertices", len(vs)))
	}
	r := Rect{vs[0], vs[0]}
	for _, v := range vs {
		if v[0] < r.Min[0] {
			r.Min[0] = v[0]
		}
		if v[1] < r.Min[1] {
			r.Min[1] = v[1]
		}
		if v[0] > r.Max[0] {
			r.Max[0] = v[0]
		}
		if v[1] > r.Max[1] {
			r.Max[1] = v[1]
		}
	}
	return r
}
