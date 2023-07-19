// Package funwork 简单的测人品
package funwork

import (
	"encoding/json"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/FloatTech/zbputils/ctxext"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"math/rand"

	fcext "github.com/FloatTech/floatbox/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type card struct {
	Name string `json:"name"`
	Info struct {
		Description        string `json:"description"`
		ReverseDescription string `json:"reverseDescription"`
		ImgURL             string `json:"imgUrl"`
	} `json:"info"`
}

type cardset = map[string]card

var (
	info     string
	cardMap  = make(cardset, 256)
	position = []string{"正位", "逆位"}
	result   map[int64]int
	signTF   map[string]int
)

func init() {
	signTF = make(map[string]int)
	result = make(map[int64]int)
	getTarot := fcext.DoOnceOnSuccess(
		func(ctx *zero.Ctx) bool { // 检查 塔罗牌文件是否存在
			data, err := os.ReadFile(engine.DataFolder() + "tarots.json")
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return false
			}
			err = json.Unmarshal(data, &cardMap)
			if err != nil {
				panic(err)
			}
			return true
		},
	)
	engine.OnFullMatch("今日人品", getTarot).SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			userPic := strconv.FormatInt(ctx.Event.UserID, 10) + time.Now().Format("20060102") + ".png"
			picDir, err := os.ReadDir(engine.DataFolder() + "randpic")
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			picDirNum := len(picDir)
			usersRandPic := fcext.RandSenderPerDayN(ctx.Event.UserID, picDirNum)
			picDirName := picDir[usersRandPic].Name()
			reg := regexp.MustCompile(`[^.]+`)
			list := reg.FindAllString(picDirName, -1)
			var mutex sync.RWMutex // 添加读写锁以保证稳定性
			mutex.Lock()
			p := rand.Intn(2)
			is := rand.Intn(77)
			i := is + 1
			card := cardMap[(strconv.Itoa(i))]
			if p == 0 {
				info = card.Info.Description
			} else {
				info = card.Info.ReverseDescription
			}
			user := ctx.Event.UserID
			userS := strconv.FormatInt(user, 10)
			now := time.Now().Format("20060102")
			// modify this possibility to 40-100, don't be to low.
			randEveryone := fcext.RandSenderPerDayN(ctx.Event.UserID, 50)
			var si = now + userS // use map to store.
			loadNotoSans := engine.DataFolder() + "NotoSansCJKsc-Regular.otf"
			if signTF[si] == 0 {
				result[user] = randEveryone + 50
				// background
				img, err := gg.LoadImage(engine.DataFolder() + "randpic" + "/" + list[0] + ".png")
				if err != nil {
					panic(err)
				}
				bgFormat := imgfactory.Limit(img, 1920, 1080)
				getBackGroundMainColorR, getBackGroundMainColorG, getBackGroundMainColorB := GetAverageColorAndMakeAdjust(bgFormat)
				mainContext := gg.NewContext(bgFormat.Bounds().Dx(), bgFormat.Bounds().Dy())
				mainContextWidth := mainContext.Width()
				mainContextHight := mainContext.Height()
				mainContext.DrawImage(bgFormat, 0, 0)
				// draw Round rectangle
				mainContext.SetFontFace(LoadFontFace(loadNotoSans, 50))
				if err != nil {
					ctx.SendChain(message.Text("Something wrong while rendering pic? font"))
					return
				}
				// shade mode || not bugs(
				mainContext.SetLineWidth(4)
				mainContext.SetRGBA255(255, 255, 255, 255)
				mainContext.DrawRoundedRectangle(0, float64(mainContextHight-150), float64(mainContextWidth), 150, 16)
				mainContext.Stroke()
				mainContext.SetRGBA255(255, 224, 216, 215)
				mainContext.DrawRoundedRectangle(0, float64(mainContextHight-150), float64(mainContextWidth), 150, 16)
				mainContext.Fill()
				// avatar,name,desc
				// draw third round rectangle
				mainContext.SetRGBA255(91, 57, 83, 255)
				mainContext.SetFontFace(LoadFontFace(loadNotoSans, 25))
				nameLength, _ := mainContext.MeasureString(ctx.CardOrNickName(ctx.Event.UserID))
				var renderLength float64
				renderLength = nameLength + 160
				if nameLength+160 <= 450 {
					renderLength = 450
				}
				mainContext.DrawRoundedRectangle(50, float64(mainContextHight-175), renderLength, 250, 20)
				mainContext.Fill()
				avatarByte, err := http.Get("https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(ctx.Event.UserID, 10) + "&s=640")
				if err != nil {
					ctx.SendChain(message.Text("Something wrong while rendering pic? avatar IO err."))
					return
				}
				// avatar
				avatarByteUni, _, _ := image.Decode(avatarByte.Body)
				avatarFormat := imgfactory.Size(avatarByteUni, 100, 100)
				mainContext.DrawImage(avatarFormat.Circle(0).Image(), 60, int(float64(mainContextHight-150)+25))
				mainContext.SetRGBA255(255, 255, 255, 255)
				mainContext.DrawString("User Info", 60, float64(mainContextHight-150)+10) // basic ui
				mainContext.SetRGBA255(155, 121, 147, 255)
				mainContext.DrawString(ctx.CardOrNickName(ctx.Event.UserID), 180, float64(mainContextHight-150)+50)
				mainContext.DrawString(fmt.Sprintf("今日人品值: %d", randEveryone+50), 180, float64(mainContextHight-150)+100)
				mainContext.Fill()
				// AOSP time and date
				mainContext.SetRGBA255(getBackGroundMainColorR, getBackGroundMainColorG, getBackGroundMainColorB, 220)
				mainContext.SetFontFace(LoadFontFace(loadNotoSans, 25))
				if err != nil {
					ctx.SendChain(message.Text("Something wrong while rendering pic?"))
					return
				}
				mainContext.Stroke()
				formatTimeDate := time.Now().Format("2006/01/02")
				formatTimeCurrent := time.Now().Format("15:04:05")
				//	formatTimeLength, _ := mainContext.MeasureString(formatTimeDate)
				formatTimeWeek := time.Now().Weekday().String()
				mainContext.SetFontFace(LoadFontFace(loadNotoSans, 30))
				mainContext.DrawStringAnchored(formatTimeCurrent, float64(mainContextWidth-50), 40, 1, 0.5)
				mainContext.DrawStringAnchored(formatTimeDate, float64(mainContextWidth-50), 90, 1, 0.5)
				mainContext.DrawStringAnchored(formatTimeWeek, float64(mainContextWidth-50), 140, 1, 0.5)
				mainContext.FillPreserve()
				mainContext.SetRGBA255(getBackGroundMainColorR, getBackGroundMainColorG, getBackGroundMainColorB, 220)
				mainContext.SetFontFace(LoadFontFace(loadNotoSans, 160))
				if err != nil {
					ctx.SendChain(message.Text("Something wrong while rendering pic?", err))
					return
				}
				mainContext.SetLineWidth(3)
				mainContext.DrawString("|", float64(mainContextWidth-40), 140)
				mainContext.Fill()
				// throw tarot card
				mainContext.SetFontFace(LoadFontFace(loadNotoSans, 20))
				if err != nil {
					ctx.SendChain(message.Text("Something wrong while rendering pic?"))
					return
				}
				mainContext.SetRGBA255(91, 57, 83, 255)
				mainContext.DrawRoundedRectangle(float64(mainContextWidth-300), float64(mainContextHight-350), 450, 300, 20)
				mainContext.Fill()
				mainContext.SetRGBA255(255, 255, 255, 255)
				mainContext.SetLineWidth(3)
				mainContext.DrawString("今日塔罗牌", float64(mainContextWidth-300)+10, float64(mainContextHight-350)+30)
				mainContext.SetRGBA255(155, 121, 147, 255)
				mainContext.DrawString(card.Name, float64(mainContextWidth-300)+10, float64(mainContextHight-350)+60)
				mainContext.DrawString(fmt.Sprintf("- %s", position[p]), float64(mainContextWidth-300)+10, float64(mainContextHight-350)+280)
				placedList := splitChineseString(info, 44)
				for i, v := range placedList {
					mainContext.DrawString(v, float64(mainContextWidth-300)+10, float64(mainContextHight-350)+90+float64(i*30))
				}
				// output
				mainContext.SetFontFace(LoadFontFace(loadNotoSans, 20))
				mainContext.SetRGBA255(186, 163, 157, 255)
				mainContext.DrawStringAnchored("Generated By Lucy (HiMoYo),Designed By MoeMagicMango", float64(mainContextWidth-20), float64(mainContextHight-30), 1, 1)
				mainContext.Fill()
				_ = mainContext.SavePNG(engine.DataFolder() + "jrrp/" + userPic)
				ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + engine.DataFolder() + "jrrp/" + userPic))
				signTF[si] = 1
			} else {
				ctx.SendChain(message.Text("今天已经测试过了哦w"), message.Image("file:///"+file.BOTPATH+"/"+engine.DataFolder()+"jrrp/"+userPic))
			}
		})
}
func splitChineseString(s string, length int) []string {
	result := make([]string, 0)
	runes := []rune(s)
	start := 0
	for i := 0; i < len(runes); i++ {
		size := utf8.RuneLen(runes[i])
		if start+size > length {
			result = append(result, string(runes[0:i]))
			runes = runes[i:]
			i, start = 0, 0
		}
		start += size
	}
	if len(runes) > 0 {
		result = append(result, string(runes))
	}
	return result
}

