package arc

import (
	"bytes"
	"fmt"
	"github.com/MoYoez/img_toolbox/FrostedGlassLike"
	"github.com/MoYoez/img_toolbox/imgCutter"
	"image"
	"image/color"
	"math"
	"os"
	"strconv"
	"time"
	"unicode"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	arcaeaRes         = "/root/Lucy_Project/workon/Lucy/data/arcaea"
	diff              string
	diffNum           int
	getSongName       string
	rating            string
	sans              font.Face
	exoMidFaces       font.Face
	exoSmallFace      font.Face
	AndrealFaceLl     font.Face
	AndrealFaceL      font.Face
	kazeRegularFacel  font.Face
	kazeRegularFaceSl font.Face
	kazeRegularFaceS  font.Face
	exoMidFace        font.Face
	exoMidFaceLs      font.Face
	exoMidFaceLL      font.Face
	exoSemiBoldFace   font.Face
	exoSmallFaceSS    font.Face
)

// TODO:
// 1. run faster(go func)
// 2. minimize memory usage
// 3. beautify b30.

type user struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Content struct {
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
		RecentScore []struct {
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
		} `json:"recent_score"`
		Songinfo []struct {
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
		} `json:"song_info"`
	} `json:"content"`
}

type arcaea struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
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
		} `json:"best30_song_info"`
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
		} `json:"best30_overflow_song_info"`
	} `json:"content"`
}

func init() {
	// main pg user b30
	// Handle font (should deal it before the function runs.)
	sans = LoadFontFace(arcaeaRes+"/resource/font/NotoSansCJKsc-Regular.otf", 30, 72) // sans-serif regular font
	exoMidFaces = LoadFontFace(arcaeaRes+"/resource/font/Exo-Medium.ttf", 40, 72)
	exoSmallFace = LoadFontFace(arcaeaRes+"/resource/font/Exo-Medium.ttf", 25, 80)
	exoSmallFaceSS = LoadFontFace(arcaeaRes+"/resource/font/Exo-Medium.ttf", 15, 72)
	AndrealFaceLl = LoadFontFace(arcaeaRes+"/resource/font/Andrea.ttf", 50, 80)
	AndrealFaceL = LoadFontFace(arcaeaRes+"/resource/font/Andrea.ttf", 30, 80)
	kazeRegularFacel = LoadFontFace(arcaeaRes+"/resource/font/Kazesawa-Regular.ttf", 30, 72)
	kazeRegularFaceSl = LoadFontFace(arcaeaRes+"/resource/font/Kazesawa-Regular.ttf", 25, 72)
	kazeRegularFaceS = LoadFontFace(arcaeaRes+"/resource/font/Kazesawa-Regular.ttf", 20, 72)
	exoMidFace = LoadFontFace(arcaeaRes+"/resource/font/Exo-Medium.ttf", 45, 72)
	exoMidFaceLs = LoadFontFace(arcaeaRes+"/resource/font/Exo-Medium.ttf", 60, 72)
	exoMidFaceLL = LoadFontFace(arcaeaRes+"/resource/font/Exo-Medium.ttf", 100, 72)
	exoSemiBoldFace = LoadFontFace(arcaeaRes+"/resource/font/Exo-SemiBold.ttf", 20, 72)
}

// StrokeStringWithOtherColor make the text shade?
func StrokeStringWithOtherColor(text string, color color.Color, h float64, w float64, bg *gg.Context) {
	bg.SetColor(color)
	for i := 1; i < 3; i++ {
		for j := 1; j < 3; j++ {
			bg.DrawString(text, h-float64(i), w-float64(j))
		}
	}
}

