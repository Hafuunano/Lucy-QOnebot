package pgr

import (
	"encoding/json"
	"github.com/FloatTech/gg"
	"image"
	"image/color"
	"math"
	"strconv"
)

var (
	phigrosB19       PhigrosStruct
	font                     = engine.DataFolder() + "rec/font.ttf"
	background               = engine.DataFolder() + "rec/illustrationLowRes/" // .png
	backgroundRender         = engine.DataFolder() + "rec/background.png"
	rank                     = engine.DataFolder() + "rec/rank/"
	ChanllengeMode           = engine.DataFolder() + "rec/challengemode/"
	icon                     = engine.DataFolder() + "rec/icon.png"
	a                float64 = 75
)

type PhigrosStruct struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Content struct {
		Phi      bool `json:"phi"`
		BestList []struct {
			Score      int     `json:"score"`
			Acc        float64 `json:"acc"`
			Level      string  `json:"level"`
			Fc         bool    `json:"fc"`
			SongId     string  `json:"songId"`
			Songname   string  `json:"songname"`
			Difficulty float64 `json:"difficulty"`
			Rks        float64 `json:"rks"`
		} `json:"bests"`
		PlayerID          string  `json:"PlayerID"`
		ChallengeModeRank int     `json:"ChallengeModeRank"`
		RankingScore      float64 `json:"RankingScore"`
	} `json:"content"`
}

// TODO: Use Original Phigros Index Path

// 绘制平行四边形 angle 角度 x, y 坐标 w 宽度 l 斜边长 (github.com/Jiang-red/go-phigros-b19 function.)
func drawTriAngle(canvas *gg.Context, angle, x, y, w, l float64) {
	// 左上角为原点
	x0, y0 := x, y
	// 右上角
	x1, y1 := x+w, y
	// 右下角
	x2 := x1 - (l * (math.Cos(angle * math.Pi / 180.0)))
	y2 := y1 + (l * (math.Sin(angle * math.Pi / 180.0)))
	// 左下角
	x3, y3 := x2-w, y2
	canvas.NewSubPath()
	canvas.MoveTo(x0, y0)
	canvas.LineTo(x1, y1)
	canvas.LineTo(x2, y2)
	canvas.LineTo(x3, y3)
	canvas.ClosePath()
}

