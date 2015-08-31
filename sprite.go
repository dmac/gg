package gg

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	transformable
	spritePlatformData
	W       float32
	H       float32
	tex     *Texture
	program uint32
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
