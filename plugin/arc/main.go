// Package arc for arc render b30
package arc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	aua "github.com/MoYoez/Arcaea_auaAPI"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"image"
	"image/jpeg"
	"os"
)

var (
	userinfo user
	r        arcaea
	engine   = control.Register("arcaea", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
		PrivateDataFolder: "arcaea",
	})
)

func init() {
	mainBG, _ := os.ReadFile(arcaeaRes + "/resource/b30/B30.png")
	// arc b30 is still in test(
	engine.OnRegex(`^[！!]arc\s*(\d+)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		id := ctx.State["regex_matched"].([]string)[1]
		sessionKey, sessionKeyInfo := aua.GetSessionQuery(os.Getenv("aualink"), os.Getenv("auakey"), id)
		playerdataByte, playerDataByteReturnMsg := aua.GetB30BySession(os.Getenv("aualink"), os.Getenv("auakey"), sessionKey)
		if playerDataByteReturnMsg != "" {
			ctx.SendChain(message.Text("SessionQuery: ", playerDataByteReturnMsg, "\nSession查询列队中，请过一段时间重新尝试呢～"))
			return
		}
		_ = json.Unmarshal(playerdataByte, &r)
		// get player info
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Get b30 data ~ trying to render \""+r.Content.AccountInfo.Name+"\" data."))
		mainBGDecoded, _, _ := image.Decode(bytes.NewReader(mainBG))
		basicBG := DrawMainUserB30(mainBGDecoded, r)
		tureResult := FinishedFullB30(basicBG, r)
		var buf bytes.Buffer
		err := jpeg.Encode(&buf, tureResult, nil)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("SessionKeyInfo: ", sessionKeyInfo), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`[！!]arc\sbind\s(.*)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		getBindInfo := ctx.State["regex_matched"].([]string)[1]
		context := IsAlphanumeric(getBindInfo)
		var userinfo user
		if !context {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回数据非法！"))
			return
		}
		dataBytes, err := aua.GetUserInfo(os.Getenv("aualink"), os.Getenv("auakey"), getBindInfo)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误.", err))
			return
		}
		_ = json.Unmarshal(dataBytes, &userinfo)
		checkStatus := userinfo.Status
		if checkStatus != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("数据返回异常，可能是接口出现问题: ERR: ", userinfo.Message))
			return
		}
		err = FormatInfo(ctx.Event.UserID, userinfo.Content.AccountInfo.Code, userinfo.Content.AccountInfo.Name).BindUserArcaeaInfo(arcAcc)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误."))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("User: `", userinfo.Content.AccountInfo.Name, "` binded, id: ", userinfo.Content.AccountInfo.Code))
	})

	engine.OnRegex(`[！!]arc\sb30$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		id, err := GetUserArcaeaInfo(arcAcc, ctx)
		if err != nil || id == "" {
			ctx.SendChain(message.Text("cannot get user bind info."))
			return
		}
		sessionKey, sessionKeyInfo := aua.GetSessionQuery(os.Getenv("aualink"), os.Getenv("auakey"), id)
		playerdataByte, playerDataByteReturnMsg := aua.GetB30BySession(os.Getenv("aualink"), os.Getenv("auakey"), sessionKey)
		if playerDataByteReturnMsg != "" {
			ctx.SendChain(message.Text("SessionQuery: ", playerDataByteReturnMsg, "\nSession查询列队中，请过一段时间重新尝试呢～"))
			return
		}
		_ = json.Unmarshal(playerdataByte, &r)
		mainBGDecoded, _, _ := image.Decode(bytes.NewReader(mainBG))
		basicBG := DrawMainUserB30(mainBGDecoded, r)
		tureResult := FinishedFullB30(basicBG, r)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, tureResult, nil)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		var SessionKeyInfoFull string
		if sessionKeyInfo != "" {
			SessionKeyInfoFull = "SessionKeyInfo: " + sessionKeyInfo
		}

		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(SessionKeyInfoFull), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`[！!]arc\schart\s([^\]]+)\s+([^\]] +)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		songName := ctx.State["regex_matched"].([]string)[1]
		songDiff := ctx.State["regex_matched"].([]string)[2]
		resultPreview, err := aua.GetSongPreview(os.Getenv("aualink"), os.Getenv("auakey"), songName, songDiff)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Reply sent, but cannot find ", songName, " (", err))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请等待一会哦~已经拿到图片请求了x"))
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, resultPreview, nil)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`[！!]arc$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		// get info.
		id, err := GetUserArcaeaInfo(arcAcc, ctx)
		if err != nil || id == "" {
			ctx.SendChain(message.Text("cannot get user bind info."))
			return
		}
		// get player info
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("好哦~正在帮你抓取最近游玩记录"))
		playerdataByte, err := aua.GetUserInfo(os.Getenv("aualink"), os.Getenv("auakey"), id)
		if err != nil {
			ctx.SendChain(message.Text("cannot get user data."))
			return
		}
		_ = json.Unmarshal(playerdataByte, &userinfo)
		checkStatus := userinfo.Status
		if checkStatus != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: \n", userinfo.Message))
			return
		}
		replyImage := RenderUserRecentLog(userinfo)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, replyImage, nil)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})
	engine.OnRegex(`[! !]arc\sbest\s([^\]]+)\s+([^\]] +)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		songName := ctx.State["regex_matched"].([]string)[1]
		songDiff := ctx.State["regex_matched"].([]string)[2]
		id, err := GetUserArcaeaInfo(arcAcc, ctx)
		if err != nil || id == "" {
			ctx.SendChain(message.Text("cannot get user bind info."))
			return
		}
		getData, err := aua.GetUserBest(os.Getenv("aualink"), os.Getenv("auakey"), id, songName, songDiff)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR:", err))
			return
		}
		_ = json.Unmarshal(getData, &userinfo)
		checkStatus := userinfo.Status
		if checkStatus != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: \n", userinfo.Message))
			return
		}
		replyImage := RenderUserRecentLog(userinfo)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, replyImage, nil)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})
}
