package FrostedGlassLike

import (
	"image"
	"image/color"
	"os"
)

// just load your file and it will reply (
func FrostedGlassLike(src os.File, radius int, err error) (dst image.Image) {
	srcImg, _, err := image.Decode(&src)
	if err != nil {
		panic(err)
	}
	bounds := srcImg.Bounds()
	dstImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := 0, 0, 0, 0
			count := 0
			for i := -radius; i <= radius; i++ {
				for j := -radius; j <= radius; j++ {
					dx, dy := x+i, y+j
					if dx < bounds.Min.X || dx >= bounds.Max.X || dy < bounds.Min.Y || dy >= bounds.Max.Y {
						continue
					}
					c := color.RGBAModel.Convert(srcImg.At(dx, dy)).(color.RGBA)
					r += int(c.R)
					g += int(c.G)
					b += int(c.B)
					a += int(c.A)
					count++
				}
			}
			dstImg.Set(x, y, color.RGBA{
				R: uint8(r / count),
				G: uint8(g / count),
				B: uint8(b / count),
				A: uint8(a / count),
			})
		}
	}
	return dstImg
}
