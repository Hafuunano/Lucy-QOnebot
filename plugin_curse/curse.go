// Package curse 骂人插件(求骂,自卫)
package curse

import (
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"

	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/process"

	"github.com/FloatTech/ZeroBot-Plugin/order"
)

const (
	minLevel = "min"
	maxLevel = "max"
)

func init() {
	limit := rate.NewManager(time.Minute, 5)
	engine := control.Register("curse", order.PrioCurse, &control.Options{
		DisableOnDefault: true,
		Help:             "骂人\n- 骂我\n- 大力骂我",
	})

	engine.OnFullMatch("骂我").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.GroupID).Acquire() {
			ctx.SendChain(message.Text("(,,´•ω•)ノ(´っω•｀。)\n 不要骂自己了啦 摸摸头"))
			return
		}
		process.SleepAbout1sTo2s()
		text := getRandomCurseByLevel(minLevel).Text
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
	})

	engine.OnFullMatch("大力骂我").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.GroupID).Acquire() {
			ctx.SendChain(message.Text("(,,´•ω•)ノ(´っω•｀。)\n 不要骂自己了啦 摸摸头"))
			return
		}
		process.SleepAbout1sTo2s()
		text := getRandomCurseByLevel(maxLevel).Text
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(text))
	})
}
