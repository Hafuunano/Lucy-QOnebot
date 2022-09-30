// funwork Hi NekoPachi!
package funwork

import (
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

const (
	ua      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.0.0 Safari/537.36"
	Referer = "https://manual-lucy.himoyo.cn"
)

var (
	engine = control.Register("funwork", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://manual-lucy.himoyo.cn",
		PrivateDataFolder: "funwork",
	})
)

// WorkON: APIWORK
func init() {
	engine.OnFullMatch("一言").Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		info, err := web.GetData("https://v1.hitokoto.cn/")
		if err != nil {
			ctx.Send(message.Text("ERROR:", err))
			return
		}
		hitokoto := gjson.Get(helper.BytesToString(info), "hitokoto").String()
		hitokotoFrom := gjson.Get(helper.BytesToString(info), "from").String()
		hitokotoFromName := gjson.Get(helper.BytesToString(info), "from_who").String()
		if hitokotoFromName == "null" {
			hitokotoFromName = "未知"
		}
		ctx.SendChain(message.Text("!~Lucy帮你找到了这个www\n一言: ", hitokoto, "\n出处: ", hitokotoFrom, "\n作者: ", hitokotoFromName))
	})

	engine.OnFullMatch("动漫一言").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		data, err := web.RequestDataWith(web.NewDefaultClient(), "http://ovooa.com/API/dmyiyan/api.php", "GET", Referer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(helper.BytesToString(data)))
	})

	engine.OnFullMatch("来份网易云热评").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			data, err := web.RequestDataWith(web.NewDefaultClient(), "https://ovooa.com/API/wyrp/api.php?type=text", "GET", Referer, ua)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Text(helper.BytesToString(data)))
		})
	engine.OnFullMatch("舔狗日记").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		data, err := web.RequestDataWith(web.NewDefaultClient(), "https://ovooa.com/API/tgrj/api.php", "GET", Referer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(helper.BytesToString(data)))
	})

	engine.OnFullMatch("哄我").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		data, err := web.RequestDataWith(web.NewDefaultClient(), "https://ovooa.com/API/wryl/api.php", "GET", Referer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(helper.BytesToString(data)))
	})
	engine.OnFullMatch("答案之书").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		data, err := web.RequestDataWith(web.NewDefaultClient(), "http://ovooa.com/API/daan/api?type=text", "GET", Referer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(helper.BytesToString(data)))
	})
}
