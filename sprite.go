package gg

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
