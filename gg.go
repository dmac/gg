package gg

import (
	"fmt"

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

type Vec2 mgl.Vec2
type Vec3 mgl.Vec3
type Rect struct {
	Min, Max Vec2
}

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

func centroid2(vs []Vec2) Vec2 {
	var sx, sy float32
	for _, v := range vs {
		sx += v[0]
		sy += v[1]
	}
	l := float32(len(vs))
	return Vec2{sx / l, sy / l}
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
