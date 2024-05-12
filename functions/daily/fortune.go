package daily

import (
	"encoding/json"
	"fmt"
	"github.com/FloatTech/zbputils/ctxext"
	Stringbreaker "github.com/MoYoez/Lucy-QOnebot/box/break"
	"github.com/MoYoez/Lucy-QOnebot/box/draw"
	"github.com/MoYoez/Lucy-QOnebot/box/emoji"
	"image"
	"image/color"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"math/rand"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"

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
	picDir, err := os.ReadDir(engine.DataFolder() + "randpic")
	if err != nil {
		panic(err)
	}
	picDirNum := len(picDir)
	reg := regexp.MustCompile(`[^.]+`)
	loadNotoSans := engine.DataFolder() + "NotoSansCJKsc-Regular.otf"
	data, err := os.ReadFile(engine.DataFolder() + "tarots.json")
	err = json.Unmarshal(data, &cardMap)

	engine.OnFullMatch("今日人品").SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			userS := strconv.FormatInt(ctx.Event.UserID, 10)
			now := time.Now().Format("20060102")
			userPic := strconv.FormatInt(ctx.Event.UserID, 10) + now + ".png"
			var avatarFormat *imgfactory.Factory
			var si = now + userS
			if signTF[si] == 0 {
				// use go func
				var avatarWaiter sync.WaitGroup
				avatarWaiter.Add(1)
				usersRandPic := fcext.RandSenderPerDayN(ctx.Event.UserID, picDirNum)
				picDirName := picDir[usersRandPic].Name()
				list := reg.FindAllString(picDirName, -1)
				// query nickname
				getNickName := ctx.CardOrNickName(ctx.Event.UserID)
				// remove emoji
				getNickName = emoji.EmojiRemover(getNickName)
				// get Callback
				p := rand.Intn(2)
				is := rand.Intn(77)
				i := is + 1
				cards := cardMap[(strconv.Itoa(i))]
				if p == 0 {
					info = cards.Info.Description
				} else {
					info = cards.Info.ReverseDescription
				}
				// modify this possibility to 40-100, don't be to low.
				randEveryone := fcext.RandSenderPerDayN(ctx.Event.UserID, 70)
				// add 30 and make it not so low.
				result[ctx.Event.UserID] = randEveryone + 30
				getNickName = Stringbreaker.BreakWords(getNickName, 24)
				// get user avatar
				go func() {
					defer avatarWaiter.Done()
					avatarByte, err := http.Get("https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(ctx.Event.UserID, 10) + "&s=640")
					if err != nil {
						ctx.SendChain(message.Text("获取不到头像( "))
						return
					}
					avatarByteUni, _, _ := image.Decode(avatarByte.Body)
					avatarFormat = imgfactory.Size(avatarByteUni, 100, 100)
				}()
				var getBackGroundMainColorR, getBackGroundMainColorG, getBackGroundMainColorB, mainContextWidth, mainContextHight int
				var mainContext *gg.Context
				// background
				img, err := gg.LoadImage(engine.DataFolder() + "randpic" + "/" + list[0] + ".png")
				if err != nil {
					panic(err)
				}
				bgFormat := imgfactory.Limit(img, 1920, 1080)
				getBackGroundMainColorR, getBackGroundMainColorG, getBackGroundMainColorB = draw.GetAverageColorAndMakeAdjust(bgFormat)
				mainContext = gg.NewContext(bgFormat.Bounds().Dx(), bgFormat.Bounds().Dy())
				mainContextWidth = mainContext.Width()
				mainContextHight = mainContext.Height()
				mainContext.DrawImage(bgFormat, 0, 0)
				// draw Round rectangle
				mainContext.SetFontFace(draw.LoadFontFace(loadNotoSans, 50))
				if err != nil {
					ctx.SendChain(message.Text("Font load err, Please contact maintainer."))
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
				// avatar,setname,desc
				// draw third round rectangle
				mainContext.SetRGBA255(91, 57, 83, 255)
				mainContext.SetFontFace(draw.LoadFontFace(loadNotoSans, 25))
				nameLength, _ := mainContext.MeasureString(getNickName)
				var renderLength float64
				renderLength = nameLength + 160
				if nameLength+160 <= 450 {
					renderLength = 450
				}
				mainContext.DrawRoundedRectangle(50, float64(mainContextHight-175), renderLength, 250, 20)
				mainContext.Fill()
				// avatar
				mainContext.SetRGBA255(255, 255, 255, 255)
				mainContext.DrawString("User Info", 60, float64(mainContextHight-150)+10) // basic ui
				mainContext.SetRGBA255(155, 121, 147, 255)
				mainContext.DrawString(getNickName, 180, float64(mainContextHight-150)+50)
				mainContext.DrawString(fmt.Sprintf("今日人品值: %d", randEveryone+40), 180, float64(mainContextHight-150)+100)
				mainContext.Fill()
				// AOSP time and date
				setInlineColor := color.NRGBA{R: uint8(getBackGroundMainColorR), G: uint8(getBackGroundMainColorG), B: uint8(getBackGroundMainColorB), A: 255}
				formatTimeDate := time.Now().Format("2006 / 01 / 02")
				formatTimeCurrent := time.Now().Format("15 : 04 : 05")
				formatTimeWeek := time.Now().Weekday().String()
				mainContext.SetFontFace(draw.LoadFontFace(loadNotoSans, 35))
				setOutlineColor := color.White
				draw.FunctionDrawBorderString(mainContext, formatTimeCurrent, 5, float64(mainContextWidth-80), 50, 1, 0.5, setInlineColor, setOutlineColor)
				draw.FunctionDrawBorderString(mainContext, formatTimeDate, 5, float64(mainContextWidth-80), 100, 1, 0.5, setInlineColor, setOutlineColor)
				draw.FunctionDrawBorderString(mainContext, formatTimeWeek, 5, float64(mainContextWidth-80), 150, 1, 0.5, setInlineColor, setOutlineColor)
				mainContext.FillPreserve()
				if err != nil {
					return
				}
				mainContext.SetFontFace(draw.LoadFontFace(loadNotoSans, 140))
				draw.FunctionDrawBorderString(mainContext, "|", 5, float64(mainContextWidth-30), 65, 1, 0.5, setInlineColor, setOutlineColor)
				// throw tarot card
				mainContext.SetFontFace(draw.LoadFontFace(loadNotoSans, 20))
				mainContext.SetRGBA255(91, 57, 83, 255)
				mainContext.DrawRoundedRectangle(float64(mainContextWidth-300), float64(mainContextHight-350), 450, 300, 20)
				mainContext.Fill()
				mainContext.SetRGBA255(255, 255, 255, 255)
				mainContext.SetLineWidth(3)
				mainContext.DrawString("今日塔罗牌", float64(mainContextWidth-300)+10, float64(mainContextHight-350)+30)
				mainContext.SetRGBA255(155, 121, 147, 255)
				mainContext.DrawString(cards.Name, float64(mainContextWidth-300)+10, float64(mainContextHight-350)+60)
				mainContext.DrawString(fmt.Sprintf("- %s", position[p]), float64(mainContextWidth-300)+10, float64(mainContextHight-350)+280)
				placedList := draw.SplitChineseString(info, 44)
				for ist, v := range placedList {
					mainContext.DrawString(v, float64(mainContextWidth-300)+10, float64(mainContextHight-350)+90+float64(ist*30))
				}
				mainContext.SetFontFace(draw.LoadFontFace(loadNotoSans, 16))
				mainContext.SetRGBA255(186, 163, 157, 255)
				avatarWaiter.Wait()
				mainContext.DrawImage(avatarFormat.Circle(0).Image(), 60, int(float64(mainContextHight-150)+25))
				mainContext.DrawStringAnchored("Generated By Lucy (HafuuNano), Design By MoeMagicMango", float64(mainContextWidth-15), float64(mainContextHight-30), 1, 1)
				mainContext.Fill()
				_ = mainContext.SavePNG(engine.DataFolder() + "jrrp/" + userPic)
				ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + engine.DataFolder() + "jrrp/" + userPic))
				signTF[si] = 1
			} else {
				ctx.SendChain(message.Text("今天已经测试过了哦w"), message.Image("file:///"+file.BOTPATH+"/"+engine.DataFolder()+"jrrp/"+userPic))
			}
		})
}