// CardRender Render By Original Image,so it's single work and do slowly.
func CardRender(canvas *gg.Context, dataOrigin []byte) *gg.Context {
	var referceLength = 0
	var referceWidth = 300 // + 1280
	var isRight = false
	var i, renderPath int
	// background render path.
	_ = json.Unmarshal(dataOrigin, &phigrosB19)
	i = 0
	renderPath = 0
	isPhi := phigrosB19.Content.Phi
	if isPhi { // while render the first set this change To Phi
		renderImage := phigrosB19.Content.BestList[i].SongId
		getRenderImage := background + renderImage + ".png"
		getImage, _ := gg.LoadImage(getRenderImage)
		// get background
		cardBackGround := DrawParallelogram(getImage)
		canvas.DrawImage(cardBackGround, referceWidth, 800+referceLength)
		// draw score path
		canvas.SetRGBA255(0, 0, 0, 160)
		drawTriAngle(canvas, 77, float64(referceWidth+500), float64(referceLength+850), 500, 210)
		canvas.Fill()
		// draw white line.
		canvas.SetColor(color.White)
		drawTriAngle(canvas, 77, float64(referceWidth+1000), float64(referceLength+850), 6, 210)
		canvas.Fill()
		// draw number format
		canvas.SetRGBA255(255, 255, 255, 245)
		drawTriAngle(canvas, 77, float64(referceWidth-26), float64(referceLength+800), 90, 55)
		canvas.Fill()
		// draw number path
		canvas.SetColor(color.Black)
		_ = canvas.LoadFontFace(font, 35)
		canvas.DrawString("Phi", float64(referceWidth-20), float64(referceLength+840))
		canvas.Fill()
		// render Diff.
		getDiff := phigrosB19.Content.BestList[i].Level
		SetDiffColor(getDiff, canvas)
		drawTriAngle(canvas, 77, float64(referceWidth-65), float64(referceLength+992), 150, 80)
		canvas.Fill()
		// render Text
		getRKS := strconv.FormatFloat(phigrosB19.Content.BestList[i].Difficulty, 'f', 1, 64)
		getRating := strconv.FormatFloat(phigrosB19.Content.BestList[i].Rks, 'f', 2, 64)
		canvas.SetColor(color.White)
		_ = canvas.LoadFontFace(font, 30)
		canvas.DrawString(getDiff+" "+getRKS, float64(referceWidth-50), float64(referceLength+1019))
		canvas.Fill()
		canvas.SetColor(color.White)
		_ = canvas.LoadFontFace(font, 45)
		canvas.DrawString(getRating, float64(referceWidth-60), float64(referceLength+1062))
		canvas.Fill()
		// render info (acc,score,name)
		getRankLink := GetRank(phigrosB19.Content.BestList[i].Score, phigrosB19.Content.BestList[i].Fc)
		loadRankImage, _ := gg.LoadImage(rank + getRankLink + ".png")
		canvas.DrawImage(loadRankImage, referceWidth+500, referceLength+920)
		canvas.Fill()
		getName := phigrosB19.Content.BestList[i].Songname
		_ = canvas.LoadFontFace(font, 35)
		canvas.DrawStringAnchored(getName, float64(referceWidth+740), float64(referceLength+890), 0.5, 0.5)
		canvas.SetColor(color.White)
		canvas.Fill()
		getScore := strconv.Itoa(phigrosB19.Content.BestList[i].Score)
		_ = canvas.LoadFontFace(font, 50)
		canvas.DrawStringAnchored(getScore, float64(referceWidth+740), float64(referceLength+950), 0.5, 0.5)
		canvas.Fill()
		canvas.SetColor(color.White)
		canvas.SetLineWidth(4)
		canvas.DrawLine(float64(referceWidth+630), float64(referceLength+990), float64(referceWidth+890), float64(referceLength+990))
		canvas.Stroke()
		_ = canvas.LoadFontFace(font, 40)
		getAcc := strconv.FormatFloat(phigrosB19.Content.BestList[i].Acc, 'f', 2, 64) + "%"
		canvas.DrawStringAnchored(getAcc, float64(referceWidth+760), float64(referceLength+1020), 0.5, 0.5)
		canvas.Fill()
		// width = referceWidth | height = 800+ referenceLength
		referceWidth = referceWidth + 1280
		referceLength += 75
		isRight = true
		i = i + 1
	}
	for ; i < len(phigrosB19.Content.BestList); i++ {
		renderImage := phigrosB19.Content.BestList[i].SongId
		getRenderImage := background + renderImage + ".png"
		getImage, _ := gg.LoadImage(getRenderImage)
		// get background
		cardBackGround := DrawParallelogram(getImage)
		canvas.DrawImage(cardBackGround, referceWidth, 800+referceLength)
		// draw score path
		canvas.SetRGBA255(0, 0, 0, 160)
		drawTriAngle(canvas, 77, float64(referceWidth+500), float64(referceLength+850), 500, 210)
		canvas.Fill()
		// draw white line.
		canvas.SetColor(color.White)
		drawTriAngle(canvas, 77, float64(referceWidth+1000), float64(referceLength+850), 6, 210)
		canvas.Fill()
		// draw number format
		canvas.SetRGBA255(255, 255, 255, 245)
		drawTriAngle(canvas, 77, float64(referceWidth-26), float64(referceLength+800), 90, 55)
		canvas.Fill()
		// draw number path
		canvas.SetColor(color.Black)
		_ = canvas.LoadFontFace(font, 35)
		canvas.DrawString("#"+strconv.Itoa(renderPath+1), float64(referceWidth-20), float64(referceLength+840))
		canvas.Fill()
		renderPath = renderPath + 1
		// render Diff.
		getDiff := phigrosB19.Content.BestList[i].Level
		SetDiffColor(getDiff, canvas)
		drawTriAngle(canvas, 77, float64(referceWidth-65), float64(referceLength+992), 150, 80)
		canvas.Fill()
		// render Text
		getRKS := strconv.FormatFloat(phigrosB19.Content.BestList[i].Difficulty, 'f', 1, 64)
		getRating := strconv.FormatFloat(phigrosB19.Content.BestList[i].Rks, 'f', 2, 64)
		canvas.SetColor(color.White)
		_ = canvas.LoadFontFace(font, 30)
		canvas.DrawString(getDiff+" "+getRKS, float64(referceWidth-50), float64(referceLength+1019))
		canvas.Fill()
		canvas.SetColor(color.White)
		_ = canvas.LoadFontFace(font, 45)
		canvas.DrawString(getRating, float64(referceWidth-60), float64(referceLength+1062))
		canvas.Fill()
		// render info (acc,score,name)
		getRankLink := GetRank(phigrosB19.Content.BestList[i].Score, phigrosB19.Content.BestList[i].Fc)
		loadRankImage, _ := gg.LoadImage(rank + getRankLink + ".png")
		canvas.DrawImage(loadRankImage, referceWidth+500, referceLength+920)
		canvas.Fill()
		getName := phigrosB19.Content.BestList[i].Songname
		_ = canvas.LoadFontFace(font, 35)
		canvas.DrawStringAnchored(getName, float64(referceWidth+740), float64(referceLength+890), 0.5, 0.5)
		canvas.SetColor(color.White)
		canvas.Fill()
		getScore := strconv.Itoa(phigrosB19.Content.BestList[i].Score)
		_ = canvas.LoadFontFace(font, 50)
		canvas.DrawStringAnchored(getScore, float64(referceWidth+740), float64(referceLength+950), 0.5, 0.5)
		canvas.Fill()
		canvas.SetColor(color.White)
		canvas.SetLineWidth(4)
		canvas.DrawLine(float64(referceWidth+630), float64(referceLength+990), float64(referceWidth+890), float64(referceLength+990))
		canvas.Stroke()
		_ = canvas.LoadFontFace(font, 40)
		getAcc := strconv.FormatFloat(phigrosB19.Content.BestList[i].Acc, 'f', 2, 64) + "%"
		canvas.DrawStringAnchored(getAcc, float64(referceWidth+760), float64(referceLength+1020), 0.5, 0.5)
		canvas.Fill()
		// width = referceWidth | height = 800+ referenceLength
		if !isRight {
			referceWidth = referceWidth + 1280
			referceLength += 75
			isRight = true
		} else {
			referceWidth = referceWidth - 1280
			referceLength -= 75
			isRight = false
			referceLength += 400
		}
	}
	canvas.Fill()
	return canvas
}

