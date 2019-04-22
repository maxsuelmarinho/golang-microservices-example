package service

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/disintegration/gift"
)

func Sepia(src image.Image, buf *bytes.Buffer) error {
	g := gift.New(
		gift.Resize(800, 0, gift.LanczosResampling),
		gift.Sepia(100),
	)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))

	g.Draw(dst, src)
	err := jpeg.Encode(buf, dst, nil)
	return err
}
