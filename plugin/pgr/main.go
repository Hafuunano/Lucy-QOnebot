package pgr // Package pgr hosted by Phigros-Library
import (
	"encoding/json"
	"image"
	"image/color"
	"net/http"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/FloatTech/floatbox/web"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/disintegration/imaging"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// too lazy,so this way is to use thrift host server (Working on HiMoYo Cloud.) (replace: now use PUA API)

// update: use PhigrosUnlimitedAPI + Phigros Library as Maintainer.

var (
	engine = control.Register("phigros", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n",
		PrivateDataFolder: "phi",
	})
)

func init() {
	engine.OnRegex(`^[! ！/]pgr\sbind\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		hash := ctx.State["regex_matched"].([]string)[1]
		userInfo := GetUserInfoTimeFromDatabase(ctx.Event.UserID)
		if userInfo+(12*60*60) > time.Now().Unix() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("12小时内仅允许绑定一次哦"))
			return
		}
		indexReply := DecHashToRaw(hash)
		// get session.
		if indexReply == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://pgr.impart.icu 获取绑定码进行绑定"))
			return
		}
		getQQID, getSessionID := RawJsonParse(indexReply)
		if getQQID != ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请求Hash中QQ号不一致，请使用自己的号重新申请"))
			return
		}
		if utf8.RuneCountInString(getSessionID) != 25 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Session 传入数值出现错误，请重新绑定"))
			return
		}
		_ = FormatUserDataBase(getQQID, getSessionID, time.Now().Unix()).BindUserDataBase()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定成功～"))
	})
	engine.OnRegex(`^[! ！/]pgr\sb19$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		data := GetUserInfoFromDatabase(ctx.Event.UserID)
		getDataSession := data.PhiSession
		if getDataSession == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://pgr.impart.icu 获取绑定码进行绑定 "))
			return
		}
		userData := GetUserInfoFromDatabase(ctx.Event.UserID)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("好哦~正在帮你请求，请稍等一下啦w~大约需要1-2分钟"))
		var dataWaiter sync.WaitGroup
		var AvatarWaiter sync.WaitGroup
		var getAvatarFormat *gg.Context
		var phidata []byte
		var setGlobalStat = true
		AvatarWaiter.Add(1)
		dataWaiter.Add(1)
		go func() {
			defer dataWaiter.Done()
			getFullLink := "https://pgrapi.impart.icu/api/phi/bests?session=" + userData.PhiSession + "&overflow=2"
			phidata, _ = web.GetData(getFullLink)
			if phidata == nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("目前 Unoffical Phigros API 暂时无法工作 请过一段时候尝试"))
				setGlobalStat = false
				return
			}
			err := json.Unmarshal(phidata, &phigrosB19)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("发生解析错误\n", err))
				return
			}
			if !phigrosB19.Status {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("w? 貌似出现了一些问题x"))
				return
			}
		}()
		go func() {
			defer AvatarWaiter.Done()
			avatarByte, err := http.Get("https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(ctx.Event.UserID, 10) + "&s=640")
			if err != nil {
				ctx.SendChain(message.Text("Something wrong while rendering pic? avatar IO err."))
				return
			}
			// avatar
			avatarByteUni, _, _ := image.Decode(avatarByte.Body)
			showUserAvatar := imaging.Resize(avatarByteUni, 250, 250, imaging.Lanczos)
			getAvatarFormat = gg.NewContext(250, 250)
			getAvatarFormat.DrawRoundedRectangle(0, 0, 248, 248, 20)
			getAvatarFormat.Clip()
			getAvatarFormat.DrawImage(showUserAvatar, 0, 0)
			getAvatarFormat.Fill()
		}()
		getRawBackground, _ := gg.LoadImage(backgroundRender)
		getMainBgRender := gg.NewContextForImage(imaging.Resize(getRawBackground, 2750, 5500, imaging.Lanczos))
		_ = getMainBgRender.LoadFontFace(font, 30)
		// header background
		drawTriAngle(getMainBgRender, a, 0, 166, 1324, 410)
		getMainBgRender.SetRGBA255(0, 0, 0, 160)
		getMainBgRender.Fill()
		drawTriAngle(getMainBgRender, a, 1318, 192, 1600, 350)
		getMainBgRender.SetRGBA255(0, 0, 0, 160)
		getMainBgRender.Fill()
		drawTriAngle(getMainBgRender, a, 1320, 164, 6, 414)
		getMainBgRender.SetColor(color.White)
		getMainBgRender.Fill()
		// header background end.
		// load icon with other userinfo.
		getMainBgRender.SetColor(color.White)
		logo, _ := gg.LoadPNG(icon)
		getImageLogo := imaging.Resize(logo, 290, 290, imaging.Lanczos)
		getMainBgRender.DrawImage(getImageLogo, 50, 216)
		fontface, _ := gg.LoadFontFace(font, 90)
		getMainBgRender.SetFontFace(fontface)
		getMainBgRender.DrawString("Phigros", 422, 336)
		getMainBgRender.DrawString("RankingScore查询", 422, 462)
		// draw userinfo path
		renderHeaderText, _ := gg.LoadFontFace(font, 54)
		getMainBgRender.SetFontFace(renderHeaderText)
		dataWaiter.Wait()
		if setGlobalStat == false {
			return
		}
		getMainBgRender.DrawString("Player: "+phigrosB19.Content.PlayerID, 1490, 300)
		getMainBgRender.DrawString("RankingScore: "+strconv.FormatFloat(phigrosB19.Content.RankingScore, 'f', 3, 64), 1490, 380)
		getMainBgRender.DrawString("ChanllengeMode: ", 1490, 460) // +56
		getColor, getLink := GetUserChallengeMode(phigrosB19.Content.ChallengeModeRank)
		if getColor != "" {
			getColorLink := ChanllengeMode + getColor + ".png"
			getColorImage, _ := gg.LoadImage(getColorLink)
			getMainBgRender.DrawImage(imaging.Resize(getColorImage, 238, 130, imaging.Lanczos), 1912, 390)
		}
		renderHeaderTextNumber, _ := gg.LoadFontFace(font, 65)
		getMainBgRender.SetFontFace(renderHeaderTextNumber)
		// white glow render
		getMainBgRender.SetRGB(1, 1, 1)
		getMainBgRender.DrawStringAnchored(getLink, 2021, 430, 0.4, 0.4)
		// avatar
		AvatarWaiter.Wait()
		getMainBgRender.DrawImage(getAvatarFormat.Image(), 2321, 230)
		getMainBgRender.Fill()
		// render
		CardRender(getMainBgRender, phidata)
		// draw bottom
		_ = getMainBgRender.LoadFontFace(font, 40)
		getMainBgRender.SetColor(color.White)
		getMainBgRender.Fill()
		getMainBgRender.DrawString("Generated By Lucy (HiMoYoBOT) | Designed By Eastown | Data From Phigros Library Project", 10, 5480)
		_ = getMainBgRender.SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("file:///"+file.BOTPATH+"/"+engine.DataFolder()+"save/"+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
	engine.OnRegex(`^[! ！/]pgr\sroll\s(\d+)`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		var wg sync.WaitGroup
		var avatarWaitGroup sync.WaitGroup
		var dataWaiter sync.WaitGroup
		var getMainBgRender *gg.Context
		var getAvatarFormat *gg.Context
		var setGlobalStat = false
		var phidata []byte
		wg.Add(1)
		avatarWaitGroup.Add(1)
		dataWaiter.Add(1)
		// get Session From Database.
		data := GetUserInfoFromDatabase(ctx.Event.UserID)
		getDataSession := data.PhiSession
		if getDataSession == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("由于Session特殊性，请前往 https://pgr.impart.icu 获取绑定码进行绑定"))
			return
		}

		//	GetSessionByPhigrosLibraryProject(getDataSession, ctx)
		// getPhigrosKey := os.Getenv("puakey")
		userData := GetUserInfoFromDatabase(ctx.Event.UserID)
		getRoll := ctx.State["regex_matched"].([]string)[1]
		getRollInt, err := strconv.ParseInt(getRoll, 10, 64)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请求roll不合法"))
			return
		}
		if getRollInt > 40 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("限制查询数为小于40"))
			return
		}
		getOverFlowNumber := getRollInt - 19
		if getOverFlowNumber <= 0 {
			getOverFlowNumber = 0
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("好哦~正在帮你请求，请稍等一下啦w~大约需要1-2分钟"))
		// data handling.
		go func() {
			defer dataWaiter.Done()
			getFullLink := "https://pgrapi.impart.icu/api/phi/bests?session=" + userData.PhiSession + "&overflow=" + strconv.Itoa(int(getOverFlowNumber))
			phidata, _ = web.GetData(getFullLink)
			if phidata == nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("目前 Unoffical Phigros Library 暂时无法工作 请过一段时候尝试"))
				setGlobalStat = false
				return
			}
			err = json.Unmarshal(phidata, &phigrosB19)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("发生解析错误: \n", err))
				return
			}
			if !phigrosB19.Status {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("w? 貌似出现了一些问题x\n", phigrosB19.Message))
				return
			}
		}()
		go func() {
			defer wg.Done()
			getRawBackground, _ := gg.LoadImage(backgroundRender)
			getMainBgRender = gg.NewContextForImage(imaging.Resize(getRawBackground, 2750, int(5250+getOverFlowNumber*200), imaging.Lanczos))
		}()
		go func() {
			defer avatarWaitGroup.Done()
			// draw Avatar, avatar from qq.
			avatarByte, err := http.Get("https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(ctx.Event.UserID, 10) + "&s=640")
			if err != nil {
				ctx.SendChain(message.Text("Something wrong while rendering pic? avatar IO err."))
				return
			}
			// avatar
			avatarByteUni, _, _ := image.Decode(avatarByte.Body)
			showUserAvatar := imaging.Resize(avatarByteUni, 250, 250, imaging.Lanczos)
			getAvatarFormat = gg.NewContext(250, 250)
			getAvatarFormat.DrawRoundedRectangle(0, 0, 248, 248, 20)
			getAvatarFormat.Clip()
			getAvatarFormat.DrawImage(showUserAvatar, 0, 0)
			getAvatarFormat.Fill()
		}()
		wg.Wait()
		_ = getMainBgRender.LoadFontFace(font, 30)
		// header background
		drawTriAngle(getMainBgRender, a, 0, 166, 1324, 410)
		getMainBgRender.SetRGBA255(0, 0, 0, 160)
		getMainBgRender.Fill()
		drawTriAngle(getMainBgRender, a, 1318, 192, 1600, 350)
		getMainBgRender.SetRGBA255(0, 0, 0, 160)
		getMainBgRender.Fill()
		drawTriAngle(getMainBgRender, a, 1320, 164, 6, 414)
		getMainBgRender.SetColor(color.White)
		getMainBgRender.Fill()
		// header background end.
		// load icon with other userinfo.
		getMainBgRender.SetColor(color.White)
		logo, _ := gg.LoadPNG(icon)
		getImageLogo := imaging.Resize(logo, 290, 290, imaging.Lanczos)
		getMainBgRender.DrawImage(getImageLogo, 50, 216)
		fontface, _ := gg.LoadFontFace(font, 90)
		getMainBgRender.SetFontFace(fontface)
		getMainBgRender.DrawString("Phigros", 422, 336)
		getMainBgRender.DrawString("RankingScore查询", 422, 462)
		dataWaiter.Wait()
		if setGlobalStat == false {
			return
		}
		// draw userinfo path
		renderHeaderText, _ := gg.LoadFontFace(font, 54)
		getMainBgRender.SetFontFace(renderHeaderText)
		// wait data until fine.
		getMainBgRender.DrawString("Player: "+phigrosB19.Content.PlayerID, 1490, 300)
		getMainBgRender.DrawString("RankingScore: "+strconv.FormatFloat(phigrosB19.Content.RankingScore, 'f', 3, 64), 1490, 380)
		getMainBgRender.DrawString("ChanllengeMode: ", 1490, 460) // +56
		getColor, getLink := GetUserChallengeMode(phigrosB19.Content.ChallengeModeRank)
		if getColor != "" {
			getColorLink := ChanllengeMode + getColor + ".png"
			getColorImage, _ := gg.LoadImage(getColorLink)
			getMainBgRender.DrawImage(imaging.Resize(getColorImage, 238, 130, imaging.Lanczos), 1912, 390)
		}
		renderHeaderTextNumber, _ := gg.LoadFontFace(font, 65)
		getMainBgRender.SetFontFace(renderHeaderTextNumber)
		// white glow render
		getMainBgRender.SetRGB(1, 1, 1)
		getMainBgRender.DrawStringAnchored(getLink, 2021, 430, 0.4, 0.4)
		avatarWaitGroup.Wait()
		getMainBgRender.DrawImage(getAvatarFormat.Image(), 2321, 230)
		getMainBgRender.Fill()
		// render
		CardRender(getMainBgRender, phidata)
		// draw bottom
		_ = getMainBgRender.LoadFontFace(font, 40)
		getMainBgRender.SetColor(color.White)
		getMainBgRender.Fill()
		getMainBgRender.DrawString("Generated By Lucy (HiMoYoBOT) | Designed By Eastown | Data From Phigros Library Project", 10, float64(5110+getOverFlowNumber*200+100))
		_ = getMainBgRender.SavePNG(engine.DataFolder() + "save/" + "roll" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("file:///"+file.BOTPATH+"/"+engine.DataFolder()+"save/"+"roll"+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
}
