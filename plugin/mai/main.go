package mai

import (
	"bytes"
	"encoding/json"
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"image"
	"image/color"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	engine = control.Register("maidx", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
		PrivateDataFolder: "maidx",
	})
)

func init() {
	engine.OnRegex(`^[！!]chun$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		dataPlayer, err := QueryChunDataFromQQ(int(uid))
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		txt := HandleChunDataByUsingText(dataPlayer)
		base64Font, err := text.RenderToBase64(txt, text.BoldFontFile, 1920, 45)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		ctx.SendChain(message.Image("base64://" + binary.BytesToString(base64Font)))
	})
	engine.OnRegex(`^[! ！/](mai|b50)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		dataPlayer, err := QueryMaiBotDataFromQQ(int(uid))
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		var data player
		_ = json.Unmarshal(dataPlayer, &data)
		renderImg := FullPageRender(data, ctx)
		_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image(Saved+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		matched := ctx.State["regex_matched"].([]string)[2]
		dataPlayer, err := QueryMaiBotDataFromUserName(matched)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		var data player
		_ = json.Unmarshal(dataPlayer, &data)
		renderImg := FullPageRender(data, ctx)
		_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image(Saved+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\splate\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPlateInfo := ctx.State["regex_matched"].([]string)[2]
		_ = FormatUserDataBase(ctx.Event.UserID, getPlateInfo).BindUserDataBase()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经将称号绑定上去了哦w"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\supload`, PictureHandler).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPic := ctx.State["image_url"].([]string)[0]
		imageData, err := web.GetData(getPic)
		if err != nil {
			return
		}
		getRaw, _, err := image.Decode(bytes.NewReader(imageData))
		if err != nil {
			panic(err)
		}
		// pic Handler
		getRenderPlatePicRaw := gg.NewContext(1260, 210)
		getRenderPlatePicRaw.DrawRoundedRectangle(0, 0, 1260, 210, 10)
		getRenderPlatePicRaw.Clip()
		getHeight := getRaw.Bounds().Dy()
		getLength := getRaw.Bounds().Dx()
		var getHeightHandler, getLengthHandler int
		switch {
		case getHeight < 210 && getLength < 1260:
			getRaw = Resize(getRaw, 1260, 210)
			getHeightHandler = 0
			getLengthHandler = 0
		case getHeight < 210:
			getRaw = Resize(getRaw, getLength, 210)
			getHeightHandler = 0
			getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
		case getLength < 1260:
			getRaw = Resize(getRaw, 1260, getHeight)
			getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
			getLengthHandler = 0
		default:
			getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
			getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
		}
		getRenderPlatePicRaw.DrawImage(getRaw, getLengthHandler, getHeightHandler)
		getRenderPlatePicRaw.Fill()
		getRenderPlatePicRaw.SavePNG(userPlate + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经存入了哦w"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sremove`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		os.Remove(userPlate + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经删掉了哦w"))
	})
	engine.OnFullMatch("/mai example render", zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		dataBytes, err := os.ReadFile(engine.DataFolder() + "example.json")
		if err != nil {
			return
		}
		var data player
		_ = json.Unmarshal(dataBytes, &data)
		// render a example pic
		var getAvatarFormat *gg.Context
		// avatarByte, err := http.Get("https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(ctx.Event.UserID, 10) + "&s=640")
		avatarByte, err := http.Get("https://cdn.sep.cc/avatar/22b242a28bb848f2629f2a636bba9c03?s=600")
		if err != nil {
			panic(err)
		}
		avatarByteUni, _, _ := image.Decode(avatarByte.Body)
		avatarFormat := imgfactory.Size(avatarByteUni, 180, 180)
		getAvatarFormat = gg.NewContext(180, 180)
		getAvatarFormat.DrawRoundedRectangle(0, 0, 178, 178, 20)
		getAvatarFormat.Clip()
		getAvatarFormat.DrawImage(avatarFormat.Image(), 0, 0)
		getAvatarFormat.Fill()
		// render Header.
		b50Render := gg.NewContext(2090, 1660)
		rawPlateData, err := gg.LoadImage(userPlate + "example.png")
		if err == nil {
			b50bg = b50Custom
			b50Render.DrawImage(rawPlateData, 595, 30)
			b50Render.Fill()
		}
		getContent, _ := gg.LoadImage(b50bg)
		b50Render.DrawImage(getContent, 0, 0)
		b50Render.Fill()
		// render user info
		b50Render.DrawImage(getAvatarFormat.Image(), 610, 50)
		b50Render.Fill()
		// render Userinfo
		b50Render.SetFontFace(nameTypeFont)
		b50Render.SetColor(color.Black)
		b50Render.DrawStringAnchored(strings.Join(strings.Split("StarKoi", ""), " "), 825, 160, 0, 0)
		b50Render.Fill()
		b50Render.SetFontFace(titleFont)
		b50Render.DrawStringAnchored(strings.Join(strings.Split("Lucy Kawaii ^^", ""), " "), 1050, 207, 0.5, 0.5)
		b50Render.Fill()
		getRating := getRatingBg(data.Rating)
		getRatingBG, err := gg.LoadImage(loadMaiPic + getRating)
		if err != nil {
			panic(err)
		}
		b50Render.DrawImage(getRatingBG, 800, 40)
		b50Render.Fill()
		// draw number
		b50Render.SetFontFace(scoreFont)
		b50Render.SetRGBA255(236, 219, 113, 255)
		b50Render.DrawStringAnchored(strconv.Itoa(data.Rating), 1056, 60, 1, 1)
		b50Render.Fill()
		// Render Card Type
		getSDLength := len(data.Charts.Sd)
		getDXLength := len(data.Charts.Dx)
		getDXinitX := 45
		getDXinitY := 1225
		getInitX := 45
		getInitY := 285
		var i int
		for i = 0; i < getSDLength; i++ {
			b50Render.DrawImage(RenderCard(data.Charts.Sd[i], i+1), getInitX, getInitY)
			getInitX += 400
			if getInitX == 2045 {
				getInitX = 45
				getInitY += 125
			}
		}
		for dx := 0; dx < getDXLength; dx++ {
			b50Render.DrawImage(RenderCard(data.Charts.Dx[dx], dx+1), getDXinitX, getDXinitY)
			getDXinitX += 400
			if getDXinitX == 2045 {
				getDXinitX = 45
				getDXinitY += 125
			}
		}
		b50Render.SavePNG(engine.DataFolder() + "save/" + "example.png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image(Saved+"example.png"))
	})
}
