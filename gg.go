package gg

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Vec2 mgl.Vec2
type Vec3 mgl.Vec3
type Vec4 mgl.Vec4
type Rect struct {
	Min, Max Vec2
}

type Poly struct {
	transformable
	polyPlatformData
	color [4]float32
	n     int32
	// TODO(dmac) color
}

func NewPoly(vertices []Vec2) *Poly {
	aabb := computeAABB(vertices)
	p := &Poly{
		transformable: transformable{
			position: aabb.Min,
			scale:    1,
		},
		color: [4]float32{0, 0, 0, 1},
		n: int32(len(vertices)),
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
	T := mgl.Translate3D(p.position[0], p.position[1], 0)
	return T.Mul4(R).Mul4(S)
}

type Sprite struct {
	transformable
	spritePlatformData
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
	T := mgl.Translate3D(s.position[0], s.position[1], 0)
	return T.Mul4(R).Mul4(S)
}

type Texture struct {
	texturePlatformData
	W int
	H int
}

type transformable struct {
	// TODO(dmac) origin Vec2
	position Vec2
	rotation float32
	scale    float32
}

func (t *transformable) SetPosition(x, y float32) {
	t.position[0] = x
	t.position[1] = y
}

func (t *transformable) SetRotation(degrees float32) {
	t.rotation = degrees
}

func (t *transformable) SetScale(s float32) {
	t.scale = s
}

func (t *transformable) Position(x, y float32) {
	t.position[0] += x
	t.position[1] += y
}

func (t *transformable) Rotate(degrees float32) {
	t.rotation += degrees
}

func (t *transformable) Scale(factor float32) {
	t.scale *= factor
}

func computeAABB(vs []Vec2) Rect {
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