// GetAvatarForGetUserInfo and resize it to 90x90.
func GetAvatarForGetUserInfo(userinfo user) (avatarByte image.Image) {
	getUncappedStatus := userinfo.Content.AccountInfo.IsCharUncapped
	getUncappedOverrideStatus := userinfo.Content.AccountInfo.IsCharUncappedOverride
	getCharacterNum := userinfo.Content.AccountInfo.Character
	// get avatar
	if getUncappedStatus != getUncappedOverrideStatus {
		// uncapped
		avatarLocation, _ := os.ReadFile(arcaeaRes + "/assets/char/" + strconv.Itoa(getCharacterNum) + "u_icon.png")
		avatarBytes, _, _ := image.Decode(bytes.NewReader(avatarLocation))
		avatarByte = imgfactory.Size(avatarBytes, 90, 90).Image()
		return
	} else {
		avatarLocation, _ := os.ReadFile(arcaeaRes + "/assets/char/" + strconv.Itoa(getCharacterNum) + "_icon.png")
		avatarBytes, _, _ := image.Decode(bytes.NewReader(avatarLocation))
		avatarByte = imgfactory.Size(avatarBytes, 90, 90).Image()
		return
	}
}

// GetSongCurrentLocation Get Location (Needs to change the main location due to it uses different location when using other platform.)  you need to Unmarshal json file first.
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

// GetB30RecentSongLocation find recent song location, you need to Unmarshal json file first.
func GetB30RecentSongLocation(userinfo user) (currentSonglocation string) {
	var (
		songID string
		dl     bool
	)
	songID = userinfo.Content.RecentScore[0].SongId
	dl = userinfo.Content.Songinfo[0].RemoteDownload
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
	RAverage := int(Average(RList))
	GAverage := int(Average(GList))
	BAverage := int(Average(BList))
	return RAverage, GAverage, BAverage
}

