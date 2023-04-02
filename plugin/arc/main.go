// Package arc for arc render b30
package arc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	aua "github.com/MoYoez/Go-ArcaeaUnlimitedAPI"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"image"
	"image/jpeg"
	"os"
)

var (
	r      arcaea
	buf    bytes.Buffer
	engine = control.Register("arcaea", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
		PrivateDataFolder: "arcaea",
	})
)

func init() {
	mainBG, _ := os.ReadFile(arcaeaRes + "/resource/b30/B30.png")
	// arc b30 is still in test(
	engine.OnRegex(`^!test arc\s*(\d+)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		id := ctx.State["regex_matched"].([]string)[1]
		// get player info
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Ok, trying to get "+id+" data."))
		playerdata, err := aua.Best30(os.Getenv("aualink"), os.Getenv("auakey"), id)
		if err != nil {
			ctx.SendChain(message.Text("cannot get user data."))
			return
		}
		playerdataByte := helper.StringToBytes(playerdata)
		_ = json.Unmarshal(playerdataByte, &r)
		checkStatus := r.Status
		if checkStatus != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Status code is not valid."))
			return
		}
		mainBGDecoded, _, _ := image.Decode(bytes.NewReader(mainBG))
		basicBG := DrawMainUserB30(mainBGDecoded, r)
		tureResult := FinishedFullB30(basicBG, r)

		err = jpeg.Encode(&buf, tureResult, nil)
		if err != nil {
			panic(err)
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`!test\sarc\sbind\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getBindInfo := ctx.State["regex_matched"].([]string)[1]
		context := isAlphanumeric(getBindInfo)
		var userinfo user
		if context == false {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回数据非法！"))
			return
		}
		checkTheContextIsNum := isNumericOrAlphanumeric(getBindInfo)
		if checkTheContextIsNum == true {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("检查到传入数值为纯数字，请选择\n1. Arcaea用户名(e.g:MoeMagicMango)\n2. ArcaeaID(e.g:594698109)"))
			nextstep := ctx.FutureEvent("message", ctx.CheckSession())
			recv, _ := nextstep.Repeat()
			for i := range recv {
				texts := i.MessageString()
				switch {
				case texts == "1":
					data, err := aua.GetUserInfo(os.Getenv("aualink"), os.Getenv("auakey"), getBindInfo)
					dataBytes := helper.StringToBytes(data)
					err = json.Unmarshal(dataBytes, &userinfo)
					err = FormatInfo(ctx.Event.UserID, userinfo.Content.AccountInfo.Code).BindUserArcaeaInfo(arcAcc)
					if err != nil {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误."))
						return
					}
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("User: `", getBindInfo, "` binded. "))
				case texts == "2":
					err := FormatInfo(ctx.Event.UserID, userinfo.Content.AccountInfo.Code).BindUserArcaeaInfo(arcAcc)
					if err != nil {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误."))
						return
					}
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("User: `", getBindInfo, "` binded. "))
				default:
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回非法！"))
					return
				}
			}
		} else {
			data, err := aua.GetUserInfo(os.Getenv("aualink"), os.Getenv("auakey"), getBindInfo)
			dataBytes := helper.StringToBytes(data)
			err = json.Unmarshal(dataBytes, &userinfo)
			err = FormatInfo(ctx.Event.UserID, userinfo.Content.AccountInfo.Code).BindUserArcaeaInfo(arcAcc)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误."))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("User: `", getBindInfo, "` binded. "))
		}
	})

	engine.OnFullMatch("!test arc b30").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		id, err := GetUserInfo(arcAcc, ctx)
		if err != nil || id == "" {
			ctx.SendChain(message.Text("cannot get user info."))
			return
		}
		// get player info
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Ok, trying to get "+id+" data."))
		playerdata, err := aua.Best30(os.Getenv("aualink"), os.Getenv("auakey"), id)
		if err != nil {
			ctx.SendChain(message.Text("cannot get user data."))
			return
		}
		playerdataByte := helper.StringToBytes(playerdata)
		_ = json.Unmarshal(playerdataByte, &r)
		checkStatus := r.Status
		if checkStatus != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Status code is not valid."))
			return
		}
		mainBGDecoded, _, _ := image.Decode(bytes.NewReader(mainBG))
		basicBG := DrawMainUserB30(mainBGDecoded, r)
		tureResult := FinishedFullB30(basicBG, r)
		err = jpeg.Encode(&buf, tureResult, nil)
		if err != nil {
			panic(err)
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})
}
