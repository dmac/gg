package gg

import mgl "github.com/go-gl/mathgl/mgl32"

type Poly struct {
	transformable
	polyPlatformData
	n int32
	// TODO(dmac) color
}

func NewPoly(vertices []Vec2) *Poly {
	aabb := computeAABB(vertices)
	p := &Poly{
		transformable: transformable{
			position: aabb.Min,
			scale:    1,
		},
		n: int32(len(vertices)),
	}
	p.init(vertices)
	return p
}

func (p *Poly) transform() mgl.Mat4 {
	S := mgl.Scale2D(p.scale, p.scale).Mat4()
	R := mgl.Rotate2D(mgl.DegToRad(p.rotation)).Mat4()
	T := mgl.Translate3D(p.position[0], p.position[1], 0)
	return T.Mul4(R).Mul4(S)
}
