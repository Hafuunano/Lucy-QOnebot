package arc

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	arcaeaRes = "/root/Lucy_Project/workon/Lucy/data/arcaea"
)

// TODO:
// 1. optimize code(load font before running.)
// 2. run faster(go func)
// 3. minimize memory usage
// 4. SQL database support.
// 5. beautify b30.

type arcaea struct {
	Status  int `json:"status"`
	Content struct {
		Best30Avg   float64 `json:"best30_avg"`
		Recent10Avg float64 `json:"recent10_avg"`
		AccountInfo struct {
			Code                   string `json:"code"`
			Name                   string `json:"name"`
			UserId                 int    `json:"user_id"`
			IsMutual               bool   `json:"is_mutual"`
			IsCharUncappedOverride bool   `json:"is_char_uncapped_override"`
			IsCharUncapped         bool   `json:"is_char_uncapped"`
			IsSkillSealed          bool   `json:"is_skill_sealed"`
			Rating                 int    `json:"rating"`
			JoinDate               int64  `json:"join_date"`
			Character              int    `json:"character"`
		} `json:"account_info"`
		Best30List []struct {
			Score             int     `json:"score"`
			Health            int     `json:"health"`
			Rating            float64 `json:"rating"`
			SongId            string  `json:"song_id"`
			Modifier          int     `json:"modifier"`
			Difficulty        int     `json:"difficulty"`
			ClearType         int     `json:"clear_type"`
			BestClearType     int     `json:"best_clear_type"`
			TimePlayed        int64   `json:"time_played"`
			NearCount         int     `json:"near_count"`
			MissCount         int     `json:"miss_count"`
			PerfectCount      int     `json:"perfect_count"`
			ShinyPerfectCount int     `json:"shiny_perfect_count"`
		} `json:"best30_list"`
		Best30Overflow []struct {
			Score             int     `json:"score"`
			Health            int     `json:"health"`
			Rating            float64 `json:"rating"`
			SongId            string  `json:"song_id"`
			Modifier          int     `json:"modifier"`
			Difficulty        int     `json:"difficulty"`
			ClearType         int     `json:"clear_type"`
			BestClearType     int     `json:"best_clear_type"`
			TimePlayed        int64   `json:"time_played"`
			NearCount         int     `json:"near_count"`
			MissCount         int     `json:"miss_count"`
			PerfectCount      int     `json:"perfect_count"`
			ShinyPerfectCount int     `json:"shiny_perfect_count"`
		} `json:"best30_overflow"`
		Best30Songinfo []struct {
			NameEn         string  `json:"name_en"`
			NameJp         string  `json:"name_jp"`
			Artist         string  `json:"artist"`
			Bpm            string  `json:"bpm"`
			BpmBase        float64 `json:"bpm_base"`
			Set            string  `json:"set"`
			SetFriendly    string  `json:"set_friendly"`
			Time           int     `json:"time"`
			Side           int     `json:"side"`
			WorldUnlock    bool    `json:"world_unlock"`
			RemoteDownload bool    `json:"remote_download"`
			Bg             string  `json:"bg"`
			Date           int     `json:"date"`
			Version        string  `json:"version"`
			Difficulty     int     `json:"difficulty"`
			Rating         int     `json:"rating"`
			Note           int     `json:"note"`
			ChartDesigner  string  `json:"chart_designer"`
			JacketDesigner string  `json:"jacket_designer"`
			JacketOverride bool    `json:"jacket_override"`
			AudioOverride  bool    `json:"audio_override"`
		} `json:"best30_songinfo"`
		Best30OverflowSonginfo []struct {
			NameEn         string  `json:"name_en"`
			NameJp         string  `json:"name_jp"`
			Artist         string  `json:"artist"`
			Bpm            string  `json:"bpm"`
			BpmBase        float64 `json:"bpm_base"`
			Set            string  `json:"set"`
			SetFriendly    string  `json:"set_friendly"`
			Time           int     `json:"time"`
			Side           int     `json:"side"`
			WorldUnlock    bool    `json:"world_unlock"`
			RemoteDownload bool    `json:"remote_download"`
			Bg             string  `json:"bg"`
			Date           int     `json:"date"`
			Version        string  `json:"version"`
			Difficulty     int     `json:"difficulty"`
			Rating         int     `json:"rating"`
			Note           int     `json:"note"`
			ChartDesigner  string  `json:"chart_designer"`
			JacketDesigner string  `json:"jacket_designer"`
			JacketOverride bool    `json:"jacket_override"`
			AudioOverride  bool    `json:"audio_override"`
		} `json:"best30_overflow_songinfo"`
	} `json:"content"`
}

