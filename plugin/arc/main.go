// Package arc for arc render b30
package arc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/url"
	"os"
	"strconv"

	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	aua "github.com/MoYoez/Arcaea_auaAPI"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	userinfo user
	//	recordinfo record
	r      arcaea
	engine = control.Register("arcaea", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
		PrivateDataFolder: "arcaea",
	})
)

func init() {
	mainBG, _ := os.ReadFile(arcaeaRes + "/resource/b30/B30.png")
	// arc b30 still in test(
	engine.OnRegex(`^[！! /](a|arc)\s*(\d+)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		id := ctx.State["regex_matched"].([]string)[1]
		sessionKey, sessionKeyInfo := aua.GetSessionQuery(os.Getenv("aualink"), os.Getenv("auakey"), id)
		playerdataByte, _ := aua.DrawRequestArc(os.Getenv("aualink")+"/arcapi/user/bests/result?session_info="+sessionKey+"&overflow=10&with_recent=false&with_song_info=true", os.Getenv("auakey"))
		getPlayerReplyStatusId := gjson.Get(string(playerdataByte), "status").Int()
		switch {
		case getPlayerReplyStatusId == -31:
			getChartNumber := gjson.Get(string(playerdataByte), "content.queried_charts").String()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(m[-31]+getChartNumber+"\n预计等待时间：1-5 分钟"))
			return
		case getPlayerReplyStatusId == -32:
			getUserSessionWaitList := gjson.Get(string(playerdataByte), "content.current_account").Int()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(m[-32]+strconv.FormatInt(getUserSessionWaitList, 10)+"\n预计等待时间："+PerdictUserWaitTime(getUserSessionWaitList)))
			return
		case getPlayerReplyStatusId != 0 && getPlayerReplyStatusId != -33:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("w？貌似出现了一些问题：Code: ", getPlayerReplyStatusId, "信息：", m[int(getPlayerReplyStatusId)]))
			return
		}
		_ = json.Unmarshal(playerdataByte, &r)
		// get player info
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Get b30 data ~ trying to render \""+r.Content.AccountInfo.Name+"\" data."))
		mainBGDecoded, _, _ := image.Decode(bytes.NewReader(mainBG))
		basicBG := DrawMainUserB30(mainBGDecoded, r)
		tureResult := FinishedFullB30(basicBG, r)
		var buf bytes.Buffer
		err := png.Encode(&buf, tureResult)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("SessionKeyInfo: ", sessionKeyInfo), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`^[！! /](a|arc)\sbind\s(.*)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		getBindInfo := ctx.State["regex_matched"].([]string)[2]
		context := IsAlphanumeric(getBindInfo)

		if !context {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回数据非法！"))
			return
		}
		dataBytes, err := aua.GetUserInfo(os.Getenv("aualink"), os.Getenv("auakey"), getBindInfo)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误.", err))
			return
		}
		err = json.Unmarshal(dataBytes, &userinfo)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("数据出现问题，", err))
			return
		}
		checkStatus := userinfo.Status
		if checkStatus != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("数据返回异常: ", m[checkStatus]))
			return
		}
		err = FormatInfo(ctx.Event.UserID, userinfo.Content.AccountInfo.Code, userinfo.Content.AccountInfo.Name).BindUserArcaeaInfo(arcAcc)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("未知错误."))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("User: `", userinfo.Content.AccountInfo.Name, "` binded, id: ", userinfo.Content.AccountInfo.Code))
	})

	engine.OnRegex(`^[！! /](a|arc)\sb30$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		id, err := GetUserArcaeaInfo(arcAcc, ctx)
		if err != nil || id == "" {
			ctx.SendChain(message.Text("找不到用户信息，请检查你是否已经在Lucy端进行绑定，方式： “！arc bind {username | userid} ” "))
			return
		}
		sessionKey, sessionKeyInfo := aua.GetSessionQuery(os.Getenv("aualink"), os.Getenv("auakey"), id)
		playerdataByte, _ := aua.DrawRequestArc(os.Getenv("aualink")+"/arcapi/user/bests/result?session_info="+sessionKey+"&overflow=10&with_recent=false&with_song_info=true", os.Getenv("auakey"))
		getPlayerReplyStatusId := gjson.Get(string(playerdataByte), "status").Int()
		switch {
		case getPlayerReplyStatusId == -31:
			getChartNumber := gjson.Get(string(playerdataByte), "content.queried_charts").Int()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(m[-31]+strconv.FormatInt(getChartNumber, 10)+"\n预计等待时间：1-4.5 分钟"))
			return
		case getPlayerReplyStatusId == -32:
			getUserSessionWaitList := gjson.Get(string(playerdataByte), "content.current_account").Int()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(m[-32]+strconv.FormatInt(getUserSessionWaitList, 10)+"\n预计等待时间："+PerdictUserWaitTime(getUserSessionWaitList)))
			return
		case getPlayerReplyStatusId != 0 && getPlayerReplyStatusId != -33:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("w？貌似出现了一些问题：Code: ", getPlayerReplyStatusId, "信息：", m[int(getPlayerReplyStatusId)]))
			return
		}
		if getPlayerReplyStatusId == -33 {
			// check the file
			generatedName := r.Content.AccountInfo.Name + FormatRawTimeStamp(int64(r.Content.QueryTime)) + ".png"
			if file.IsExist(engine.DataFolder() + "save/" + generatedName) {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(m[-33]), message.Image("file:///"+file.BOTPATH+"/"+engine.DataFolder()+"save/"+r.Content.AccountInfo.Name+FormatRawTimeStamp(int64(r.Content.QueryTime))+".png"))
				return
			}
		}
		_ = json.Unmarshal(playerdataByte, &r)
		mainBGDecoded, _, _ := image.Decode(bytes.NewReader(mainBG))
		basicBG := DrawMainUserB30(mainBGDecoded, r)
		tureResult := FinishedFullB30(basicBG, r)
		getRawData, err := os.Create(engine.DataFolder() + "save/" + r.Content.AccountInfo.Name + FormatRawTimeStamp(int64(r.Content.QueryTime)) + ".png")
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似图片保存失败了（"))
			return
		}
		err = png.Encode(getRawData, tureResult)
		// get link

		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		var SessionKeyInfoFull string
		if sessionKeyInfo != "" {
			SessionKeyInfoFull = m[-33]
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(SessionKeyInfoFull), message.Image("file:///"+file.BOTPATH+"/"+engine.DataFolder()+"save/"+r.Content.AccountInfo.Name+FormatRawTimeStamp(int64(r.Content.QueryTime))+".png"))
	})

	engine.OnRegex(`^[！! /](a|arc)\spreview\s(\w+)(?:\s(\w+))?$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		songName := ctx.State["regex_matched"].([]string)[2]
		songDiff := ctx.State["regex_matched"].([]string)[3]
		// default is ftr
		if songDiff == "" {
			songDiff = "ftr"
		}
		// render , Chinese Words may cause unsafe when requesting,
		EscapeNameCN := url.QueryEscape(songName)
		replyImage, err := aua.GetSongPreview(os.Getenv("aualink"), os.Getenv("auakey"), EscapeNameCN, songDiff)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(err))
			return
		}
		var buf bytes.Buffer
		err = png.Encode(&buf, replyImage)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		} else {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经拿到图片数据~"))
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`^[！! /](a|arc)$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		// get info.
		id, err := GetUserArcaeaInfo(arcAcc, ctx)
		if err != nil || id == "" {
			ctx.SendChain(message.Text("找不到用户信息，请检查你是否已经在Lucy端进行绑定，方式： “！arc bind {username | userid} ” "))
			return
		}
		// get player info
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("好哦~正在帮你抓取最近游玩记录"))
		playerdataByte, err := aua.GetUserInfo(os.Getenv("aualink"), os.Getenv("auakey"), id)
		if err != nil {
			ctx.SendChain(message.Text("获取用户信息失败"))
			return
		}
		_ = json.Unmarshal(playerdataByte, &userinfo)
		checkStatus := userinfo.Status
		if checkStatus != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: \n", m[userinfo.Status]))
			return
		}
		replyImage := RenderUserRecentLog(userinfo)
		var buf bytes.Buffer
		err = png.Encode(&buf, replyImage)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})

	engine.OnRegex(`^[！! /](a|arc)\sinfo\s(\w+)(?:\s(\w+))?$`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		songName := ctx.State["regex_matched"].([]string)[2]
		songDiff := ctx.State["regex_matched"].([]string)[3]
		id, err := GetUserArcaeaInfo(arcAcc, ctx)
		if err != nil || id == "" {
			ctx.SendChain(message.Text("找不到用户信息，请检查你是否已经在Lucy端进行绑定，方式： “！arc bind {username | userid} ” "))
			return
		}
		if songDiff == "" {
			songDiff = "ftr"
		}
		getData, err := aua.GetUserBest(os.Getenv("aualink"), os.Getenv("auakey"), id, url.QueryEscape(songName), songDiff)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("发生错误：", err))
			return
		}
		var record recordBest
		_ = json.Unmarshal(getData, &record)
		checkStatus := record.Status
		if checkStatus != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("发生错误: ", m[checkStatus]))
			return
		}
		replyImage := RenderUserBestInfo(record)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, replyImage, nil)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("发生错误: ", err))
			return
		}
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})

	engine.OnFullMatch("!arc example render", zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getjson := engine.DataFolder() + "example.json"
		getdata, _ := os.ReadFile(getjson)
		_ = json.Unmarshal(getdata, &r)
		mainBGDecoded, _, _ := image.Decode(bytes.NewReader(mainBG))
		basicBG := DrawMainUserB30(mainBGDecoded, r)
		tureResult := FinishedFullB30(basicBG, r)
		var buf bytes.Buffer
		_ = png.Encode(&buf, tureResult)
		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64Str))
	})
}
