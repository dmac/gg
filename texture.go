package gg

import (
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture struct {
	W int
	H int
	t uint32
}

func NewTextureFromImage(img image.Image) *Texture {
	var buf []uint8
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			buf = append(buf, uint8(r/256))
			buf = append(buf, uint8(g/256))
			buf = append(buf, uint8(b/256))
			buf = append(buf, uint8(a/256))
		}
	}
	tex := &Texture{
		W: img.Bounds().Dx(),
		H: img.Bounds().Dy(),
	}
	gl.GenTextures(1, &tex.t)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, tex.t)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(img.Bounds().Dx()),
		int32(img.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(buf),
	)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	return tex
}
