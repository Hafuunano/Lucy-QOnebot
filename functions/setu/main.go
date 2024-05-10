package setu

import (
	"encoding/base64"
	"github.com/FloatTech/AnimeAPI/bilibili"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
	"time"
)

var (
	ua            = web.RandUA()
	limitForPhoto = rate.NewManager[int64](time.Minute*3, 9)
	engine        = control.Register("setu", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "setu - provide some photos here.",
		PrivateDataFolder: "setu",
	})
)

func init() {
	engine.OnFullMatch("来份二次元").SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		if !limitForPhoto.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("要不休息一会再试试吧/"))
			return
		}
		realLink, _ := bilibili.GetRealURL("https://img.moehu.org/pic.php?id=img1")
		data, err := web.RequestDataWith(web.NewDefaultClient(), realLink, "GET", "https://sina.com", ua, nil)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		messageID := ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64.StdEncoding.EncodeToString(data)))
		time.Sleep(time.Second * 20)
		ctx.DeleteMessage(messageID)
	})

	engine.OnFullMatch("来份星空").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !limitForPhoto.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("要不休息一会再试试吧/"))
			return
		}
		realLink, _ := bilibili.GetRealURL("https://img.moehu.org/pic.php?id=xingk")
		data, err := web.RequestDataWith(web.NewDefaultClient(), realLink, "GET", "https://sina.com", ua, nil)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64.StdEncoding.EncodeToString(data)))
	})

	engine.OnFullMatch("来份兽耳").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !limitForPhoto.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("要不休息一会再试试吧/"))
			return
		}
		realLink, _ := bilibili.GetRealURL("https://img.moehu.org/pic.php?id=kemonomimi")
		data, err := web.RequestDataWith(web.NewDefaultClient(), realLink, "GET", "https://sina.com", ua, nil)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		messageID := ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64.StdEncoding.EncodeToString(data)))
		time.Sleep(time.Second * 20)
		ctx.DeleteMessage(messageID)
	})

	engine.OnFullMatch("来份白毛").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !limitForPhoto.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("要不休息一会再试试吧/"))
			return
		}
		realLink, _ := bilibili.GetRealURL("https://img.moehu.org/pic.php?id=yin")
		data, err := web.RequestDataWith(web.NewDefaultClient(), realLink, "GET", "https://sina.com", ua, nil)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		messageID := ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+base64.StdEncoding.EncodeToString(data)))
		time.Sleep(time.Second * 20)
		ctx.DeleteMessage(messageID)
	})
	engine.OnFullMatch("来份猫猫表情包").SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		if !limitForPhoto.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("要不休息一会再试试吧/"))
			return
		}
		data, err := web.GetData("https://img.moehu.org/pic.php?id=miao&return=json")
		if err != nil {
			ctx.SendChain(message.Text("ERR:", err))
			return
		}
		picURL := gjson.Get(string(data), "acgurl").String()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image(picURL))
	})
	engine.OnFullMatch("来份兽耳酱表情包").SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		if !limitForPhoto.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("要不休息一会再试试吧/"))
			return
		}
		data, err := web.GetData("https://img.moehu.org/pic.php?id=kemomimi&return=json")
		if err != nil {
			ctx.SendChain(message.Text("ERR:", err))
			return
		}
		picURL := gjson.Get(string(data), "acgurl").String()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image(picURL))
	})
}