// GetSongCurrentLocation Get Location (Needs to change the main location due to it uses different location when using other platform.)
func GetSongCurrentLocation(r arcaea, idLocated int, b40 bool) (currentSonglocation string) {
	// get the song downloaded status
	var dl bool
	var songID string
	if b40 {
		dl = r.Content.Best30OverflowSonginfo[idLocated].RemoteDownload
		songID = r.Content.Best30Overflow[idLocated].SongId
	} else {
		dl = r.Content.Best30Songinfo[idLocated].RemoteDownload
		songID = r.Content.Best30List[idLocated].SongId
	}
	// get song's current location.
	if dl {
		currentSonglocation = arcaeaRes + "/assets/song/dl_" + songID + "/base.jpg"
	} else {
		currentSonglocation = arcaeaRes + "/assets/song/" + songID + "/base.jpg"
	}

	return currentSonglocation
}

// GetAverageColor different from k-means algorithm,it uses origin plugin's algorithm.
func GetAverageColor(image image.Image) (int, int, int) {
	var RList []int
	var GList []int
	var BList []int
	width, height := image.Bounds().Size().X, image.Bounds().Size().Y
	for x := 0; x < width/5; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := image.At(x, y).RGBA()
			RList = append(RList, int(r>>8))
			GList = append(GList, int(g>>8))
			BList = append(BList, int(b>>8))
		}
	}
	RAverage := int(average(RList))
	GAverage := int(average(GList))
	BAverage := int(average(BList))
	return RAverage, GAverage, BAverage
}

// average sum all the numbers and divide by the length of the list.
func average(numbers []int) float64 {
	var sum float64
	for _, num := range numbers {
		sum += float64(num)
	}
	return math.Round(sum / float64(len(numbers)))
}

// IsDark judge which font color I prefer to use (black or white)
func IsDark(rgb color.RGBA) bool {
	var (
		r = float64(rgb.R) * 0.299
		g = float64(rgb.G) * 0.587
		b = float64(rgb.B) * 0.114
	)
	lum := r + g + b
	return lum < 192
}

// FastResizeImage resize the image to the size I want.(Use imgfactory to resize the image)
func FastResizeImage(imgPath string, w int, h int) image.Image {
	file, err := os.ReadFile(imgPath)
	if err != nil {
		panic(err)
	}
	// decode img
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}
	// resize img
	return imgfactory.Limit(img, w, h)
}

// FormatNumber Format Num like (000'000'000)
func FormatNumber(number int) string {
	if number < 1000 {
		return fmt.Sprintf("%d", number)
	}
	return FormatNumber(number/1000) + fmt.Sprintf("'%03d", number%1000)
}

// FormatTimeStamp Format TimeStamp to 2006-01-02 15:04:05
func FormatTimeStamp(timeStamp int64) string {
	t := time.Unix(timeStamp/1000, 0)
	return t.Format("2006-01-02 15:04:05")
}

