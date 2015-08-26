package gg

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

const vertexShader = `#version 330

in vec3 vertex_position;
out vec3 frag_color;

void main() {
        gl_Position = vec4(vertex_position, 1);
        frag_color = vertex_position;
}
`

const fragmentShader = `#version 330

in vec3 frag_color;
out vec4 color;

void main() {
        color = vec4(frag_color, 1.0);
}
`

type Vec3 mgl.Vec3

type Poly struct {
	polyPlatformData
	vertices []Vec3
}

func NewPoly(vertices []Vec3) *Poly {
	p := &Poly{
		vertices: vertices,
	}
	p.init()
	return p
}