// GetUserChallengeMode Challenge Mode Type Reply
func GetUserChallengeMode(num int) (challenge string, link string) {
	var colors string
	for i := 1; i < 7; i++ {
		if i*100 > num {
			getCurrentRankLevel := i - 1
			switch {
			case getCurrentRankLevel == 1:
				colors = "Green"
				link = strconv.Itoa(num - (i-1)*100)
			case getCurrentRankLevel == 2:
				colors = "Blue"
				link = strconv.Itoa(num - (i-1)*100)
			case getCurrentRankLevel == 3:
				colors = "Red"
				link = strconv.Itoa(num - (i-1)*100)
			case getCurrentRankLevel == 4:
				colors = "Gold"
				link = strconv.Itoa(num - (i-1)*100)
			case getCurrentRankLevel == 5:
				colors = "Rainbow"
				link = strconv.Itoa(num - (i-1)*100)
			default:
				colors = ""
				link = "无记录"
			}
			break
		}
	}
	return colors, link
}

// DrawParallelogram Draw Card TriAnglePath
func DrawParallelogram(img image.Image) image.Image {
	length := 690.0
	dc := gg.NewContext(img.Bounds().Dx(), img.Bounds().Dy())
	picLengthWidth := img.Bounds().Dx()
	picLengthHeight := img.Bounds().Dy()
	getClipWidth := float64(picLengthWidth) - (length * 0.65) // get start point
	dc.NewSubPath()
	dc.MoveTo(getClipWidth, 0)
	dc.LineTo(float64(picLengthWidth), 0)
	dc.LineTo(length*0.65, float64(picLengthHeight))
	dc.LineTo(0, float64(picLengthHeight))
	dc.ClosePath()
	dc.Clip()
	dc.DrawImage(img, 0, 0)
	//getResizeImage := imaging.Resize(dc.Image(), 350, 270, imaging.Lanczos)
	getResizeImage := dc.Image()
	return getResizeImage
}

// GetRank get this rank.
func GetRank(num int, isFC bool) string {
	var rankPicName string

	switch {
	case num == 1000000:
		rankPicName = "phi"
	case num >= 960000:
		rankPicName = "v"
	case num >= 920000:
		rankPicName = "s"
	case num >= 880000:
		rankPicName = "a"
	case num >= 820000:
		rankPicName = "b"
	case num >= 700000:
		rankPicName = "c"
	default:
		rankPicName = "f"
	}
	if isFC && num != 1000000 {
		rankPicName = "fc"
	}
	return rankPicName
}

// SetDiffColor Set Diff Color.
func SetDiffColor(diff string, canvas *gg.Context) {
	switch {
	case diff == "IN":
		canvas.SetRGB255(189, 45, 36)
	case diff == "HD":
		canvas.SetRGB255(50, 115, 179)
	case diff == "AT":
		canvas.SetRGB255(56, 56, 56)
	case diff == "EZ":
		canvas.SetRGB255(79, 200, 134)
	}
}