// DrawScoreCard draw the detailed info of the score.
func DrawScoreCard(songCover image.Image, songNum int, r arcaea) image.Image {
	// can optimize this part, I just copy it from the original plugin(arcaeabot)
	var diffs = make(map[int]string)
	diffs[0] = "PST.png"
	diffs[1] = "PRS.png"
	diffs[2] = "FTR.png"
	diffs[3] = "BYD.png"
	scRaw := gg.NewContextForImage(songCover)
	sc := gg.NewContext(scRaw.W()*5/2, scRaw.H())
	sc.DrawRoundedRectangle(0, 0, float64(sc.W()), float64(sc.H()), 32)
	sc.Clip()
	R, G, B := GetAverageColor(songCover)
	sc.DrawImage(songCover, sc.W()-scRaw.W(), 0)
	grad := gg.NewLinearGradient(0, 0, float64(sc.W()), 0)
	grad.AddColorStop(0, color.NRGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: 255})
	grad.AddColorStop(1-float64(scRaw.W())/float64(sc.W()), color.NRGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: 255})
	grad.AddColorStop(1, color.NRGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: 0})
	sc.SetFillStyle(grad)
	sc.DrawRectangle(0, 0, float64(sc.W()), float64(sc.H()))
	sc.Fill()
	songCoverHandled := sc.Image()
	diffNum := r.Content.Best30List[songNum].Difficulty
	diff := diffs[diffNum]
	fullDiffLink := arcaeaRes + "/resource/diff/" + diff
	resizedDiffImage := FastResizeImage(fullDiffLink, 14, 58)
	// check the song name's length.
	getSongName := r.Content.Best30Songinfo[songNum].NameEn
	if len(getSongName) >= 19 {
		getSongName = getSongName[:18] + "..."
	}
	mainPicHandler := gg.NewContextForImage(songCoverHandled)
	// Handle font (should deal it before the function runs.)
	fontByteRaw, _ := os.ReadFile(arcaeaRes + "/resource/font/NotoSansCJKsc-Regular.otf")
	exoMid, _ := os.ReadFile(arcaeaRes + "/resource/font/Exo-Medium.ttf")
	kazeRegular, _ := os.ReadFile(arcaeaRes + "/resource/font/Kazesawa-Regular.ttf")
	andrealFont, _ := os.ReadFile(arcaeaRes + "/resource/font/Andrea.ttf")
	fontExact, _ := opentype.Parse(fontByteRaw) // sans-serif regular font
	exoMidExact, _ := opentype.Parse(exoMid)
	kazeRegularExact, _ := opentype.Parse(kazeRegular)
	andrealFontExact, _ := opentype.Parse(andrealFont)
	sans, _ := opentype.NewFace(fontExact, &opentype.FaceOptions{Size: 30, DPI: 72, Hinting: font.HintingFull}) // sans-serif regular font
	exoMidFace, _ := opentype.NewFace(exoMidExact, &opentype.FaceOptions{Size: 40, DPI: 72, Hinting: font.HintingFull})
	exoSmallFace, _ := opentype.NewFace(exoMidExact, &opentype.FaceOptions{Size: 25, DPI: 80, Hinting: font.HintingFull})
	AndrealFaceLl, _ := opentype.NewFace(andrealFontExact, &opentype.FaceOptions{Size: 50, DPI: 80, Hinting: font.HintingFull})
	kazeRegularFacel, _ := opentype.NewFace(kazeRegularExact, &opentype.FaceOptions{Size: 30, DPI: 72, Hinting: font.HintingFull})
	kazeRegularFaceSl, _ := opentype.NewFace(kazeRegularExact, &opentype.FaceOptions{Size: 25, DPI: 72, Hinting: font.HintingFull})
	kazeRegularFaceS, _ := opentype.NewFace(kazeRegularExact, &opentype.FaceOptions{Size: 20, DPI: 72, Hinting: font.HintingFull})
	// font part end.
	mainPicHandler.SetFontFace(sans)
	// check the bg is dark? if false, use black color font, else use white color font.
	var setMainColor color.NRGBA
	if IsDark(color.RGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: 255}) {
		setMainColor = color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: 205}
		mainPicHandler.SetColor(setMainColor)
	} else {
		setMainColor = color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: 205}
		mainPicHandler.SetColor(setMainColor)
	}
	mainPicHandler.DrawString(getSongName, 45, 52)                  // Write song name.
	mainPicHandler.DrawString("#"+strconv.Itoa(songNum+1), 580, 45) // Write nums
	mainPicHandler.Fill()
	mainPicHandler.SetFontFace(exoMidFace)
	mainPicHandler.DrawString(FormatNumber(r.Content.Best30List[songNum].Score), 45, 100) // draw score.
	mainPicHandler.Fill()
	mainPicHandler.DrawImage(resizedDiffImage, 24, 24) // diff path
	// origin plugin's code judge the bg's color and get average color again, but I think it's useless.(lmao)
	// draw andrea style.
	tableImage, _ := os.ReadFile(arcaeaRes + "/resource/b30/table.png")
	tableImageResized, _, _ := image.Decode(bytes.NewReader(tableImage)) // that's origin lmao
	mainPicHandler.DrawImage(tableImageResized, 0, -10)
	// draw p,f,l,etc.
	mainPicHandler.SetFontFace(exoSmallFace)
	mainPicHandler.SetColor(setMainColor)
	mainPicHandler.DrawString("P", 48, 144)
	mainPicHandler.DrawString("F", 48, 189)
	mainPicHandler.DrawString("L", 48, 234)
	mainPicHandler.Fill()
	// draw ptt,date,etc.
	mainPicHandler.SetFontFace(exoSmallFace)
	mainPicHandler.SetColor(setMainColor)
	mainPicHandler.DrawString("PTT", 250, 130)
	mainPicHandler.DrawString("DATE", 250, 200)
	mainPicHandler.Fill()
	mainPicHandler.SetFontFace(AndrealFaceLl)
	mainPicHandler.SetColor(setMainColor)
	mainPicHandler.DrawString(">", 330, 160)
	mainPicHandler.Fill()
	// draw json unmashalled data.
	// draw size 30 font(P,F,L,etc.)
	mainPicHandler.SetFontFace(kazeRegularFacel)
	mainPicHandler.SetColor(setMainColor)
	mainPicHandler.DrawString(strconv.Itoa(r.Content.Best30List[songNum].PerfectCount), 75, 142)
	mainPicHandler.DrawString(strconv.Itoa(r.Content.Best30List[songNum].NearCount), 75, 187)
	mainPicHandler.DrawString(strconv.Itoa(r.Content.Best30List[songNum].MissCount), 75, 232)
	// format time
	mainPicHandler.DrawString(FormatTimeStamp(r.Content.Best30List[songNum].TimePlayed), 250, 240)
	mainPicHandler.Fill()
	// draw size 20 font(shiny count).
	mainPicHandler.SetFontFace(kazeRegularFaceS)
	mainPicHandler.SetColor(setMainColor)
	mainPicHandler.DrawString("+"+strconv.Itoa(r.Content.Best30List[songNum].ShinyPerfectCount), 153, 140)
	mainPicHandler.Fill()
	// draw ptt.
	mainPicHandler.SetFontFace(kazeRegularFaceSl)
	mainPicHandler.SetColor(setMainColor)
	mainPicHandler.DrawString(strconv.FormatFloat(float64(r.Content.Best30Songinfo[songNum].Rating)/10, 'f', 1, 64), 250, 162)
	mainPicHandler.DrawString(strconv.FormatFloat(r.Content.Best30List[songNum].Rating, 'f', 3, 64), 360, 155)
	mainPicHandler.Fill()
	return mainPicHandler.Image()
}

