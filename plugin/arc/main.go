// Package arc for arc render b30
package arc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
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
	engine = control.Register("arcaea", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
		PrivateDataFolder: "arcaea",
	})
)

func init() {
	mainBG, _ := os.ReadFile(arcaeaRes + "/resource/b30/B30.png")
	// arc b30 is still in test(
	engine.OnRegex(`^!test arc\s*(\d+)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		id := ctx.State["regex_matched"].([]string)[1]
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
		// get player info
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Ok, trying to render \""+r.Content.AccountInfo.Name+"\" data."))
		mainBGDecoded, _, _ := image.Decode(bytes.NewReader(mainBG))
		basicBG := DrawMainUserB30(mainBGDecoded, r)
		tureResult := FinishedFullB30(basicBG, r)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, tureResult, nil)
		if err != nil {
			panic(err)
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`!test\sarc\sbind\s(.*)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		getBindInfo := ctx.State["regex_matched"].([]string)[1]
		context := isAlphanumeric(getBindInfo)
		var userinfo user
		if context == false {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回数据非法！"))
			return
		}
		checkTheContextIsNum := isNumericOrAlphanumeric(getBindInfo)
		if checkTheContextIsNum {
			// I don't know why I do this? useless(
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("检查到传入数值为纯数字，请选择\n1. Arcaea用户名 (e.g:MoeMagicMango)\n2. ArcaeaID (e.g:594698109)"))
			nextstep := ctx.FutureEvent("message", ctx.CheckSession())
			recv, _ := nextstep.Repeat()
			for i := range recv {
				texts := i.MessageString()
				switch {
				case texts == "1":
					data, err := aua.GetUserInfo(os.Getenv("aualink"), os.Getenv("auakey"), getBindInfo)
					dataBytes := helper.StringToBytes(data)
					err = json.Unmarshal(dataBytes, &userinfo)
					err = FormatInfo(ctx.Event.UserID, userinfo.Content.AccountInfo.Code, getBindInfo).BindUserArcaeaInfo(arcAcc)
					if err != nil {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误."))
						return
					}
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("User: `", getBindInfo, "` binded, id: ", userinfo.Content.AccountInfo.Code))
				case texts == "2":
					data, err := aua.GetUserInfo(os.Getenv("aualink"), os.Getenv("auakey"), getBindInfo)
					dataBytes := helper.StringToBytes(data)
					err = json.Unmarshal(dataBytes, &userinfo)
					err = FormatInfo(ctx.Event.UserID, getBindInfo, userinfo.Content.AccountInfo.Name).BindUserArcaeaInfo(arcAcc)
					if err != nil {
						ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误."))
						return
					}
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("User: `", getBindInfo, "` binded, id: ", userinfo.Content.AccountInfo.Code))
				default:
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回非法！"))
				}
			}
		} else {
			data, err := aua.GetUserInfo(os.Getenv("aualink"), os.Getenv("auakey"), getBindInfo)
			dataBytes := helper.StringToBytes(data)
			err = json.Unmarshal(dataBytes, &userinfo)
			err = FormatInfo(ctx.Event.UserID, userinfo.Content.AccountInfo.Code, getBindInfo).BindUserArcaeaInfo(arcAcc)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误."))
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("User: `", getBindInfo, "` binded, id: ", userinfo.Content.AccountInfo.Code))
		}
	})

	engine.OnFullMatch("!test arc b30").SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		id, err := GetUserArcaeaInfo(arcAcc, ctx)
		if err != nil || id == "" {
			ctx.SendChain(message.Text("cannot get user bind info."))
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
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, tureResult, nil)
		if err != nil {
			panic(err)
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`!test\sarc\schart\s([^\]]+)\s+([^\]]+)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		songName := ctx.State["regex_matched"].([]string)[1]
		songDiff := ctx.State["regex_matched"].([]string)[2]
		if songDiff == "" {
			songDiff = "ftr"
		}
		resultPreview, err := aua.GetSongPreview(os.Getenv("aualink"), os.Getenv("auakey"), songName, songDiff)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Unknown ERR:", err))
			return
		}
		if resultPreview == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Reply sent, but cannot find ", songName, " ("))
			return
		}
		var buf bytes.Buffer
		resultEncodingToImage, _, _ := image.Decode(bytes.NewReader(helper.StringToBytes(resultPreview)))
		err = jpeg.Encode(&buf, resultEncodingToImage, nil)
		if err != nil {
			panic(err)
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})
	// TODO: arcaea single song check (
}