// Average sum all the numbers and divide by the length of the list.
func Average(numbers []int) float64 {
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
func DrawScoreCard(songCover image.Image, songNum int, r arcaea, b40 bool) image.Image {
	// can optimize this part, I just copy it from the original plugin(arcaeabot)
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
	if b40 {
		diffNum = r.Content.Best30Overflow[songNum].Difficulty
	} else {
		diffNum = r.Content.Best30List[songNum].Difficulty
	}
	switch {
	case diffNum == 0:
		diff = "PST.png"
	case diffNum == 1:
		diff = "PRS.png"
	case diffNum == 2:
		diff = "FTR.png"
	case diffNum == 3:
		diff = "BYD.png"
	}
	fullDiffLink := arcaeaRes + "/resource/diff/" + diff
	resizedDiffImage := FastResizeImage(fullDiffLink, 14, 58)
	// check the song name's length.
	if b40 {
		getSongName = r.Content.Best30OverflowSonginfo[songNum].NameEn
	} else {
		getSongName = r.Content.Best30Songinfo[songNum].NameEn
	}
	if len(getSongName) >= 19 {
		getSongName = getSongName[:18] + "..."
	}
	mainPicHandler := gg.NewContextForImage(songCoverHandled)
	// Handle font (should deal it before the function runs.)
	// font part end.
	mainPicHandler.SetFontFace(sans)
	StrokeStringWithOtherColor("#"+strconv.Itoa(songNum+1), color.Black, 590, 45, mainPicHandler)
	mainPicHandler.SetColor(color.White)
	mainPicHandler.DrawString("#"+strconv.Itoa(songNum+1), 590, 45) // Write nums
	// check the bg is dark? if false, use black color font, else use white color font.
	var setMainColor color.NRGBA
	if IsDark(color.RGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: 255}) {
		setMainColor = color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: 205}
		mainPicHandler.SetColor(setMainColor)
	} else {
		setMainColor = color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: 205}
		mainPicHandler.SetColor(setMainColor)
	}
	mainPicHandler.DrawString(getSongName, 45, 52) // Write song name.
	mainPicHandler.FillPreserve()
	mainPicHandler.SetFontFace(exoMidFaces)
	if b40 {
		mainPicHandler.DrawString(FormatNumber(r.Content.Best30Overflow[songNum].Score), 45, 100) // draw score.
	} else {
		mainPicHandler.DrawString(FormatNumber(r.Content.Best30List[songNum].Score), 45, 100) // draw score.
	}
	mainPicHandler.FillPreserve()
	mainPicHandler.DrawImage(resizedDiffImage, 24, 24) // diff path
	// origin plugin's code judge the bg's color and get Average color again, but I think it's useless.(lmao)
	// draw andrea style.
	tableImage, _ := os.ReadFile(arcaeaRes + "/resource/b30/table.png")
	tableImageResized, _, _ := image.Decode(bytes.NewReader(tableImage)) // that's origin lmao
	mainPicHandler.DrawImage(tableImageResized, 0, -10)
	// draw p,f,l,etc.
	mainPicHandler.SetFontFace(exoSmallFace)
	mainPicHandler.DrawString("P", 48, 144)
	mainPicHandler.DrawString("F", 48, 189)
	mainPicHandler.DrawString("L", 48, 234)
	mainPicHandler.DrawString("PTT", 250, 130)
	mainPicHandler.DrawString("DATE", 250, 200)
	mainPicHandler.FillPreserve()
	mainPicHandler.SetFontFace(AndrealFaceLl)
	mainPicHandler.DrawString(">", 330, 160)
	mainPicHandler.FillPreserve()
	// draw json unmashalled data.
	// draw size 30 font(P,F,L,etc.)
	if b40 {
		mainPicHandler.SetFontFace(kazeRegularFacel)
		mainPicHandler.DrawString(strconv.Itoa(r.Content.Best30Overflow[songNum].PerfectCount), 75, 142)
		mainPicHandler.DrawString(strconv.Itoa(r.Content.Best30Overflow[songNum].NearCount), 75, 187)
		mainPicHandler.DrawString(strconv.Itoa(r.Content.Best30Overflow[songNum].MissCount), 75, 232)
		// format time
		mainPicHandler.DrawString(FormatTimeStamp(r.Content.Best30Overflow[songNum].TimePlayed), 250, 240)
		mainPicHandler.FillPreserve()
		// draw size 20 font(shiny count).
		mainPicHandler.SetFontFace(kazeRegularFaceS)
		mainPicHandler.DrawString("+"+strconv.Itoa(r.Content.Best30Overflow[songNum].ShinyPerfectCount), 153, 140)
		mainPicHandler.FillPreserve()
		// draw ptt.
		mainPicHandler.SetFontFace(kazeRegularFaceSl)
		mainPicHandler.DrawString(strconv.FormatFloat(float64(r.Content.Best30OverflowSonginfo[songNum].Rating)/10, 'f', 1, 64), 250, 162)
		mainPicHandler.DrawString(strconv.FormatFloat(r.Content.Best30Overflow[songNum].Rating, 'f', 3, 64), 360, 155)
		mainPicHandler.Fill()
	} else {
		mainPicHandler.SetFontFace(kazeRegularFacel)
		mainPicHandler.DrawString(strconv.Itoa(r.Content.Best30List[songNum].PerfectCount), 75, 142)
		mainPicHandler.DrawString(strconv.Itoa(r.Content.Best30List[songNum].NearCount), 75, 187)
		mainPicHandler.DrawString(strconv.Itoa(r.Content.Best30List[songNum].MissCount), 75, 232)
		// format time
		mainPicHandler.DrawString(FormatTimeStamp(r.Content.Best30List[songNum].TimePlayed), 250, 240)
		mainPicHandler.FillPreserve()
		// draw size 20 font(shiny count).
		mainPicHandler.SetFontFace(kazeRegularFaceS)
		mainPicHandler.DrawString("+"+strconv.Itoa(r.Content.Best30List[songNum].ShinyPerfectCount), 153, 140)
		mainPicHandler.FillPreserve()
		// draw ptt.
		mainPicHandler.SetFontFace(kazeRegularFaceSl)
		mainPicHandler.DrawString(strconv.FormatFloat(float64(r.Content.Best30Songinfo[songNum].Rating)/10, 'f', 1, 64), 250, 162)
		mainPicHandler.DrawString(strconv.FormatFloat(r.Content.Best30List[songNum].Rating, 'f', 3, 64), 360, 155)
		mainPicHandler.Fill()
	}
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
func ChoicePttBackground(ptt float64) string {
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
	mainBGHandler := gg.NewContextForImage(mainBg)
	// draw avatar
	mainBGHandler.DrawImage(GetAvatar(r), 75, 130)
	// draw name
	mainBGHandler.SetFontFace(exoMidFaceLL)
	mainBGHandler.SetColor(color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: 255})
	mainBGHandler.DrawString(r.Content.AccountInfo.Name, 355, 250)
	mainBGHandler.FillPreserve()
	// draw userId
	mainBGHandler.SetFontFace(exoMidFaceLs)
	mainBGHandler.DrawString("ID:"+r.Content.AccountInfo.Code, 380, 364)
	mainBGHandler.FillPreserve()
	// draw ptt.
	pttImage, _ := os.ReadFile(arcaeaRes + "/resource/ptt/" + ChoicePttBackground(float64(r.Content.AccountInfo.Rating)))
	pttImageFormat, _, _ := image.Decode(bytes.NewReader(pttImage))
	mainBGHandler.SetFontFace(exoMidFace)
	mainBGHandler.DrawImage(imgfactory.Size(pttImageFormat, 150, 150).Image(), 195, 295)
	if r.Content.AccountInfo.Rating == -1 {
		rating = "--"
	} else {
		rating = strconv.FormatFloat(float64(r.Content.AccountInfo.Rating)/100, 'f', 2, 64)
	}
	StrokeStringWithOtherColor(rating, color.Black, 215, 385, mainBGHandler)
	mainBGHandler.SetColor(color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: 255})
	mainBGHandler.DrawString(rating, 215, 385)
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
	// TODO: Maybe I can try to not use resize for the Scorecards?
	for i := 0; i < len(r.Content.Best30List); i++ {
		getSongLocation := GetSongCurrentLocation(r, i, false)
		// get song's Average color
		bgReader, err := os.ReadFile(getSongLocation)
		if err != nil {
			panic(err)
		}
		bgDecode, _, _ := image.Decode(bytes.NewReader(bgReader))
		mainB30Handler.DrawImage(imgfactory.Size(DrawScoreCard(imgfactory.Limit(bgDecode, 270, 270), i, r, false), 560, 256).Image(), initWidth, initHighth)
		mainB30Handler.Fill()
		initWidth += 60 + 560
		if initWidth >= 1860 {
			initWidth = 100
			initHighth += 60 + 224
		}
	}
	if len(r.Content.Best30Overflow) > 0 {
		divider := FastResizeImage(arcaeaRes+"/resource/b30/Divider.png", 2000, 500)
		mainB30Handler.DrawImage(divider, 300, 3540)
		for i := 0; i < 9; i++ {
			getSongLocation := GetSongCurrentLocation(r, i, true)
			// get song's Average color
			bgReader, err := os.ReadFile(getSongLocation)
			if err != nil {
				panic(err)
			}
			bgDecode, _, _ := image.Decode(bytes.NewReader(bgReader))
			initB10Length := initHighth + 180
			mainB30Handler.DrawImage(imgfactory.Size(DrawScoreCard(imgfactory.Limit(bgDecode, 270, 270), i, r, true), 560, 256).Image(), initWidth, initB10Length)
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

// LoadFontFace load font face once before running, to work it quickly and save memory.
func LoadFontFace(filePath string, size float64, dpi float64) font.Face {
	fontFile, _ := os.ReadFile(filePath)
	fontFileParse, _ := opentype.Parse(fontFile)
	fontFace, _ := opentype.NewFace(fontFileParse, &opentype.FaceOptions{Size: size, DPI: dpi, Hinting: font.HintingFull})
	return fontFace
}

// IsAlphanumeric Check the context is num or english.
func IsAlphanumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}

// RenderUserRecentLog render user recent log.
func RenderUserRecentLog(userinfo user) image.Image {
	// first of all, get raw and format it.
	var rawFormatUserPic string
	if userinfo.Content.Songinfo[0].Side == 0 {
		rawFormatUserPic = arcaeaRes + "/resource/recent/RawV3Bg_0.png"
	} else {
		rawFormatUserPic = arcaeaRes + "/resource/recent/RawV3Bg_1.png"
	}
	dataRawResize := FastResizeImage(rawFormatUserPic, 540, 781)
	// frostlike bg
	bgLocation, err := os.ReadFile(GetB30RecentSongLocation(userinfo)) // get pic and render it as frost like(Blur.)
	imageBgLocation, _, _ := image.Decode(bytes.NewReader(bgLocation)) // decode it.
	resultimg := FrostedGlassLike.FrostedGlassLike(imageBgLocation, 10, err)
	handledImg := imgCutter.CropImage(resultimg, 600, 867)
	mainBG := gg.NewContextForImage(handledImg)
	mainBG.DrawImage(dataRawResize, 160, 86)
	// render image(e.g. Track Lost,ACC,so on.)
	avatarimg := GetAvatarForGetUserInfo(userinfo)
	mainBG.DrawImage(avatarimg, 220, 170)
	pttImage, _ := os.ReadFile(arcaeaRes + "/resource/ptt/" + ChoicePttBackground(float64(userinfo.Content.AccountInfo.Rating)))
	pttImageFormat, _, _ := image.Decode(bytes.NewReader(pttImage))
	mainBG.DrawImage(imgfactory.Size(pttImageFormat, 70, 70).Image(), 260, 210)
	if userinfo.Content.AccountInfo.Rating == -1 {
		rating = "--"
	} else {
		rating = strconv.FormatFloat(float64(userinfo.Content.AccountInfo.Rating)/100, 'f', 2, 64)
	}
	mainBG.SetFontFace(exoSemiBoldFace)
	// draw black, then draw white.
	StrokeStringWithOtherColor(rating, color.Black, 270, 252, mainBG)
	mainBG.SetColor(color.White)
	mainBG.DrawString(rating, 270, 252)
	mainBG.SetColor(color.Black)
	mainBG.DrawString("ArcID: "+userinfo.Content.AccountInfo.Code, 330, 240)
	mainBG.SetFontFace(sans)
	mainBG.DrawString(userinfo.Content.AccountInfo.Name, 330, 210)
	mainBG.DrawImage(imgfactory.Limit(imageBgLocation, 230, 230), 315, 306)
	mainBG.SetFontFace(AndrealFaceL)
	var userRecentSongName string
	userRecentSongName = userinfo.Content.Songinfo[0].NameJp
	if userRecentSongName == "" {
		userRecentSongName = userinfo.Content.Songinfo[0].NameEn
	}
	mainBG.DrawStringAnchored(userRecentSongName, 435, 555, 0.5, 0.5)
	getRecentDiffStatus := userinfo.Content.RecentScore[0].Difficulty
	var replyRecentDiffStatus string
	var replyRecentDIffStatusColor color.Color
	switch {
	case getRecentDiffStatus == 0:
		replyRecentDiffStatus = "Past"
		replyRecentDIffStatusColor = color.NRGBA{R: uint8(20), G: uint8(165), B: uint8(215), A: uint8(255)}
	case getRecentDiffStatus == 1:
		replyRecentDiffStatus = "Present"
		replyRecentDIffStatusColor = color.NRGBA{R: uint8(120), G: uint8(155), B: uint8(80), A: uint8(255)}
	case getRecentDiffStatus == 2:
		replyRecentDiffStatus = "Future"
		replyRecentDIffStatusColor = color.NRGBA{R: uint8(115), G: uint8(35), B: uint8(100), A: uint8(255)}
	case getRecentDiffStatus == 3:
		replyRecentDiffStatus = "Beyond"
		replyRecentDIffStatusColor = color.NRGBA{R: uint8(166), G: uint8(20), B: uint8(49), A: uint8(255)}
	}
	mainBG.SetFontFace(kazeRegularFaceSl)
	mainBG.SetColor(replyRecentDIffStatusColor)
	mainBG.DrawStringAnchored(fmt.Sprintf("%s | %s", replyRecentDiffStatus, strconv.FormatFloat(float64(userinfo.Content.Songinfo[0].Rating)/10, 'f', 1, 64)), 435, 590, 0.5, 0.5)
	var getClearType map[int]string
	var getClearTypeNum int
	var getTypeLocation string
	getClearType = map[int]string{
		0: "TL",
		1: "NC",
		2: "FR",
		3: "PM",
		4: "EC",
		5: "HC",
	}
	getClearTypeNum = userinfo.Content.RecentScore[0].ClearType
	switch {
	case getClearTypeNum == 1 || getClearTypeNum == 4 || getClearTypeNum == 5:
		getTypeLocation = arcaeaRes + "/resource/recent/clear_normal.png"
	case getClearTypeNum == 2:
		getTypeLocation = arcaeaRes + "/resource/recent/clear_full.png"
	case getClearTypeNum == 3:
		getTypeLocation = arcaeaRes + "/resource/recent/clear_pure.png"
	case getClearTypeNum == 0:
		getTypeLocation = arcaeaRes + "/resource/recent/clear_fail.png"
	}
	typeFileByte, err := os.ReadFile(getTypeLocation)
	if err != nil {
		panic(err)
	}
	typeFileByteImage, _, _ := image.Decode(bytes.NewReader(typeFileByte))
	typeFileByteImageContext := gg.NewContextForImage(typeFileByteImage)
	typeCovert2Image := imgfactory.Size(typeFileByteImage, 400, 400/typeFileByteImageContext.W()/typeFileByteImageContext.H()).Image()
	mainBG.DrawImage(typeCovert2Image, 240, 635)
	mainBG.SetFontFace(exoMidFaces)
	mainBG.SetColor(color.Black)
	mainBG.DrawStringAnchored(fmt.Sprintf("%s [%s]", FormatNumber(userinfo.Content.RecentScore[0].Score), getClearType[getClearTypeNum]), 440, 700, 0.5, 0.5)
	// draw rest of them, fuck.
	mainBG.SetFontFace(exoSmallFaceSS)
	mainBG.DrawStringAnchored(fmt.Sprintf("Play PTT:  %s", strconv.FormatFloat(userinfo.Content.RecentScore[0].Rating, 'f', 1, 64)), 220, 800, 0, 0)
	mainBG.DrawStringAnchored(fmt.Sprintf("Play Time: %s", FormatTimeStamp(userinfo.Content.RecentScore[0].TimePlayed)), 220, 830, 0, 0)
	// status.
	mainBG.SetFontFace(exoSemiBoldFace)
	mainBG.DrawStringAnchored(fmt.Sprintf("Pure:  %d(+%d)", userinfo.Content.RecentScore[0].PerfectCount, userinfo.Content.RecentScore[0].ShinyPerfectCount), 490, 780, 0, 0)
	mainBG.DrawStringAnchored(fmt.Sprintf("Far:  %d", userinfo.Content.RecentScore[0].NearCount), 490, 810, 0, 0)
	mainBG.DrawStringAnchored(fmt.Sprintf("Lost:  %d", userinfo.Content.RecentScore[0].MissCount), 490, 840, 0, 0)
	mainBG.Fill()
	return mainBG.Image()
}