// GetAvatar and resize it to 250x250.
func GetAvatar(r arcaea) (avatarByte image.Image) {
	getUncappedStatus := r.Content.AccountInfo.IsCharUncapped
	getUncappedOverrideStatus := r.Content.AccountInfo.IsCharUncappedOverride
	getCharacterNum := r.Content.AccountInfo.Character
	// get avatar
	if getUncappedStatus != getUncappedOverrideStatus {
		// uncapped
		avatarLocation, _ := os.ReadFile(arcaeaRes + "/assets/char/" + strconv.Itoa(getCharacterNum) + "u_icon.png")
		avatarBytes, _, _ := image.Decode(bytes.NewReader(avatarLocation))
		avatarByte = imgfactory.Size(avatarBytes, 250, 250).Image()
		return
	} else {
		avatarLocation, _ := os.ReadFile(arcaeaRes + "/assets/char/" + strconv.Itoa(getCharacterNum) + "_icon.png")
		avatarBytes, _, _ := image.Decode(bytes.NewReader(avatarLocation))
		avatarByte = imgfactory.Size(avatarBytes, 250, 250).Image()
		return
	}
}

// ChoicePttBackground choice ptt background.
func choicePttBackground(ptt float64) string {
	if ptt == -1 {
		return "rating_off.png"
	}
	ptt /= 100
	switch {
	case ptt < 3:
		return "rating_0.png"
	case ptt < 7:
		return "rating_1.png"
	case ptt < 10:
		return "rating_2.png"
	case ptt < 11:
		return "rating_3.png"
	case ptt < 12:
		return "rating_4.png"
	case ptt < 12.5:
		return "rating_5.png"
	case ptt < 13:
		return "rating_6.png"
	default:
		return "rating_7.png"
	}
}

