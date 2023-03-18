package colorPicker

import (
	"image"
	"image/color"
)

// GetMainColor Idea from stack Overflow
// (summarize main color and be in a map,when needed,decode an img file and compare which is main color)
func GetMainColor(img image.Image) color.Color {
	model := color.Palette{
		color.RGBA{A: 255},                         // 黑色
		color.RGBA{R: 255, A: 255},                 // 红色
		color.RGBA{G: 255, A: 255},                 // 绿色
		color.RGBA{B: 255, A: 255},                 // 蓝色
		color.RGBA{R: 255, G: 255, A: 255},         // 黄色
		color.RGBA{G: 255, B: 255, A: 255},         // 青色
		color.RGBA{R: 255, B: 255, A: 255},         // 品红色
		color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色
	}
	counts := make(map[color.Color]int)
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			c := model.Convert(img.At(x, y))
			counts[c]++
		}
	}
	var maxColor color.Color
	var maxCount int
	for c, count := range counts {
		if count > maxCount {
			maxColor = c
			maxCount = count
		}
	}
	return maxColor
}
