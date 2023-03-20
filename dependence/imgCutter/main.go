package imgCutter

import (
	"golang.org/x/image/draw"
	"image"
)

// Python code with CatmullRom but written in golang.
func cropImage(img image.Image, width, height int) (destimg image.Image) {
	sourceWidth := img.Bounds().Dx()
	sourceHeight := img.Bounds().Dy()
	sourceAspect := float64(sourceWidth) / float64(sourceHeight)
	destAspect := float64(width) / float64(height)
	var rect image.Rectangle
	if sourceAspect == destAspect {
		rect = image.Rect(0, 0, width, height)
	} else if sourceAspect > destAspect {
		destHeight := int(float64(sourceHeight) * destAspect / sourceAspect)
		destY := (height - destHeight) / 2
		rect = image.Rect(0, destY, width, destY+destHeight)
	} else {
		destWidth := int(float64(sourceWidth) * sourceAspect / destAspect)
		destX := (width - destWidth) / 2
		rect = image.Rect(destX, 0, destX+destWidth, height)
	}
	destImg := image.NewRGBA(rect)
	draw.Draw(destImg, destImg.Bounds(), img, rect.Min, draw.Src)
	if sourceWidth < width || sourceHeight < height {
		scale := 1.0
		if destAspect > sourceAspect {
			scale = float64(width) / float64(sourceWidth)
		} else {
			scale = float64(height) / float64(sourceHeight)
		}
		newWidth := int(float64(sourceWidth) * scale)
		newHeight := int(float64(sourceHeight) * scale)
		newRect := image.Rect(0, 0, newWidth, newHeight)
		newImg := image.NewRGBA(newRect)
		draw.CatmullRom.Scale(newImg, newRect, destImg, destImg.Bounds(), draw.Over, nil)
		destImg = newImg
	}
	return destImg
}