// LoadFontFace load font face once before running, to work it quickly and save memory.
func LoadFontFace(filePath string, size float64) font.Face {
	fontFile, _ := os.ReadFile(filePath)
	fontFileParse, _ := opentype.Parse(fontFile)
	fontFace, _ := opentype.NewFace(fontFileParse, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	return fontFace
}

// GetAverageColorAndMakeAdjust different from k-means algorithm,it uses origin plugin's algorithm.(Reduce the cost of averge color usage.)
func GetAverageColorAndMakeAdjust(image image.Image) (int, int, int) {
	var RList []int
	var GList []int
	var BList []int
	width, height := image.Bounds().Size().X, image.Bounds().Size().Y
	// use the center of the bg, to make it more quickly and save memory and usage.
	for x := width / 6; x < width/2; x++ {
		for y := height / 4; y < height/2; y++ {
			r, g, b, _ := image.At(x, y).RGBA()
			RList = append(RList, int(r>>8))
			GList = append(GList, int(g>>8))
			BList = append(BList, int(b>>8))
		}
	}
	RAverage := int(Average(RList))
	GAverage := int(Average(GList))
	BAverage := int(Average(BList))
	// get Average and make adjust
	if BAverage > 200 {
		BAverage = int(math.Round(float64(BAverage-BAverage/6))) - 15
	} else {
		BAverage = int(math.Round(float64(BAverage+BAverage/6))) + 15
	}
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
