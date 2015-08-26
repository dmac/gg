package gg

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

const vertexShader = `#version 330

uniform mat4 proj, view, model;
in vec3 vertex_position;
out vec3 frag_color;

void main() {
	gl_Position = proj * view * model * vec4(vertex_position, 1);
}
`

const fragmentShader = `#version 330

out vec4 color;

void main() {
        color = vec4(1.0, 0.0, 1.0, 1.0);
}
`

type Vec2 mgl.Vec2
type Vec3 mgl.Vec3

type Poly struct {
	polyPlatformData
	center Vec2
	n int32
	scale float32
	rotation float32
}

func NewPoly(vertices []Vec2) *Poly {
	p := &Poly{
		center: centroid2(vertices),
		scale: 1,
		n: int32(len(vertices)),
	}
	p.init(vertices)
	return p
}

func (p *Poly) SetScale(s float32) {
	p.scale = s
}

func (p *Poly) Rotate(degrees float32) {
	p.rotation += degrees
}

func (p *Poly) Translate(x, y float32) {
	p.center[0] += x
	p.center[1] += y
}

func centroid2(vs []Vec2) Vec2 {
	var sx, sy float32
	for _, v := range vs {
		sx += v[0]
		sy += v[1]
	}
	l := float32(len(vs))
	return Vec2{sx/l, sy/l}
}
