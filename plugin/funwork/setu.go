// 使用国内源的iw233 + 新背景实现
package funwork

import (
	"math/rand"

	"time"

	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var (
	limit = rate.NewManager[int64](time.Minute*3, 8)
)

func init() {
	engine.OnFullMatch("来份二次元").SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.UserID).Acquire() {
			return
		}
		data, err := web.RequestDataWith(web.NewDefaultClient(), "https://mirlkoi.ifast3.vipnps.vip/api.php?sort=random&type=json", "GET", Referer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		picURL := gjson.Get(string(data), "pic.0").String()

		messageID := ctx.SendChain(message.Image(picURL))
		time.Sleep(time.Second * 20)
		ctx.DeleteMessage(messageID)
	})

	engine.OnFullMatch("来份星空").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.UserID).Acquire() {
			return
		}
		data, err := web.RequestDataWith(web.NewDefaultClient(), "https://mirlkoi.ifast3.vipnps.vip/api.php?sort=xing&type=json", "GET", Referer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		picURL := gjson.Get(string(data), "pic.0").String()
		ctx.SendChain(message.Image(picURL))
	})

	engine.OnFullMatch("来份兽耳").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.UserID).Acquire() {
			return
		}
		data, err := web.RequestDataWith(web.NewDefaultClient(), "https://mirlkoi.ifast3.vipnps.vip/api.php?sort=cat&type=json", "GET", Referer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		picURL := gjson.Get(string(data), "pic.0").String()
		messageID := ctx.SendChain(message.Image(picURL))
		time.Sleep(time.Second * 20)
		ctx.DeleteMessage(messageID)
	})

	engine.OnFullMatch("来份白毛").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.UserID).Acquire() {
			return
		}
		data, err := web.RequestDataWith(web.NewDefaultClient(), "https://mirlkoi.ifast3.vipnps.vip/api.php?sort=yin&type=json", "GET", Referer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		picURL := gjson.Get(string(data), "pic.0").String()
		messageID := ctx.SendChain(message.Image(picURL))
		time.Sleep(time.Second * 20)
		ctx.DeleteMessage(messageID)
	})

	engine.OnFullMatch("随机老婆").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.UserID).Acquire() {
			return
		}
		data, err := web.RequestDataWith(web.NewDefaultClient(), "http://ovooa.com/API/Lovely/api?type=text", "GET", Referer, ua)
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return
		}
		messageID := ctx.SendChain(message.Image(helper.BytesToString(data)))
		time.Sleep(time.Second * 20)
		ctx.DeleteMessage(messageID)
	})
	engine.OnFullMatch("涩涩", zero.OnlyToMe).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.UserID).Acquire() {
			return
		}
		if rand.Intn(4) == 1 {
			data, err := web.RequestDataWith(web.NewDefaultClient(), "http://iw233.fgimax2.fgnwctvip.com/API/Ghs.php?type=json", "GET", Referer, ua)
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			picURL := gjson.Get(string(data), "pic").String()
			messageID := ctx.SendChain(message.Text(picURL))
			time.Sleep(time.Second * 20)
			ctx.DeleteMessage(messageID)
		} else {
			ctx.Send(message.Text([]string{"看什么看！咱没有涩图 哼!", "只有笨蛋才看涩图", "咱想把你变成涩图", "好孩子是不会看涩图的", "敲~笨蛋 不许色色", "咱觉得你需要通过别的方式放松哦，而不是看涩图"}[rand.Intn(6)]))
		}
	})
}