// DrawMainUserB30 draw main user b30.
func DrawMainUserB30(mainBg image.Image, r arcaea) image.Image {
	// main pg user b30
	exoMid, _ := os.ReadFile(arcaeaRes + "/resource/font/Exo-Medium.ttf")
	exoMidExact, _ := opentype.Parse(exoMid)
	mainBGHandler := gg.NewContextForImage(mainBg)
	// draw avatar
	mainBGHandler.DrawImage(GetAvatar(r), 75, 130)
	// draw name
	exoMidFaceLL, _ := opentype.NewFace(exoMidExact, &opentype.FaceOptions{Size: 100, DPI: 72, Hinting: font.HintingFull})
	mainBGHandler.SetFontFace(exoMidFaceLL)
	mainBGHandler.SetColor(color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: 255})
	mainBGHandler.DrawString(r.Content.AccountInfo.Name, 355, 250)
	mainBGHandler.FillPreserve()
	// draw userId
	exoMidFaceLs, _ := opentype.NewFace(exoMidExact, &opentype.FaceOptions{Size: 60, DPI: 72, Hinting: font.HintingFull})
	mainBGHandler.SetFontFace(exoMidFaceLs)
	mainBGHandler.DrawString("ID:"+r.Content.AccountInfo.Code, 380, 364)
	mainBGHandler.FillPreserve()
	// draw ptt.
	pttImage, _ := os.ReadFile(arcaeaRes + "/resource/ptt/" + choicePttBackground(float64(r.Content.AccountInfo.Rating)))
	pttImageFormat, _, _ := image.Decode(bytes.NewReader(pttImage))
	exoMidFace, _ := opentype.NewFace(exoMidExact, &opentype.FaceOptions{Size: 45, DPI: 72, Hinting: font.HintingFull})
	mainBGHandler.SetFontFace(exoMidFace)
	mainBGHandler.DrawImage(imgfactory.Size(pttImageFormat, 150, 150).Image(), 195, 295)
	mainBGHandler.SetColor(color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: 255})
	mainBGHandler.DrawStringAnchored(strconv.FormatFloat(float64(r.Content.AccountInfo.Rating)/100, 'f', 2, 64), 270, 365, 0.5, 0.5)
	mainBGHandler.FillPreserve()
	mainBGHandler.SetFontFace(exoMidFaceLL)
	mainBGHandler.DrawString("Best 30:"+strconv.FormatFloat(r.Content.Best30Avg, 'f', 3, 64), 200, 560)
	mainBGHandler.DrawString("Recent 10:"+strconv.FormatFloat(r.Content.Recent10Avg, 'f', 3, 64), 1000, 560)
	mainBGHandler.Fill()
	return mainBGHandler.Image()
}

// FinishedFullB30 draw full b30.
func FinishedFullB30(mainB30 image.Image, r arcaea) image.Image {
	// main pg user b30
	mainB30Handler := gg.NewContextForImage(mainB30)
	// draw b30
	initWidth := 100
	initHighth := 640
	for i := 0; i < len(r.Content.Best30List); i++ {
		getSongLocation := GetSongCurrentLocation(r, i, false)
		// get song's average color
		bgReader, err := os.ReadFile(getSongLocation)
		if err != nil {
			panic(err)
		}
		bgDecode, _, _ := image.Decode(bytes.NewReader(bgReader))
		mainB30Handler.DrawImage(imgfactory.Size(DrawScoreCard(imgfactory.Limit(bgDecode, 270, 270), i, r), 560, 256).Image(), initWidth, initHighth)
		mainB30Handler.Fill()
		initWidth += 60 + 560
		if initWidth >= 1860 {
			initWidth = 100
			initHighth += 60 + 224
		}
	}
	if len(r.Content.Best30List) >= 30 {
		divider := FastResizeImage(arcaeaRes+"/resource/b30/Divider.png", 2000, 500)
		mainB30Handler.DrawImage(divider, 300, 3540)
		for i := 0; i < 9; i++ {
			getSongLocation := GetSongCurrentLocation(r, i, true)
			// get song's average color
			bgReader, err := os.ReadFile(getSongLocation)
			if err != nil {
				panic(err)
			}
			bgDecode, _, _ := image.Decode(bytes.NewReader(bgReader))
			initB10Length := initHighth + 180
			mainB30Handler.DrawImage(imgfactory.Size(DrawScoreCard(imgfactory.Limit(bgDecode, 270, 270), i, r), 560, 256).Image(), initWidth, initB10Length)
			mainB30Handler.Fill()
			initWidth += 60 + 560
			if initWidth >= 1860 {
				initWidth = 100
				initHighth += 60 + 224
			}
		}
	}
	return mainB30Handler.Image()
}
